package mysql

import (
	"fmt"
	"github.com/romberli/db-operator/pkg/message"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"github.com/romberli/db-operator/module/implement/mysql/parameter"
	"github.com/romberli/db-operator/pkg/util/ssh"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/log"
	"github.com/spf13/viper"

	msgMySQL "github.com/romberli/db-operator/pkg/message/mysql"
)

const (
	installSuccessMessage = "install mysql server completed."

	defaultUseSudo = true

	addrTemplate                         = "%s:%d"
	defaultConfigFileName                = "/etc/my.cnf"
	defaultConfigFileBackupNameTemplate  = "/etc/my.cnf.%s"
	configFileNameTemplate               = "my.cnf.%d"
	mysqldMultiTitleTemplate             = "mysqld%d"
	mysqldSingleInstanceSectionTemplate  = "[mysqld]"
	mysqldMultiInstanceSectionTemplate   = "[mysqld%d]"
	mysqldMultiInstanceIsRunningTemplate = "MySQL server from group: mysqld%d is running"

	getMySQLPIDListCommandTemplate     = `/usr/bin/ps -ef | /usr/bin/grep mysqld | /usr/bin/grep %d | /usr/bin/grep %s | /usr/bin/grep -v grep | /usr/bin/awk -F' ' '{print \$2}'`
	initMySQLInstanceCommandTemplate   = "%s/bin/mysqld --defaults-file=/tmp/my.cnf.%d --initialize --basedir=%s --datadir=%s/data --user=%s"
	getDefaultRootPassCommandTemplate  = `grep 'A temporary password is generated for root@localhost' %s/%s/mysql.err | awk -F' ' '{print \$NF}'`
	startSingleInstanceCommandTemplate = "%s/bin/mysqld --defaults-file=/tmp/my.cnf.%d --basedir=%s --datadir=%s/data --user=%s &"
	startMultiInstanceCommandTemplate  = "export PATH=$PATH:%s/bin && mysqld_multi start %d"
	checkMultiInstanceCommandTemplate  = `export PATH=$PATH:%s/bin && mysqld_multi report %d | /usr/bin/grep \"MySQL server from group\"`
	initMySQLUserCommandTemplate       = `%s/bin/mysql --connect-expired-password -uroot -p'%s' -S %s/run/mysql.sock -e \"%s\"`

	shutdownSQL     = "shutdown ;"
	startReplicaSQL = "start replica ;"
	stopReplicaSQL  = "stop replica ;"

	changeMasterSQLTemplate    = "change master to master_host='%s', master_port=%d, master_user='%s', master_password='%s', master_auto_position=1 ;"
	SlaveIOThreadRunningField  = "Slave_IO_Running"
	SlaveSQLThreadRunningField = "Slave_SQL_Running"
	IsRunningValue             = "Yes"

	maxRetryCount = 5
	retryInterval = 2 * time.Second
)

type Engine struct {
	dboRepo      *DBORepo
	ose          *OSExecutor
	mysqlVersion *version.Version
	Mode         mode.Mode              `json:"mode"`
	Addrs        []string               `json:"addrs"`
	MySQLServer  *parameter.MySQLServer `json:"mysql_server"`
	PMMClient    *parameter.PMMClient   `json:"pmm_client"`
}

// NewEngine returns a new *Engine
func NewEngine(dboRepo *DBORepo, mysqlVersion *version.Version, m mode.Mode, addrs []string, mysqlServer *parameter.MySQLServer, pmmClient *parameter.PMMClient) *Engine {
	return newEngine(dboRepo, mysqlVersion, m, addrs, mysqlServer, pmmClient)
}

// NewEngineWithDefault returns a new *Engine with default values
func NewEngineWithDefault(mysqlVersion *version.Version, m mode.Mode, addrs []string, mysqlServer *parameter.MySQLServer, pmmClient *parameter.PMMClient) *Engine {
	return newEngine(
		NewDBORepoWithDefault(),
		mysqlVersion,
		m,
		addrs,
		mysqlServer,
		pmmClient,
	)
}

// newEngine returns a new *Engine
func newEngine(dboRepo *DBORepo, mysqlVersion *version.Version, m mode.Mode, addrs []string, mysqlServer *parameter.MySQLServer, pmmClient *parameter.PMMClient) *Engine {
	return &Engine{
		dboRepo:      dboRepo,
		mysqlVersion: mysqlVersion,
		Mode:         m,
		Addrs:        addrs,
		MySQLServer:  mysqlServer,
		PMMClient:    pmmClient,
	}
}

// Install installs mysql to the hosts
func (e *Engine) Install(operationID int) error {
	err := linux.SortAddrs(e.Addrs)
	if err != nil {
		return err
	}

	var (
		sourceHostIP  string
		sourcePortNum int
	)

	for i, addr := range e.Addrs {
		var (
			isSource          bool
			hostIP            string
			portNumStr        string
			portNum           int
			operationDetailID int
		)

		hostIP, portNumStr, err = net.SplitHostPort(addr)
		if err != nil {
			return errors.Trace(err)
		}
		if hostIP == constant.EmptyString || portNumStr == constant.EmptyString {
			return errors.Errorf("mysql Engine.Install(): addr must be formatted as host:port, %s is invalid", addr)
		}
		portNum, err = strconv.Atoi(portNumStr)
		if err != nil {
			return errors.Trace(err)
		}

		if i == constant.ZeroInt {
			isSource = true
			sourceHostIP = hostIP
			sourcePortNum = portNum
		}

		// init operation detail
		operationDetailID, err = e.dboRepo.InitOperationDetail(operationID, e.MySQLServer.HostIP, e.MySQLServer.PortNum)
		if err != nil {
			return err
		}
		// install single instance
		err = e.InstallSingleInstance(hostIP, portNum, isSource)
		if err != nil {
			updateErr := e.dboRepo.UpdateOperationDetail(operationDetailID, defaultFailedStatus, err.Error())
			if updateErr != nil {
				log.Errorf(constant.LogWithStackString, message.NewMessage(msgMySQL.ErrMySQLEngineUpdateOperationDetail,
					updateErr, operationID, operationDetailID, hostIP, portNum, defaultFailedStatus))
			}

			return err
		}

		if !isSource && e.Mode == mode.AsyncReplication || e.Mode == mode.SemiSyncReplication {
			// configure mysql replica
			err = e.ConfigureReplica(addr, sourceHostIP, sourcePortNum)
			if err != nil {
				updateErr := e.dboRepo.UpdateOperationDetail(operationDetailID, defaultFailedStatus, err.Error())
				if updateErr != nil {
					log.Errorf(constant.LogWithStackString, message.NewMessage(msgMySQL.ErrMySQLEngineUpdateOperationDetail,
						updateErr, operationID, operationDetailID, hostIP, portNum, defaultFailedStatus))
				}

				return err
			}
		}

		updateErr := e.dboRepo.UpdateOperationDetail(operationDetailID, defaultSuccessStatus, installSuccessMessage)
		if updateErr != nil {
			log.Errorf(constant.LogWithStackString, message.NewMessage(msgMySQL.ErrMySQLEngineUpdateOperationDetail,
				updateErr, operationID, operationDetailID, hostIP, portNum, defaultSuccessStatus))
		}

		log.Infof(message.NewMessage(msgMySQL.InfoMySQLEngineInitInstance, operationID, operationDetailID, hostIP, portNum).Error())
	}

	if e.Mode == mode.GroupReplication {
		// configure mysql group replication
		return e.ConfigureGroupReplication()
	}

	return nil
}

// InstallSingleInstance installs the single instance
func (e *Engine) InstallSingleInstance(hostIP string, portNum int, isSource bool) error {
	// reset MySQL Sever Parameter
	err := e.MySQLServer.InitWithHostInfo(hostIP, portNum, isSource)
	if err != nil {
		return err
	}
	// init os
	err = e.InitOS()
	if err != nil {
		return err
	}
	// init mysql instance
	err = e.InitMySQLInstance()
	if err != nil {
		return err
	}
	// init pmm client
	err = e.InitPMMClient()
	if err != nil {
		return err
	}

	return nil
}

// InitOS initializes the os
func (e *Engine) InitOS() error {
	err := e.InitOSExecutor()
	if err != nil {
		return err
	}

	return e.ose.Init()
}

// InitOSExecutor initializes the ssh connection
func (e *Engine) InitOSExecutor() error {
	sshConn, err := linux.NewSSHConn(
		e.MySQLServer.HostIP,
		constant.DefaultSSHPort,
		viper.GetString(config.MySQLUserOSUserKey),
		viper.GetString(config.MySQLUserOSPassKey),
		defaultUseSudo,
	)
	if err != nil {
		return err
	}

	e.ose = NewOSExecutor(ssh.NewConn(sshConn), e.mysqlVersion, e.MySQLServer)

	return nil
}

// InitMySQLInstance initializes the mysql instance
func (e *Engine) InitMySQLInstance() error {
	// prepare mysql multi instance config file
	err := e.prepareMultiInstanceConfigFile()
	if err != nil {
		return err
	}
	// init single instance
	rootPass, err := e.initMySQLInstance()
	if err != nil {
		return err
	}
	// start mysql single instance asynchronously
	go func() {
		err = e.startInstanceWithMySQLD()
		if err != nil {
			log.Errorf("mysql Engine.InitMySQLInstance(): start mysql instance failed. hostIP: %s, portNum: %d, error:\n%+v", e.MySQLServer.HostIP, e.MySQLServer.PortNum, err)
		}
	}()
	time.Sleep(retryInterval)
	// check instance status
	err = e.checkInstanceWithPID()
	if err != nil {
		return err
	}
	// sleep for a while to wait for mysql to be ready
	time.Sleep(retryInterval)
	// init mysql user
	err = e.initMySQLUser(rootPass)
	if err != nil {
		return err
	}
	// stop mysql single instance
	err = e.stopInstance()
	if err != nil {
		return err
	}
	err = e.waitForShuttingDown()
	// start mysql multi instance
	err = e.startInstanceWithMySQLDMulti()
	if err != nil {
		return err
	}
	time.Sleep(retryInterval)
	// check mysql multi instance
	isRunning, err := e.checkInstanceWithMySQLDMulti()
	if err != nil {
		return err
	}
	if !isRunning {
		return errors.Errorf("mysql Engine.InitMySQLInstance(): mysql multi instance is not running. hostIP: %s, portNum: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum)
	}

	return nil
}

// ConfigureReplica configures the replication
func (e *Engine) ConfigureReplica(addr, sourceHostIP string, sourcePortNum int) error {
	if addr == fmt.Sprintf(addrTemplate, sourceHostIP, sourcePortNum) {
		// this is the source node, do nothing
		return nil
	}

	conn, err := mysql.NewConn(addr, constant.EmptyString, constant.DefaultRootUserName, e.MySQLServer.RootPass)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("mysql Engine.ConfigureReplica(): close mysql connection failed. error:\n%+v", err)
		}
	}()

	sql := fmt.Sprintf(changeMasterSQLTemplate, sourceHostIP, sourcePortNum, e.MySQLServer.ReplicationUser, e.MySQLServer.ReplicationPass)
	_, err = conn.Execute(sql)
	if err != nil {
		return err
	}

	_, err = conn.Execute(startReplicaSQL)
	if err != nil {
		return err
	}

	var status string
	// check io thread
	for i := constant.ZeroInt; i < maxRetryCount; i++ {
		result, err := conn.GetReplicationSlavesStatus()
		if err != nil {
			return err
		}

		status, err = result.GetStringByName(constant.ZeroInt, SlaveIOThreadRunningField)
		if err != nil {
			return err
		}
		if status != IsRunningValue {
			log.Warnf("mysql Engine.ConfigureReplica(): slave io thread is not running, will be retry soon. hostIP: %s, portNum: %d, retryCount: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum, i)
			time.Sleep(retryInterval)
			continue
		}
		status, err = result.GetStringByName(constant.ZeroInt, SlaveSQLThreadRunningField)
		if err != nil {
			return err
		}
		if status == IsRunningValue {
			break
		}

		log.Warnf("mysql Engine.ConfigureReplica(): slave sql thread is not running, will be retry soon. hostIP: %s, portNum: %d, retryCount: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum, i)
		time.Sleep(retryInterval)
		continue
	}

	if status != IsRunningValue {
		return errors.Errorf("slave io thread is not running")
	}

	// check sql thread
	for i := constant.ZeroInt; i < maxRetryCount; i++ {
		result, err := conn.GetReplicationSlavesStatus()
		if err != nil {
			return err
		}

		status, err = result.GetStringByName(constant.ZeroInt, SlaveSQLThreadRunningField)
		if err != nil {
			return err
		}

		if status == IsRunningValue {
			break
		}

		log.Warnf("mysql Engine.ConfigureReplica(): slave sql thread is not running, will be retry soon. hostIP: %s, portNum: %d, retryCount: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum, i)
		time.Sleep(retryInterval)
		continue
	}

	if status != IsRunningValue {
		return errors.Errorf("mysql Engine.ConfigureReplica(): slave sql thread is not running")
	}

	return nil
}

// InitPMMClient initializes the pmm client
func (e *Engine) InitPMMClient() error {
	pmmExecutor := NewPMMExecutor(e.ose.Conn, e.MySQLServer.HostIP, e.MySQLServer.PortNum, e.PMMClient)

	return pmmExecutor.Init()
}

// ConfigureGroupReplication configures the group replication
func (e *Engine) ConfigureGroupReplication() error {
	// TODO: implement this
	return errors.Errorf("mysql Engine.ConfigureGroupReplication(): group replication has not been implemented")
}

// initMySQLInstance initializes the mysql instance
func (e *Engine) initMySQLInstance() (string, error) {
	// prepare init config file
	err := e.prepareInitConfigFile()
	if err != nil {
		return constant.EmptyString, err
	}
	// init mysql instance
	cmd := fmt.Sprintf(initMySQLInstanceCommandTemplate, e.MySQLServer.BinaryDirBase,
		e.MySQLServer.PortNum, e.MySQLServer.BinaryDirBase, e.MySQLServer.DataDirBase, defaultMySQLUser)
	err = e.ose.Conn.ExecuteCommandWithoutOutput(cmd)
	if err != nil {
		return constant.EmptyString, err
	}

	// return the default root password
	return e.getDefaultMySQLRootPass()
}

// startInstanceWithMySQLD starts the instance with mysqld
func (e *Engine) startInstanceWithMySQLD() error {
	cmd := fmt.Sprintf(startSingleInstanceCommandTemplate, e.MySQLServer.BinaryDirBase, e.MySQLServer.PortNum,
		e.MySQLServer.BinaryDirBase, e.MySQLServer.DataDirBase, defaultMySQLUser)

	return e.ose.Conn.ExecuteCommandWithoutOutput(cmd)
}

// checkInstanceWithPID checks the instance with pid
func (e *Engine) checkInstanceWithPID() error {
	for i := constant.ZeroInt; i < maxRetryCount; i++ {
		pidList, err := e.ose.GetMySQLPIDList()
		if err != nil {
			return err
		}

		if len(pidList) > constant.ZeroInt {
			return nil
		}

		log.Warnf("mysql Engine.checkInstanceWithPID(): no mysqld pid found, will be retry soon. hostIP: %s, portNum: %d, retryCount: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum, i)
		time.Sleep(retryInterval)
	}

	return errors.Errorf("mysql Engine.checkInstanceWithPID(): maximum retry count of checking mysql pid exceeded, but still no mysqld pid found. hostIP: %s, portNum: %d, maxRetryCount: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum, maxRetryCount)
}

// initMySQLUser initializes the user
func (e *Engine) initMySQLUser(rootPass string) error {
	sql, err := e.MySQLServer.GetInitUserSQL()
	if err != nil {
		return err
	}

	command := fmt.Sprintf(initMySQLUserCommandTemplate, e.MySQLServer.BinaryDirBase, rootPass, e.MySQLServer.DataDirBase, sql)

	return e.ose.Conn.ExecuteCommandWithoutOutput(command)
}

// waitForShuttingDown waits for the instance to shut down
func (e *Engine) waitForShuttingDown() error {
	for i := constant.ZeroInt; i < maxRetryCount; i++ {
		pidList, err := e.ose.GetMySQLPIDList()
		if err != nil {
			return err
		}

		if len(pidList) == constant.ZeroInt {
			return nil
		}

		log.Warnf("mysql Engine.waitForShuttingDown(): mysqld pid found, will be retry soon. hostIP: %s, portNum: %d, maxRetryCount: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum, i)
		time.Sleep(retryInterval)
	}

	return errors.Errorf("mysql Engine.waitForShuttingDown(): maximum retry count of waiting for shutting down exceeded, but still found mysqld pid. hostIP: %s, portNum: %d, maxRetryCount: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum, maxRetryCount)
}

// prepareMultiInstanceConfigFile prepares mysql config file
func (e *Engine) prepareMultiInstanceConfigFile() error {
	// check if the config file exists
	exists, err := e.ose.Conn.PathExists(defaultConfigFileName)
	if err != nil {
		return err
	}

	title := e.getMultiInstanceTitle()

	if !exists {
		// the config file does not exist, generate a new one
		configBytes, err := e.MySQLServer.GetConfigWithTitle(title, e.mysqlVersion, e.Mode)
		if err != nil {
			return err
		}

		return e.transferConfigContent(configBytes, fmt.Sprintf(configFileNameTemplate, e.MySQLServer.PortNum), defaultConfigFileName)
	}

	// the config file exists
	// backup the config file
	err = e.ose.Conn.Copy(defaultConfigFileName, fmt.Sprintf(defaultConfigFileBackupNameTemplate, time.Now().Format(constant.TimeLayoutSecondDash)))
	if err != nil {
		return err
	}
	// get the config file content
	existingContent, err := e.ose.Conn.Cat(defaultConfigFileName)
	if err != nil {
		return err
	}

	if strings.Contains(existingContent, mysqldSingleInstanceSectionTemplate) {
		// mysqld section exists
		return errors.New("mysql Engine.prepareMultiInstanceConfigFile(): mysqld section exists, db operator does not support converting the single instance to multi instance")
	}

	if !strings.Contains(existingContent, fmt.Sprintf(mysqldMultiInstanceSectionTemplate, e.MySQLServer.PortNum)) {
		// the instance section does not exist
		newContent, err := e.MySQLServer.GetMySQLDConfigWithTitle(title, e.mysqlVersion, e.Mode)
		if err != nil {
			return err
		}
		// append the instance section to the config file
		content := existingContent + string(newContent)
		fileName := fmt.Sprintf(configFileNameTemplate, e.MySQLServer.PortNum)

		return e.transferConfigContent([]byte(content), fileName, defaultConfigFileName)
	}

	// the instance section exists, do nothing
	return nil
}

// stopInstance stops the instance
func (e *Engine) stopInstance() error {
	// connect to the mysql instance
	conn, err := mysql.NewConn(
		fmt.Sprintf("%s:%d", e.MySQLServer.HostIP, e.MySQLServer.PortNum),
		constant.EmptyString,
		constant.DefaultRootUserName,
		e.MySQLServer.RootPass,
	)
	if err != nil {
		return err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("Engine.stopInstance(): close mysql connection failed. error:\n%+v", err)
		}
	}()

	// stop the mysql instance
	_, err = conn.Execute(shutdownSQL)

	return err
}

// startInstanceWithMySQLDMulti starts the instance with mysqld_multi
func (e *Engine) startInstanceWithMySQLDMulti() error {
	cmd := fmt.Sprintf(startMultiInstanceCommandTemplate, e.MySQLServer.BinaryDirBase, e.MySQLServer.PortNum)

	return e.ose.Conn.ExecuteCommandWithoutOutput(cmd)
}

// checkInstanceWithMySQLDMulti checks the instance with mysqld_multi
func (e *Engine) checkInstanceWithMySQLDMulti() (bool, error) {
	cmd := fmt.Sprintf(checkMultiInstanceCommandTemplate, e.MySQLServer.BinaryDirBase, e.MySQLServer.PortNum)
	instanceRunning := fmt.Sprintf(mysqldMultiInstanceIsRunningTemplate, e.MySQLServer.PortNum)

	for i := constant.ZeroInt; i < maxRetryCount; i++ {
		output, err := e.ose.Conn.ExecuteCommand(cmd)
		if err != nil {
			return false, err
		}

		if output == instanceRunning {
			return true, nil
		}

		log.Warnf("mysql Engine.checkInstanceWithMySQLDMulti(): mysqld multi instance is not running. hostIP: %s, portNum: %d, retryCount: %d", e.MySQLServer.HostIP, e.MySQLServer.PortNum, i)
		time.Sleep(retryInterval)
		continue
	}

	return false, nil
}

// prepareInitConfigFile prepares mysql config file for initializing mysql instance
func (e *Engine) prepareInitConfigFile() error {
	configBytes, err := e.MySQLServer.GetConfig(e.mysqlVersion, e.Mode)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf(configFileNameTemplate, e.MySQLServer.PortNum)
	fileDest := filepath.Join(constant.DefaultTmpDir, fileName)

	return e.transferConfigContent(configBytes, fileName, fileDest)
}

// transferConfigContent transfers the config content to the remote host
func (e *Engine) transferConfigContent(configContent []byte, fileNameSource, filePathDest string) error {
	fileSource, err := os.CreateTemp(viper.GetString(config.MySQLInstallationTemporaryDirKey), fileNameSource)
	if err != nil {
		return err
	}
	defer func() {
		err = fileSource.Close()
		if err != nil {
			log.Errorf("Engine.transferConfigContent(): close file source failed. error:\n%+v", err)
		}
		err = os.Remove(fileSource.Name())
		if err != nil {
			log.Errorf("Engine.transferConfigContent(): remove file source failed. error:\n%+v", err)
		}
	}()

	_, err = fileSource.Write(configContent)
	if err != nil {
		return err
	}
	err = e.ose.Conn.CopySingleFileToRemote(fileSource.Name(), filePathDest, constant.DefaultTmpDir)
	if err != nil {
		return err
	}

	return e.ose.Conn.Chown(filePathDest, defaultMySQLUser, defaultMySQLUser)
}

// getMultiInstanceTitle gets the title of the instance
func (e *Engine) getMultiInstanceTitle() string {
	return fmt.Sprintf(mysqldMultiTitleTemplate, e.MySQLServer.PortNum)
}

// getDefaultMySQLRootPass gets the default mysql root password
func (e *Engine) getDefaultMySQLRootPass() (string, error) {
	output, err := e.ose.Conn.ExecuteCommand(fmt.Sprintf(getDefaultRootPassCommandTemplate, e.MySQLServer.DataDirBase, logDirName))
	if err != nil {
		return constant.EmptyString, err
	}

	return strings.TrimSpace(output), nil
}

// getSourceNode gets the source node
func (e *Engine) getSourceNode() (string, int, error) {
	err := linux.SortAddrs(e.Addrs)
	if err != nil {
		return constant.EmptyString, constant.ZeroInt, err
	}

	addr := e.Addrs[constant.ZeroInt]
	hostIP, portNumStr, err := net.SplitHostPort(addr)
	if err != nil {
		return constant.EmptyString, constant.ZeroInt, err
	}

	portNum, err := strconv.Atoi(portNumStr)
	if err != nil {
		return constant.EmptyString, constant.ZeroInt, err
	}

	return hostIP, portNum, nil
}
