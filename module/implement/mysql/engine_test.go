package mysql

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"github.com/romberli/db-operator/module/implement/mysql/parameter"
	"github.com/romberli/db-operator/pkg/util/ssh"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testMySQLInstallationPackageDir = "/data/software/mysql"

	testHostIP          = "192.168.137.21"
	testPortNum         = 3306
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
	testServerID                        = 3306137011
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
	testPMMClientVersion   = "2.24.0"
	testPMMServiceName     = "192-168-137-11-3306"
	testReplicationSetName = constant.EmptyString
)

var (
	testAddr  = fmt.Sprintf("%s:%d", testHostIP, testPortNum)
	testAddrs = []string{testAddr}

	testOSVersion    = version.Must(version.NewVersion(testOSVersionStr))
	testMySQLVersion = version.Must(version.NewVersion(testMySQLVersionStr))

	testMySQLServer *parameter.MySQLServer
	testPMMClient   *parameter.PMMClient

	testEngine *Engine
)

func init() {
	testInitViper()
	testMySQLServer = testInitMySQLServer()
	testPMMClient = testInitPMMClient()
	testSSHConn = testInitSSHConn()

	testEngine = testInitEngine()
}

func testInitViper() {
	viper.Set(config.MySQLInstallationPackageDirKey, testMySQLInstallationPackageDir)
	viper.Set(config.PMMServerAddrKey, testPMMServerAddr)
	viper.Set(config.PMMClientVersionKey, testPMMClientVersion)
}

func testInitMySQLServer() *parameter.MySQLServer {
	return parameter.NewMySQLServer(
		testVersion,
		testHostIP,
		testPortNum,
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
		testServerID,
		testBinlogExpireLogsSeconds,
		testBinlogExpireLogsDays,
		testBackupDir,
		testMaxConnections,
		testInnodbBufferPoolSize,
		testInnodbIOCapacity,
	)
}

func testInitPMMClient() *parameter.PMMClient {
	return parameter.NewPMMClientWithDefault()
}

func testInitSSHConn() *ssh.Conn {
	sshConn, err := linux.NewSSHConn(testHostIP, constant.DefaultSSHPort, testOSUser, testOSPass, testUseSudo)
	if err != nil {
		panic(err)
	}

	return ssh.NewConn(sshConn)
}

func testInitEngine() *Engine {
	return NewEngineWithDefault(testSSHConn, testMode, testAddrs, testMySQLServer, testPMMClient)
}

func TestEngine_All(t *testing.T) {
	TestEngine_InitMySQLInstance(t)
	TestEngine_ConfigureMySQLCluster(t)
	TestEngine_InitPMMClient(t)
	TestEngine_Install(t)
}

func TestEngine_InitMySQLInstance(t *testing.T) {
	asst := assert.New(t)

	err := testEngine.InitMySQLInstance(true)
	asst.Nil(err, "test InitMySQLInstance() failed")
	// create connection
	conn, err := mysql.NewConn(testAddr, constant.EmptyString, testClientUser, testClientPass)
	asst.Nil(err, "test InitMySQLInstance() failed")
	defer func() {
		err = conn.Close()
		asst.Nil(err, "test InitMySQLInstance() failed")
	}()
	// check version
	ok := conn.CheckInstanceStatus()
	asst.True(ok, "test InitMySQLInstance() failed")
}

func TestEngine_ConfigureMySQLCluster(t *testing.T) {
	asst := assert.New(t)

	err := testEngine.ConfigureMySQLCluster(testAddr)
	asst.Nil(err, "test ConfigureMySQLCluster() failed")
	// create connection
	conn, err := mysql.NewConn(testAddr, constant.EmptyString, testClientUser, testClientPass)
	asst.Nil(err, "test ConfigureMySQLCluster() failed")
	defer func() {
		err = conn.Close()
		asst.Nil(err, "test ConfigureMySQLCluster() failed")
	}()
	// check version
	ok := conn.CheckInstanceStatus()
	asst.True(ok, "test ConfigureMySQLCluster() failed")
}

func TestEngine_InitPMMClient(t *testing.T) {

}

func TestEngine_Install(t *testing.T) {

}
