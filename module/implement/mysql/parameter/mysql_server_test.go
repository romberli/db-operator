package parameter

import (
	"testing"

	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"github.com/romberli/go-util/common"
	"github.com/stretchr/testify/assert"
)

const (
	testConfigPath = "/tmp/test.cnf"
)

var (
	testMySQLServer *MySQLServer
)

func init() {
	initTestMySQLServer()
}

func initTestMySQLServer() {
	testMySQLServer = newMySQLServer(
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
		testInnodbIOCapacityMax,
	)
}

func TestMySQLServer_All(t *testing.T) {
	TestMySQLServer_GetConfig(t)
	TestMySQLServer_WriteConfig(t)
	TestMySQLServer_Marshal(t)
	TestMySQLServer_Unmarshal(t)
}

func TestMySQLServer_GetConfig(t *testing.T) {
	asst := assert.New(t)

	config, err := testMySQLServer.GetConfig(testMySQLVersion, mode.AsyncReplication)
	asst.Nil(err, common.CombineMessageWithError("test GetConfig() failed", err))
	t.Log(string(config))
}

func TestMySQLServer_WriteConfig(t *testing.T) {
	asst := assert.New(t)

	config, err := testMySQLServer.GetConfig(testMySQLVersion, mode.GroupReplication)
	asst.Nil(err, common.CombineMessageWithError("test WriteConfig() failed", err))
	err = testMySQLServer.WriteConfig(testConfigPath, config)
	asst.Nil(err, common.CombineMessageWithError("test WriteConfig() failed", err))
	t.Log(string(config))
}

func TestMySQLServer_Marshal(t *testing.T) {
	asst := assert.New(t)

	jsonBytes, err := testMySQLServer.Marshal()
	asst.Nil(err, common.CombineMessageWithError("test Marshal() failed", err))
	t.Log(string(jsonBytes))
}

func TestMySQLServer_Unmarshal(t *testing.T) {
	asst := assert.New(t)

	jsonBytes := []byte(`{"host_ip": "192.168.123.1", "ddd": 1}`)
	err := testMySQLServer.Unmarshal(jsonBytes)
	asst.Nil(err, common.CombineMessageWithError("test Unmarshal() failed", err))
	t.Logf("host_ip: %s, port_num: %d, server_id: %d", testMySQLServer.HostIP, testMySQLServer.PortNum, testMySQLServer.ServerID)
}
