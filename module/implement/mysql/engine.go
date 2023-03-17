package mysql

import (
	"fmt"
	"github.com/romberli/go-util/middleware/mysql"
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
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/log"
	"github.com/spf13/viper"
)

const (
	addrTemplate                         = "%s:%d"
	defaultConfigFileName                = "/etc/my.cnf"
	defaultConfigFileBackupNameTemplate  = "/etc/my.cnf.%s"
	configFileNameTemplate               = "my.cnf.%d"
	mysqldMultiTitleTemplate             = "mysqld%d"
	mysqldSingleInstanceSectionTemplate  = "[mysqld]"
	mysqldMultiInstanceSectionTemplate   = "[mysqld%d]"
	mysqldMultiInstanceIsRunningTemplate = "mysqld%d is running"

	initMySQLInstanceCommandTemplate   = "%s/bin/mysqld --defaults-file=/tmp/my.cnf.%d --initialize --basedir=%s --datadir=%s/data --user=%s"
	getDefaultRootPassCommandTemplate  = "grep 'A temporary password is generated for root@localhost' %s/%s/mysql.err | awk -F' ' '{print $NF}'"
	startSingleInstanceCommandTemplate = "%s/bin/mysqld --defaults-file=/tmp/my.cnf.%d --basedir=%s --datadir=%s/data --user=%s"
	startMultiInstanceCommandTemplate  = "%s/bin/mysqld_multi start %d"
	checkMultiInstanceCommandTemplate  = "%s/bin/mysqld_multi report %d "
	initMySQLUserCommandTemplate       = `%s/bin/mysql --connect-expired-password -uroot -p'%s' -S %s/run/mysql.sock -e "%s"`

	shutdownSQL     = "shutdown ;"
	startReplicaSQL = "start replica ;"
	stopReplicaSQL  = "stop replica ;"

	changeMasterSQLTemplate = "change master to master_host='%s', master_port=%d, master_user='%s', master_password='%s', master_auto_position=1 ;"
	SlaveThreadIsRunning    = "Slave_IO_Running: YES, Slave_SQL_Running: YES"
)

type Engine struct {
	dboRepo      *DBORepo
	mysqlVersion *version.Version
	Mode         mode.Mode              `json:"mode"`
	Addrs        []string               `json:"addrs"`
	MySQLServer  *parameter.MySQLServer `json:"mysql_server"`
	PMMClient    *parameter.PMMClient   `json:"pmm_client"`
}

// NewEngine returns a new *Engine
func NewEngine(dboRepo *DBORepo, m mode.Mode, addrs []string, mysqlServer *parameter.MySQLServer, pmmClient *parameter.PMMClient) *Engine {
	return newEngine(dboRepo, m, addrs, mysqlServer, pmmClient)
}

// NewEngineWithDefault returns a new *Engine with default values
func NewEngineWithDefault(m mode.Mode, addrs []string, mysqlServer *parameter.MySQLServer, pmmClient *parameter.PMMClient) *Engine {
	return NewEngine(
		NewDBORepoWithDefault(),
		m,
		addrs,
		mysqlServer,
		pmmClient,
	)
}

// newEngine returns a new *Engine
func newEngine(dboRepo *DBORepo, m mode.Mode, addrs []string, mysqlServer *parameter.MySQLServer, pmmClient *parameter.PMMClient) *Engine {
	return &Engine{
		dboRepo:     dboRepo,
		Mode:        m,
		Addrs:       addrs,
		MySQLServer: mysqlServer,
		PMMClient:   pmmClient,
	}
}

// Install installs mysql to the hosts
func (e *Engine) Install() error {

	var (
		sourceHostIP  string
		sourcePortNum int
	)

	for i, addr := range e.Addrs {
		var isSource bool

		hostIP, portNumStr, err := net.SplitHostPort(addr)
		if err != nil {
			return errors.Trace(err)
		}
		if hostIP == constant.EmptyString || portNumStr == constant.EmptyString {
			return errors.Errorf("addr must be formatted as host:port, %s is invalid", addr)
		}
		portNum, err := strconv.Atoi(portNumStr)
		if err != nil {
			return errors.Trace(err)
		}

		// set MySQL Sever Parameter
		e.MySQLServer.SetHostIP(hostIP)
		e.MySQLServer.SetPortNum(portNum)

		if i == constant.ZeroInt {
			isSource = true
			sourceHostIP = hostIP
			sourcePortNum = portNum
		}

		// init os
		err = e.InitOS(hostIP)
		if err != nil {
			return err
		}
		// init mysql instance
		err = e.InitMySQLInstance()
		if err != nil {
			return err
		}

		if !isSource {
			// configure mysql replica
			err = e.ConfigureMySQLCluster(addr, sourceHostIP, sourcePortNum)
			if err != nil {
				return err
			}
		}
		// init pmm client
		err = e.InitPMMClient()
		if err != nil {
			return err
		}
	}

	return nil
}

// InitOS initializes the os
func (e *Engine) InitOS(hostIP string) error {
	sshConn := linux.NewSSHConn(hostIP)
	osExecutor := NewOSExecutor(e.mysqlVersion, e.MySQLServer)

	return osExecutor.Init()
}

// InitMySQLInstance initializes the mysql instance
func (e *Engine) InitMySQLInstance(isSource bool) error {
	// prepare mysql multi instance config file
	err := e.prepareMultiInstanceConfigFile(isSource)
	// init single instance
	rootPass, err := e.initMySQLInstance()
	if err != nil {
		return err
	}
	// start mysql single instance
	err = e.startInstanceWithMySQLD()
	if err != nil {
		return err
	}
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
	// start mysql multi instance
	err = e.startInstanceWithMySQLDMulti()
	if err != nil {
		return err
	}
	// check mysql multi instance
	isRunning, err := e.checkInstanceWithMySQLDMulti()
	if err != nil {
		return err
	}
	if !isRunning {
		return errors.Errorf("mysql multi instance is not running. port_num: %d", e.MySQLServer.PortNum)
	}

	return nil
}

// initMySQLInstance initializes the mysql instance
func (e *Engine) initMySQLInstance() (string, error) {
	// prepare init config file
	err := e.prepareInitConfigFile()
	if err != nil {
		return constant.EmptyString, err
	}
	// init mysql instance
	cmd := fmt.Sprintf(initMySQLInstanceCommandTemplate, e.MySQLServer.BinaryDirBase, e.MySQLServer.PortNum, e.MySQLServer.BinaryDirBase, e.MySQLServer.DataDirBase, defaultMySQLUser)
	err = e.sshConn.ExecuteCommandWithoutOutput(cmd)
	if err != nil {
		return constant.EmptyString, err
	}

	// return the default root password
	return e.getDefaultMySQLRootPass()
}

// startInstanceWithMySQLD starts the instance with mysqld
func (e *Engine) startInstanceWithMySQLD() error {
	cmd := fmt.Sprintf(startSingleInstanceCommandTemplate, e.MySQLServer.BinaryDirBase, e.MySQLServer.PortNum, e.MySQLServer.BinaryDirBase, e.MySQLServer.DataDirBase, defaultMySQLUser)

	return e.sshConn.ExecuteCommandWithoutOutput(cmd)
}

// initMySQLUser initializes the user
func (e *Engine) initMySQLUser(rootPass string) error {
	sql, err := e.MySQLServer.GetInitUserSQL()
	if err != nil {
		return err
	}

	command := fmt.Sprintf(initMySQLUserCommandTemplate, e.MySQLServer.BinaryDirBase, rootPass, e.MySQLServer.DataDirBase, sql)

	return e.sshConn.ExecuteCommandWithoutOutput(command)
}

// prepareMultiInstanceConfigFile prepares mysql config file
func (e *Engine) prepareMultiInstanceConfigFile(isSource bool) error {
	// check if the config file exists
	exists, err := e.sshConn.PathExists(defaultConfigFileName)
	if err != nil {
		return err
	}

	if !exists {
		// the config file does not exist, generate a new one
		if isSource {
			e.MySQLServer.SetSemiSyncSourceEnabled(constant.OneInt)
			e.MySQLServer.SetSemiSyncReplicaEnabled(constant.ZeroInt)
		}
		configBytes, err := e.MySQLServer.GetConfig(e.mysqlVersion, e.Mode)
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf(configFileNameTemplate, e.MySQLServer.PortNum)
		fileDest := filepath.Join(constant.DefaultTmpDir, fileName)

		return e.transferConfigContent(string(configBytes), fileName, fileDest)
	}

	// the config file exists
	// backup the config file
	err = e.sshConn.Copy(defaultConfigFileName,
		fmt.Sprintf(defaultConfigFileBackupNameTemplate, time.Now().Format(constant.TimeLayoutSecondDash)))
	if err != nil {
		return err
	}
	// get the config file content
	existingContent, err := e.sshConn.Cat(defaultConfigFileName)
	if err != nil {
		return err
	}

	if strings.Contains(existingContent, mysqldSingleInstanceSectionTemplate) {
		// mysqld section exists
		return errors.New("mysqld section exists, db operator does not support converting the single instance to multi instance")
	}

	if !strings.Contains(existingContent, fmt.Sprintf(mysqldMultiInstanceSectionTemplate, e.MySQLServer.PortNum)) {
		// the instance section does not exist
		newContent, err := e.MySQLServer.GetMySQLDConfigWithTitle(e.getTitle(), e.mysqlVersion, e.Mode)
		if err != nil {
			return err
		}
		// append the instance section to the config file
		content := existingContent + string(newContent)
		fileName := fmt.Sprintf(configFileNameTemplate, e.MySQLServer.PortNum)

		return e.transferConfigContent(content, fileName, defaultConfigFileName)
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
			log.Errorf("Engine.initMySQLUser(): close mysql connection failed. error:\n%+v", err)
		}
	}()

	// stop the mysql instance
	_, err = conn.Execute(shutdownSQL)

	return err
}

// startInstanceWithMySQLDMulti starts the instance with mysqld_multi
func (e *Engine) startInstanceWithMySQLDMulti() error {
	cmd := fmt.Sprintf(startMultiInstanceCommandTemplate, e.MySQLServer.BinaryDirBase, e.MySQLServer.PortNum)

	return e.sshConn.ExecuteCommandWithoutOutput(cmd)
}

// checkInstanceWithMySQLDMulti checks the instance with mysqld_multi
func (e *Engine) checkInstanceWithMySQLDMulti() (bool, error) {
	cmd := fmt.Sprintf(checkMultiInstanceCommandTemplate, e.MySQLServer.BinaryDirBase, e.MySQLServer.PortNum)

	output, err := e.sshConn.ExecuteCommand(cmd)
	if err != nil {
		return false, err
	}

	return output == fmt.Sprintf(mysqldMultiInstanceIsRunningTemplate, e.MySQLServer.PortNum), nil
}

// ConfigureMySQLCluster configures the mysql cluster
func (e *Engine) ConfigureMySQLCluster(addr, sourceHostIP string, sourcePortNum int) error {
	switch e.Mode {
	case mode.Standalone:
		return nil
	case mode.AsyncReplication, mode.SemiSyncReplication:

		return e.configureReplication(addr, sourceHostIP, sourcePortNum)
	case mode.GroupReplication:
		return e.configureGroupReplication()
	default:
		return errors.Errorf("unsupported mode %d", e.Mode)
	}
}

// InitPMMClient initializes the pmm client
func (e *Engine) InitPMMClient() error {
	pmmExecutor := NewPMMExecutor(e.sshConn, e.MySQLServer.HostIP, e.MySQLServer.PortNum, e.PMMClient)

	return pmmExecutor.Init()
}

// prepareInitConfigFile prepares mysql config file for initializing mysql instance
func (e *Engine) prepareInitConfigFile() error {
	configBytes, err := e.MySQLServer.GetConfig(e.mysqlVersion, e.Mode)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf(configFileNameTemplate, e.MySQLServer.PortNum)
	fileDest := filepath.Join(constant.DefaultTmpDir, fileName)

	return e.transferConfigContent(string(configBytes), fileName, fileDest)
}

// transferConfigContent transfers the config content to the remote host
func (e *Engine) transferConfigContent(configContent string, fileNameSource, filePathDest string) error {
	fileSource, err := os.CreateTemp(viper.GetString(config.MySQLInstallationTemporaryDirKey), fileNameSource)
	if err != nil {
		return err
	}
	defer func() {
		err = fileSource.Close()
		if err != nil {
			log.Errorf("Engine.prepareInitConfigFile(): close file source failed. error:\n%+v", err)
		}
		err = os.Remove(fileSource.Name())
		if err != nil {
			log.Errorf("Engine.prepareInitConfigFile(): remove file source failed. error:\n%+v", err)
		}
	}()

	_, err = fileSource.WriteString(configContent)
	if err != nil {
		return err
	}
	err = e.sshConn.CopySingleFileToRemote(fileSource.Name(), filePathDest, constant.DefaultTmpDir)
	if err != nil {
		return err
	}

	return e.sshConn.Chown(defaultConfigFileName, defaultMySQLUser, defaultMySQLUser)
}

// getTitle gets the title of the instance
func (e *Engine) getTitle() string {
	return fmt.Sprintf(mysqldMultiTitleTemplate, e.MySQLServer.PortNum)
}

// getDefaultMySQLRootPass gets the default mysql root password
func (e *Engine) getDefaultMySQLRootPass() (string, error) {
	output, err := e.sshConn.ExecuteCommand(fmt.Sprintf(getDefaultRootPassCommandTemplate, e.MySQLServer.DataDirBase, logDirName))
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

// configureReplication configures the replication
func (e *Engine) configureReplication(addr, sourceHostIP string, sourcePortNum int) error {
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
			log.Errorf("Engine.configureReplication(): close mysql connection failed. error:\n%+v", err)
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

	result, err := conn.GetReplicationSlavesStatus()
	if err != nil {
		return err
	}
	status, err := result.GetString(constant.ZeroInt, constant.ZeroInt)
	if err != nil {
		return err
	}
	if !strings.Contains(status, SlaveThreadIsRunning) {
		return errors.Errorf("slave thread is not running")
	}

	return nil
}

// configureGroupReplication configures the group replication
func (e *Engine) configureGroupReplication() error {
	return nil
}
