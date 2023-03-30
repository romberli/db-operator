package mysql

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/romberli/db-operator/pkg/util/ssh"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/stretchr/testify/assert"
)

const (
	testMySQLVersionCommand = `export PATH=\$PATH:%s/bin && mysql --version`
)

var (
	testConn       *ssh.Conn
	testOSExecutor *OSExecutor
)

func init() {
	testInitViper()
	testServerID = testInitServerID(testHostIP1, testPortNum1)
	testMySQLServer = testInitMySQLServer(testHostIP1, testPortNum1, testServerID)
	testPMMClient = testInitPMMClient()
	testConn = testInitSSHConn(testHostIP1)

	testOSExecutor = testInitOSExecutor()
}

func testInitSSHConn(hostIP string) *ssh.Conn {
	sshConn, err := linux.NewSSHConn(hostIP, constant.DefaultSSHPort, testOSUser, testOSPass, testUseSudo)
	if err != nil {
		panic(err)
	}

	return ssh.NewConn(sshConn)
}

func testInitOSExecutor() *OSExecutor {
	ose := NewOSExecutor(testConn, testMySQLVersion, testMySQLServer)
	err := ose.InitExecutor()
	if err != nil {
		panic(err)
	}

	return ose
}

func testClearMySQL(addrs ...string) error {
	for _, addr := range addrs {
		hostIP, portNumStr, err := net.SplitHostPort(addr)
		if err != nil {
			return err
		}

		portNum, err := strconv.Atoi(portNumStr)
		if err != nil {
			return err
		}

		testConn = testInitSSHConn(testHostIP1)
		testOSExecutor = testInitOSExecutor()
		err = testOSExecutor.mysqlServer.InitWithHostInfo(hostIP, portNum, false)
		if err != nil {
			return err
		}

		// kill mysql server process
		pidList, err := testOSExecutor.GetMySQLPIDList()
		if err != nil {
			return err
		}
		for _, pid := range pidList {
			cmd := fmt.Sprintf(testKillMySQLDCommand, pid)
			err = testOSExecutor.Conn.ExecuteCommandWithoutOutput(cmd)
			if err != nil {
				return err
			}
		}

		// remove the mysql server files
		err = testOSExecutor.Conn.RemoveAll(testOSExecutor.mysqlServer.BinaryDirBase)
		if err != nil {
			return err
		}
		err = testOSExecutor.Conn.RemoveAll(testOSExecutor.mysqlServer.DataDirBase)
		if err != nil {
			return err
		}
		err = testOSExecutor.Conn.RemoveAll(testOSExecutor.mysqlServer.LogDirBase)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestOSExecutor_All(t *testing.T) {
	TestOSExecutor_Precheck(t)
	TestOSExecutor_GetMySQLPIDList(t)
	TestOSExecutor_InstallRPM(t)
	TestOSExecutor_InitUserAndGroup(t)
	TestOSExecutor_InitDir(t)
	TestOSExecutor_InstallMySQLBinary(t)
}

func TestOSExecutor_Precheck(t *testing.T) {
	asst := assert.New(t)

	err := testClearMySQL()
	asst.Nil(err, "test Precheck() failed")
	err = testOSExecutor.Precheck()
	asst.Nil(err, "test Precheck() failed")
}

func TestOSExecutor_GetMySQLPIDList(t *testing.T) {
	asst := assert.New(t)

	pidList, err := testOSExecutor.GetMySQLPIDList()
	asst.Nil(err, "test GetMySQLPIDList() failed")
	asst.Equal(constant.ZeroInt, len(pidList), "test GetMySQLPIDList() failed")
}

func TestOSExecutor_InstallRPM(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.InstallRPM()
	asst.Nil(err, "test InstallRPM() failed")
	pathExists, err := testOSExecutor.Conn.PathExists(libNCursesPath)
	asst.Nil(err, "test InstallRPM() failed")
	asst.True(pathExists, "test InstallRPM() failed")
	pathExists, err = testOSExecutor.Conn.PathExists(libTInfoPath)
	asst.Nil(err, "test InstallRPM() failed")
	asst.True(pathExists, "test InstallRPM() failed")
}

func TestOSExecutor_InitUserAndGroup(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.InitUserAndGroup()
	asst.Nil(err, "test InitUserAndGroup() failed")
	err = testOSExecutor.Conn.ExecuteCommandWithoutOutput(checkMySQLGroupCommand)
	asst.Nil(err, "test InitUserAndGroup() failed")
	err = testOSExecutor.Conn.ExecuteCommandWithoutOutput(checkMySQLUserCommand)
	asst.Nil(err, "test InitUserAndGroup() failed")
}

func TestOSExecutor_InitDir(t *testing.T) {
	asst := assert.New(t)

	err := testClearMySQL()
	asst.Nil(err, "test Precheck() failed")
	err = testOSExecutor.InitDir()
	asst.Nil(err, "test InitDir() failed")
	pathExists, err := testOSExecutor.Conn.PathExists(testOSExecutor.mysqlServer.BackupDir)
	asst.Nil(err, "test InitDir() failed")
	asst.True(pathExists, "test InitDir() failed")
	pathExists, err = testOSExecutor.Conn.PathExists(testOSExecutor.mysqlServer.DataDirBase)
	asst.Nil(err, "test InitDir() failed")
	asst.True(pathExists, "test InitDir() failed")
	pathExists, err = testOSExecutor.Conn.PathExists(testOSExecutor.mysqlServer.LogDirBase)
	asst.Nil(err, "test InitDir() failed")
	asst.True(pathExists, "test InitDir() failed")
}

func TestOSExecutor_InstallMySQLBinary(t *testing.T) {
	asst := assert.New(t)

	err := testClearMySQL()
	asst.Nil(err, "test Precheck() failed")
	err = testOSExecutor.InstallMySQLBinary()
	asst.Nil(err, "test InstallMySQLBinary() failed")
	pathExists, err := testOSExecutor.Conn.PathExists(testOSExecutor.mysqlServer.BinaryDirBase)
	asst.Nil(err, "test InstallMySQLBinary() failed")
	asst.True(pathExists, "test InstallMySQLBinary() failed")
}

func TestOSExecutor_ConfigurePathEnv(t *testing.T) {
	asst := assert.New(t)

	err := testOSExecutor.ConfigurePathEnv()
	asst.Nil(err, "test ConfigurePathEnv() failed")
	cmd := fmt.Sprintf(testMySQLVersionCommand, testOSExecutor.mysqlServer.BinaryDirBase)
	output, err := testOSExecutor.Conn.ExecuteCommand(cmd)
	asst.Nil(err, "test ConfigurePathEnv() failed")
	t.Log(output)
}

func TestOSExecutor_Init(t *testing.T) {
	asst := assert.New(t)

	err := testClearMySQL()
	asst.Nil(err, "test Precheck() failed")
	err = testOSExecutor.Init()
	asst.Nil(err, "test Init() failed")
	pathExists, err := testOSExecutor.Conn.PathExists(testOSExecutor.mysqlServer.BinaryDirBase)
	asst.Nil(err, "test InstallMySQLBinary() failed")
	asst.True(pathExists, "test InstallMySQLBinary() failed")
}
