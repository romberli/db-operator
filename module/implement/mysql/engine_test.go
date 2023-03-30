package mysql

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/global"
	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"github.com/romberli/db-operator/module/implement/mysql/parameter"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	testDBDBOMySQLAddr = "192.168.137.11:3306"
	testDBDBOMySQLName = "dbo"
	testDBDBOMySQLUser = "root"
	testDBDBOMySQLPass = "root"

	testMySQLInstallationPackageDir = "/data/software/mysql"

	testHostIP1         = "192.168.137.21"
	testPortNum1        = 3306
	testHostIP2         = "192.168.137.21"
	testPortNum2        = 3307
	testHostIP3         = "192.168.137.21"
	testPortNum3        = 3308
	testDataDirBaseName = "/data/mysql/data"
	testLogDirBaseName  = "/data/mysql/data"
	testRootPaas        = "root"
	testAdminUser       = "admin"
	testAdminPass       = "admin"
	testClientUser      = "root"
	testClientPass      = "root"
	testMySQLDMultiUser = "mysqld_multi"
	testMySQLDMultiPass = "mysqld_multi"
	testMonitorUser     = "pmm"
	testMonitorPass     = "pmm"
	testReplicationUser = "replication"
	testReplicationPass = "replication"
	testDASUser         = "das"
	testDASPass         = "das"

	defaultTitle                        = "mysqld"
	testVersion                         = "8.0.32"
	testBinaryDirBase                   = "/data/mysql/mysql8.0.32"
	testSemiSyncSourceEnabled           = 1
	testSemiSyncReplicaEnabled          = 0
	testSemiSyncSourceTimeout           = 30000
	testGroupReplicationConsistency     = "eventual"
	testGroupReplicationFlowControlMode = "disabled"
	testGroupReplicationMemberWeight    = 50
	testBinlogExpireLogsSeconds         = 604800
	testBinlogExpireLogsDays            = 7
	testBackupDir                       = "/data/backup"
	testMaxConnections                  = 128
	testInnodbBufferPoolSize            = "1G"
	testInnodbIOCapacity                = 1000

	testOSUser  = "dba"
	testOSPass  = "dba"
	testUseSudo = true

	testArch            = "aarch64"
	testOSVersionStr    = "9.0"
	testMySQLVersionStr = "8.0.32"
	testMode            = mode.AsyncReplication

	testPMMServerAddr      = "192.168.137.11:443"
	testPMMServerUser      = "admin"
	testPMMServerPass      = "admin"
	testPMMClientVersion   = "2.24.0"
	testPMMServiceName     = "192-168-137-11-3306"
	testReplicationSetName = constant.EmptyString

	testDateCommand          = "date +%Y%m%d-%H%M%S"
	testGetPIDCommand        = `ps -ef | grep mysqld | grep %d | grep -v grep | awk -F' ' '{print \$3}'`
	testKillMySQLDCommand    = "kill -9 %d"
	testRemoveDataDirCommand = "rm -rf %s"
	testRemoveBaseDirCommand = "rm -rf %s"
	testRemove

	testCheckSlaveStatusSQL = "show slave status"
)

var (
	testAddr1 = fmt.Sprintf("%s:%d", testHostIP1, testPortNum1)
	testAddr2 = fmt.Sprintf("%s:%d", testHostIP2, testPortNum2)
	testAddr3 = fmt.Sprintf("%s:%d", testHostIP3, testPortNum3)
	testAddrs = []string{testAddr1, testAddr2}

	testOSVersion    = version.Must(version.NewVersion(testOSVersionStr))
	testMySQLVersion = version.Must(version.NewVersion(testMySQLVersionStr))

	testServerID    int
	testMySQLServer *parameter.MySQLServer
	testEngine      *Engine
)

func init() {
	testInitViper()
	testInitDBOMySQLPool()
	testServerID = testInitServerID(testHostIP1, testPortNum1)
	testMySQLServer = testInitMySQLServer(testHostIP1, testPortNum1, testServerID)
	testPMMClient = testInitPMMClient()
	testConn = testInitSSHConn(testHostIP1)
	testOSExecutor = testInitOSExecutor()

	testEngine = testInitEngine()
}

func testInitViper() {
	// set global
	viper.Set(config.DBDBOMySQLAddrKey, testDBDBOMySQLAddr)
	viper.Set(config.DBDBOMySQLNameKey, testDBDBOMySQLName)
	viper.Set(config.DBDBOMySQLUserKey, testDBDBOMySQLUser)
	viper.Set(config.DBDBOMySQLPassKey, testDBDBOMySQLPass)
	viper.Set(config.DBPoolMaxConnectionsKey, mysql.DefaultMaxConnections)
	viper.Set(config.DBPoolInitConnectionsKey, mysql.DefaultInitConnections)
	viper.Set(config.DBPoolMaxIdleConnectionsKey, mysql.DefaultMaxIdleConnections)
	viper.Set(config.DBPoolMaxIdleTimeKey, mysql.DefaultMaxIdleTime)
	viper.Set(config.DBPoolMaxWaitTimeKey, mysql.DefaultMaxWaitTime)
	viper.Set(config.DBPoolMaxRetryCountKey, mysql.DefaultMaxRetryCount)
	viper.Set(config.DBPoolKeepAliveIntervalKey, mysql.DefaultKeepAliveInterval)

	viper.Set(config.MySQLInstallationPackageDirKey, testMySQLInstallationPackageDir)
	viper.Set(config.MySQLUserOSUserKey, testOSUser)
	viper.Set(config.MySQLUserOSPassKey, testOSPass)
	viper.Set(config.MySQLUserMonitorUserKey, testMonitorUser)
	viper.Set(config.MySQLUserMonitorPassKey, testMonitorPass)
	viper.Set(config.PMMServerAddrKey, testPMMServerAddr)
	viper.Set(config.PMMServerUserKey, testPMMServerUser)
	viper.Set(config.PMMServerPassKey, testPMMServerPass)
	viper.Set(config.PMMClientVersionKey, testPMMClientVersion)
}

func testInitDBOMySQLPool() {
	if global.DBOMySQLPool == nil {
		err := global.InitDBOMySQLPool()
		if err != nil {
			panic(err)
		}
	}
}

func testInitServerID(hostIP string, portNum int) int {
	ipList := strings.Split(hostIP, constant.DotString)
	serverIDStr := fmt.Sprintf(parameter.DefaultServerIDTemplate, portNum, ipList[constant.TwoInt], ipList[constant.ThreeInt])
	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		panic(err)
	}

	return serverID
}

func testInitMySQLServer(hostIP string, portNum, serverID int) *parameter.MySQLServer {
	return parameter.NewMySQLServer(
		testVersion,
		hostIP,
		portNum,
		testRootPaas,
		testAdminUser,
		testAdminPass,
		testClientUser,
		testClientPass,
		testMySQLDMultiUser,
		testMySQLDMultiPass,
		testReplicationUser,
		testReplicationPass,
		testMonitorUser,
		testMonitorPass,
		testDASUser,
		testDASPass,
		defaultTitle,
		testBinaryDirBase,
		testDataDirBaseName,
		testLogDirBaseName,
		testSemiSyncSourceEnabled,
		testSemiSyncReplicaEnabled,
		testSemiSyncSourceTimeout,
		testGroupReplicationConsistency,
		testGroupReplicationFlowControlMode,
		testGroupReplicationMemberWeight,
		serverID,
		testBinlogExpireLogsSeconds,
		testBinlogExpireLogsDays,
		testBackupDir,
		testMaxConnections,
		testInnodbBufferPoolSize,
		testInnodbIOCapacity,
	)
}

func testInitEngine() *Engine {
	return NewEngineWithDefault(testMySQLVersion, testMode, testAddrs, testMySQLServer, testPMMClient)
}

func testInitInstance(addrs []string) error {
	err := linux.SortAddrs(addrs)
	if err != nil {
		return err
	}

	for i, addr := range addrs {
		var (
			isSource   bool
			hostIP     string
			portNumStr string
			portNum    int
		)

		if i == constant.ZeroInt {
			isSource = true
		}

		hostIP, portNumStr, err = net.SplitHostPort(addr)
		if err != nil {
			return errors.Trace(err)
		}
		if hostIP == constant.EmptyString || portNumStr == constant.EmptyString {
			return errors.Errorf("invalid addr: %s", addr)
		}
		portNum, err = strconv.Atoi(portNumStr)
		if err != nil {
			return errors.Trace(err)
		}
		// set MySQL Sever Parameter
		testServerID = testInitServerID(hostIP, portNum)
		testConn = testInitSSHConn(hostIP)
		testMySQLServer = testInitMySQLServer(hostIP, portNum, testServerID)
		err = testMySQLServer.InitWithHostInfo(hostIP, portNum, isSource)
		if err != nil {
			return err
		}

		// clear previous mysql server
		err = testClearMySQL()
		if err != nil {
			return err
		}

		// init os
		err = testEngine.InitOS()
		if err != nil {
			return err
		}
		// init mysql instance
		err = testEngine.InitMySQLInstance()
		if err != nil {
			return err
		}
	}

	return nil
}

func TestEngine_All(t *testing.T) {
	TestEngine_InitOSExecutor(t)
	TestEngine_Install(t)
	TestEngine_InstallSingeInstance(t)
	TestEngine_InitMySQLInstance(t)
	TestEngine_ConfigureReplication(t)
	TestEngine_InitPMMClient(t)
	TestEngine_ConfigureGroupReplication(t)
}

func TestEngine_InitOSExecutor(t *testing.T) {
	asst := assert.New(t)

	err := testEngine.InitOSExecutor()
	asst.Nil(err, "test InitOSExecutor() failed")
	asst.NotNil(testEngine.ose.Conn, "test InitOSExecutor() failed")
	output, err := testEngine.ose.Conn.ExecuteCommand(testDateCommand)
	asst.Nil(err, "test InitOSExecutor() failed")
	t.Logf("output: %s", output)
}

func TestEngine_Install(t *testing.T) {
	asst := assert.New(t)

	// clear previous mysql server
	err := testClearMySQL(testAddrs...)
	asst.Nil(err, "test InstallSingleInstance() failed")
	// install single instance
	err = testEngine.Install(1)
	asst.Nil(err, "test InstallSingleInstance() failed")
	// clear mysql server
	// err = testClearMySQL()
	asst.Nil(err, "test InstallSingleInstance() failed")
}

func TestEngine_InstallSingeInstance(t *testing.T) {
	asst := assert.New(t)

	// clear previous mysql server
	err := testClearMySQL()
	asst.Nil(err, "test InstallSingleInstance() failed")
	// install single instance
	err = testEngine.InstallSingleInstance(testHostIP1, testPortNum1, true)
	asst.Nil(err, "test InstallSingleInstance() failed")
	// clear mysql server
	//err = testClearMySQL()
	asst.Nil(err, "test InstallSingleInstance() failed")
}

func TestEngine_InitMySQLInstance(t *testing.T) {
	asst := assert.New(t)

	// init os executor
	err := testEngine.InitOSExecutor()
	asst.Nil(err, "test InitOSExecutor() failed")
	asst.NotNil(testEngine.ose.Conn, "test InitOSExecutor() failed")
	// clear previous mysql server
	err = testClearMySQL()
	asst.Nil(err, "test InitMySQLInstance() failed")
	// init os
	err = testEngine.InitOS()
	asst.Nil(err, "test InitMySQLInstance() failed")
	// init mysql instance
	err = testEngine.InitMySQLInstance()
	asst.Nil(err, "test InitMySQLInstance() failed")
	// create connection
	conn, err := mysql.NewConn(testAddr1, constant.EmptyString, testClientUser, testClientPass)
	asst.Nil(err, "test InitMySQLInstance() failed")
	defer func() {
		err = conn.Close()
		asst.Nil(err, "test InitMySQLInstance() failed")
	}()
	// check version
	ok := conn.CheckInstanceStatus()
	asst.True(ok, "test InitMySQLInstance() failed")
	// clear mysql server
	err = testClearMySQL()
	asst.Nil(err, "test InitMySQLInstance() failed")
}

func TestEngine_ConfigureReplication(t *testing.T) {
	asst := assert.New(t)

	err := testInitInstance(testAddrs)
	asst.Nil(err, "test ConfigureReplica() failed")

	err = testEngine.ConfigureReplica(testAddr2, testHostIP1, testPortNum1)
	asst.Nil(err, "test ConfigureReplica() failed")
	// create connection
	conn, err := mysql.NewConn(testAddr2, constant.EmptyString, testClientUser, testClientPass)
	asst.Nil(err, "test ConfigureReplica() failed")
	defer func() {
		err = conn.Close()
		asst.Nil(err, "test ConfigureReplica() failed")
	}()
	// check version
	result, err := conn.Conn.Execute(testCheckSlaveStatusSQL)
	asst.Nil(err, "test ConfigureReplica() failed")
	asst.True(result.RowNumber() == 1, "test ConfigureReplica() failed")
}

func TestEngine_InitPMMClient(t *testing.T) {

}

func TestEngine_ConfigureGroupReplication(t *testing.T) {

}
