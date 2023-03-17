package parameter

import (
	"testing"

	"github.com/romberli/go-util/common"
	"github.com/stretchr/testify/assert"
)

const (
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
	testReplicationUser = "replication"
	testReplicationPass = "replication"
	testMonitorUser     = "pmm"
	testMonitorPass     = "pmm"
	testDASUser         = "das"
	testDASPass         = "das"

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
	testInnodbIOCapacityMax             = 2000
)

var (
	testClient *Common
)

func init() {
	initTestClient()
}

func initTestClient() {
	testClient = NewCommon(testPortNum, testDataDirBaseName, testClientUser, testClientPass, testMySQLDMultiUser, testMySQLDMultiPass)
}

func TestClient_All(t *testing.T) {
	TestClient_GetConfig(t)
}

func TestClient_GetConfig(t *testing.T) {
	asst := assert.New(t)

	config, err := testClient.GetConfig()
	asst.Nil(err, common.CombineMessageWithError("test GetConfig() failed", err))
	t.Log(string(config))
}
