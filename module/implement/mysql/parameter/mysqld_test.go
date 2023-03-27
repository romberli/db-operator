package parameter

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"github.com/romberli/go-util/common"
	"github.com/stretchr/testify/assert"
)

var (
	testMySQLD       *MySQLD
	testMySQLVersion *version.Version
)

func init() {
	initTestMySQLD()
	initTestVersion()
}

func initTestMySQLD() {
	testMySQLD = NewMySQLD(
		testVersion,
		testHostIP,
		testPortNum,
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

func initTestVersion() *version.Version {
	var err error

	testMySQLVersion, err = version.NewVersion(mysql80)
	if err != nil {
		panic(common.CombineMessageWithError("init test version failed", err))
	}

	return testMySQLVersion
}

func TestMySQLD_All(t *testing.T) {
	TestMySQLD_GetConfig(t)
}

func TestMySQLD_GetConfig(t *testing.T) {
	asst := assert.New(t)

	config, err := testMySQLD.GetConfig(testMySQLVersion, mode.AsyncReplication)
	asst.Nil(err, common.CombineMessageWithError("test TestMySQLD_GetConfig() failed", err))
	t.Log(string(config))
}
