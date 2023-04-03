package mysql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testService *Service
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
	testService = testInitService()
}

func testInitService() *Service {
	return NewServiceWithDefault(testEngine)
}

func TestService_Install(t *testing.T) {
	asst := assert.New(t)

	// clear previous mysql server
	err := testClearMySQL(testAddrs...)
	asst.Nil(err, "test InstallSingleInstance() failed")
	// install
	err = testService.Install()
	asst.Nil(err, "test Install() failed")
	// clear previous mysql server
	// err = testClearMySQL(testAddrs...)
	asst.Nil(err, "test InstallSingleInstance() failed")
}
