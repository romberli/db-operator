package mysql

import (
	"github.com/romberli/db-operator/pkg/util/ssh"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testSSHConn    *ssh.Conn
	testOSExecutor *OSExecutor
)

func init() {
	testInitViper()
	testMySQLServer = testInitMySQLServer()
	testPMMClient = testInitPMMClient()
	testSSHConn = testInitSSHConn()

	testOSExecutor = testInitOSExecutor()
}

func testInitOSExecutor() *OSExecutor {
	ose := NewOSExecutor(testMySQLVersion, testSSHConn, testMySQLServer)
	err := ose.InitExecutor()
	if err != nil {
		panic(err)
	}

	return ose
}

func TestOSExecutor_All(t *testing.T) {
	TestOSExecutor_Precheck(t)
	TestOSExecutor_InstallRPM(t)
	TestOSExecutor_InitUserAndGroup(t)
	TestOSExecutor_InitDir(t)
	TestOSExecutor_InstallMySQLBinary(t)
}

func TestOSExecutor_Precheck(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.Precheck()
	asst.Nil(err, "test Precheck() failed")
}

func TestOSExecutor_InstallRPM(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.InstallRPM()
	asst.Nil(err, "test InstallRPM() failed")
	pathExists, err := testOSExecutor.sshConn.PathExists(libNCursesPath)
	asst.Nil(err, "test InstallRPM() failed")
	asst.True(pathExists, "test InstallRPM() failed")
	pathExists, err = testOSExecutor.sshConn.PathExists(libTInfoPath)
	asst.Nil(err, "test InstallRPM() failed")
	asst.True(pathExists, "test InstallRPM() failed")
}

func TestOSExecutor_InitUserAndGroup(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.InitUserAndGroup()
	asst.Nil(err, "test InitUserAndGroup() failed")
	err = testOSExecutor.sshConn.ExecuteCommandWithoutOutput(checkMySQLGroupCommand)
	asst.Nil(err, "test InitUserAndGroup() failed")
	err = testOSExecutor.sshConn.ExecuteCommandWithoutOutput(checkMySQLUserCommand)
	asst.Nil(err, "test InitUserAndGroup() failed")
}

func TestOSExecutor_InitDir(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.InitDir()
	asst.Nil(err, "test InitDir() failed")
	pathExists, err := testOSExecutor.sshConn.PathExists(testOSExecutor.mysqlServer.BackupDir)
	asst.Nil(err, "test InitDir() failed")
	asst.True(pathExists, "test InitDir() failed")
	pathExists, err = testOSExecutor.sshConn.PathExists(testOSExecutor.mysqlServer.DataDirBase)
	asst.Nil(err, "test InitDir() failed")
	asst.True(pathExists, "test InitDir() failed")
	pathExists, err = testOSExecutor.sshConn.PathExists(testOSExecutor.mysqlServer.LogDirBase)
	asst.Nil(err, "test InitDir() failed")
	asst.True(pathExists, "test InitDir() failed")
}

func TestOSExecutor_InstallMySQLBinary(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.InstallMySQLBinary()
	asst.Nil(err, "test InstallMySQLBinary() failed")
	pathExists, err := testOSExecutor.sshConn.PathExists(testOSExecutor.mysqlServer.BinaryDirBase)
	asst.Nil(err, "test InstallMySQLBinary() failed")
	asst.True(pathExists, "test InstallMySQLBinary() failed")
}

func TestOSExecutor_Init(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.Init()
	asst.Nil(err, "test Init() failed")
	pathExists, err := testOSExecutor.sshConn.PathExists(testOSExecutor.mysqlServer.BinaryDirBase)
	asst.Nil(err, "test InstallMySQLBinary() failed")
	asst.True(pathExists, "test InstallMySQLBinary() failed")
}
