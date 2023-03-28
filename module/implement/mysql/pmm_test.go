package mysql

import (
	"testing"

	"github.com/romberli/db-operator/module/implement/mysql/parameter"
	"github.com/romberli/go-util/constant"
	"github.com/stretchr/testify/assert"
)

var (
	testPMMClient   *parameter.PMMClient
	testPMMExecutor *PMMExecutor
)

func init() {
	testInitViper()
	testConn = testInitSSHConn(testHostIP1)
	testPMMClient = testInitPMMClient()
	testPMMExecutor = testInitPMMExecutor(testHostIP1, testPortNum1)
}

func testInitPMMClient() *parameter.PMMClient {
	return parameter.NewPMMClientWithDefault()
}

func testInitPMMExecutor(hostIP string, portNum int) *PMMExecutor {
	return NewPMMExecutor(testConn, hostIP, portNum, testPMMClient)
}

func TestPMMExecutor_All(t *testing.T) {
	TestPMMExecutor_CheckPMMClient(t)
	TestPMMExecutor_Install(t)
	TestPMMExecutor_CheckConfiguration(t)
	TestPMMExecutor_ConfigureServer(t)
	TestPMMExecutor_CheckServiceExists(t)
	TestPMMExecutor_AddService(t)
}

func TestPMMExecutor_CheckPMMClient(t *testing.T) {
	asst := assert.New(t)

	ok, err := testPMMExecutor.CheckPMMClient()
	asst.Nil(err, "test CheckPMMClient() failed")
	asst.True(ok, "test CheckPMMClient() failed")
}

func TestPMMExecutor_Install(t *testing.T) {
	asst := assert.New(t)

	arch, err := testPMMExecutor.sshConn.GetArch()
	asst.Nil(err, "test Install() failed")
	if arch == constant.X64Arch {
		err = testPMMExecutor.Install()
		asst.Nil(err, "test Install() failed")
	} else {
		t.Skip("skip test Install() for non-x64 arch")
	}
}

func TestPMMExecutor_CheckConfiguration(t *testing.T) {
	asst := assert.New(t)

	configured, err := testPMMExecutor.CheckConfiguration()
	asst.Nil(err, "test CheckConfiguration() failed")
	asst.True(configured, "test CheckConfiguration() failed")
}

func TestPMMExecutor_ConfigureServer(t *testing.T) {
	asst := assert.New(t)

	configured, err := testPMMExecutor.CheckConfiguration()
	asst.Nil(err, "test ConfigureServer() failed")
	if !configured {
		err := testPMMExecutor.ConfigureServer()
		asst.Nil(err, "test ConfigureServer() failed")
		err = testPMMExecutor.StartClient()
		asst.Nil(err, "test ConfigureServer() failed")
	} else {
		t.Skip("skip test ConfigureServer() for configured server")
	}
}

func TestPMMExecutor_CheckServiceExists(t *testing.T) {
	asst := assert.New(t)

	ok, err := testPMMExecutor.CheckServiceExists()
	asst.Nil(err, "test CheckServiceExists() failed")
	asst.True(ok, "test CheckServiceExists() failed")
}

func TestPMMExecutor_AddService(t *testing.T) {
	asst := assert.New(t)

	ok, err := testPMMExecutor.CheckServiceExists()
	asst.Nil(err, "test AddService() failed")
	if !ok {
		err = testPMMExecutor.AddService()
		asst.Nil(err, "test AddService() failed")
		ok, err = testPMMExecutor.CheckServiceExists()
		asst.Nil(err, "test CheckServiceExists() failed")
		asst.True(ok, "test CheckServiceExists() failed")
	} else {
		t.Skip("skip test AddService() for existing service")
	}
}
