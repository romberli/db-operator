package mysql

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/spf13/viper"

	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/module/implement/mysql/parameter"
	"github.com/romberli/db-operator/pkg/util/ssh"
)

const (
	os9VersionStr           = "9.0.0"
	mysql8026VersionStr     = "8.0.26"
	mysql8032VersionStr     = "8.0.32"
	minX64MySQLVersionStr   = mysql8026VersionStr
	minAArchMySQLVersionStr = mysql8032VersionStr
	// mysql 8.0
	mysqlServerBinaryPackageNameTemplateV1 = "mysql-%s-linux-glibc2.12-x86_64.tar.xz"
	mysqlServerBinaryPackageNameTemplateV2 = "mysql-%s-linux-glibc2.17-aarch64.tar.gz"
	decompressCommandV1                    = "tar xf %s -C %s"
	decompressCommandV2                    = "tar zxf %s -C %s"
	tarGZExt                               = ".tar.gz"
	tarXZExt                               = ".tar.xz"

	dataDirName     = "data"
	logDirName      = "log"
	binlogDirName   = "binlog"
	relaylogDirName = "relaylog"
	tmpDirName      = "tmp"
	runDirName      = "run"

	defaultMySQLUser  = "mysql"
	defaultMySQLGroup = "mysql"

	bashProfilePath = "/root/.bash_profile"

	yumInstallCommand       = "/usr/bin/yum install -y ncurses-c++-libs ncurses-libs"
	libNCursesPath          = "/usr/lib64/libncurses.so.5"
	libTInfoPath            = "/usr/lib64/libtinfo.so.5"
	lnLibNCursesCommand     = "/usr/bin/ln -s /usr/lib64/libncurses.so.6.2 /usr/lib64/libncurses.so.5"
	lnLibTInfoCommand       = "/usr/bin/ln -s /usr/lib64/libtinfo.so.6.2 /usr/lib64/libtinfo.so.5"
	checkMySQLGroupCommand  = "/usr/bin/id -g mysql"
	checkMySQLUserCommand   = "/usr/bin/id -u mysql"
	createMySQLGroupCommand = "/usr/sbin/groupadd -g 1001 mysql"
	createMySQLUserCommand  = "/usr/sbin/useradd -u 1001 -g mysql mysql"
	checkPathEnvCommand     = "/usr/bin/grep PATH %s | /usr/bin/grep -c %s/bin | /usr/bin/grep -v grep"
	addPathEnvCommand       = `/usr/bin/echo 'export PATH=\$PATH:%s/bin' >> %s`
)

var (
	minAArchMySQLVersion = version.Must(version.NewVersion(minAArchMySQLVersionStr))
	minX64MySQLVersion   = version.Must(version.NewVersion(minX64MySQLVersionStr))
	os9Version           = version.Must(version.NewVersion(os9VersionStr))
)

type OSExecutor struct {
	*ssh.Conn
	mysqlVersion *version.Version
	mysqlServer  *parameter.MySQLServer

	arch      string
	osVersion *version.Version
}

// NewOSExecutor returns a new *OSExecutor
func NewOSExecutor(sshConn *ssh.Conn, mysqlVersion *version.Version, mysqlServer *parameter.MySQLServer) *OSExecutor {
	return newOSExecutor(sshConn, mysqlVersion, mysqlServer)
}

// newOSExecutor returns a new *OSExecutor
func newOSExecutor(sshConn *ssh.Conn, mysqlVersion *version.Version, mysqlServer *parameter.MySQLServer) *OSExecutor {
	return &OSExecutor{
		Conn:         sshConn,
		mysqlVersion: mysqlVersion,
		mysqlServer:  mysqlServer,
	}
}

// Init initializes the os
func (ose *OSExecutor) Init() error {
	// init executor
	err := ose.InitExecutor()
	// precheck
	err = ose.Precheck()
	if err != nil {
		return err
	}
	// Install rpm
	err = ose.InstallRPM()
	if err != nil {
		return err
	}
	// init user and group
	err = ose.InitUserAndGroup()
	if err != nil {
		return err
	}
	// init dir
	err = ose.InitDir()
	if err != nil {
		return err
	}
	// Install mysql binary
	err = ose.InstallMySQLBinary()
	if err != nil {
		return err
	}
	// Configure path env
	err = ose.ConfigurePathEnv()
	if err != nil {
		return err
	}

	return nil
}

// InitExecutor initializes the os executor
func (ose *OSExecutor) InitExecutor() error {
	var err error
	// get os version
	ose.osVersion, err = ose.Conn.GetOSVersion()
	if err != nil {
		return err
	}
	// get arch
	ose.arch, err = ose.Conn.GetArch()
	if err != nil {
		return err
	}

	return nil
}

// Precheck checks the os
func (ose *OSExecutor) Precheck() error {
	// check minimum version
	if (ose.arch == constant.AArch64Arch && ose.mysqlVersion.LessThan(minAArchMySQLVersion)) ||
		(ose.arch == constant.X64Arch && ose.mysqlVersion.LessThan(minX64MySQLVersion)) {
		return errors.Errorf("the minimum mysql version on %s is %s, %s not valid", ose.arch, minAArchMySQLVersion.String(), ose.mysqlVersion.String())
	}
	// check if mysql pid exists
	pidList, err := ose.GetMySQLPIDList()
	if err != nil {
		return err
	}
	if len(pidList) > constant.ZeroInt {
		return errors.Errorf("mysql pid exists, installation aborted. pid list: %v", pidList)
	}

	// check if mysql installation package exists
	installationPackagePath := filepath.Join(viper.GetString(config.MySQLInstallationPackageDirKey), ose.getMySQLServerBinaryPackageName())
	exists, err := linux.PathExists(installationPackagePath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("mysql installation package does not exist. installation package path: %s", installationPackagePath)
	}
	// check if the mysql data directory exists
	dataDir := filepath.Join(ose.mysqlServer.DataDirBase, dataDirName)
	output, err := ose.Conn.ListPath(dataDir)
	if err == nil && len(output) > constant.ZeroInt {
		return errors.Errorf("mysql data directory exists and is not empty, installation aborted. data directory: %s", dataDir)
	}
	// check if the mysql binlog directory exists
	binlogDir := filepath.Join(ose.mysqlServer.LogDirBase, binlogDirName)
	output, err = ose.Conn.ListPath(binlogDir)
	if err == nil && len(output) > constant.ZeroInt {
		return errors.Errorf("mysql binlog directory exists and is not empty, installation aborted. binlog directory: %s", dataDir)
	}
	// check if the mysql relaylog directory exists
	relaylogDir := filepath.Join(ose.mysqlServer.LogDirBase, relaylogDirName)
	output, err = ose.Conn.ListPath(relaylogDir)
	if err == nil && len(output) > constant.ZeroInt {
		return errors.Errorf("mysql relaylog directory exists and is not empty, installation aborted. relaylog directory: %s", dataDir)
	}

	return nil
}

// GetMySQLPIDList gets the mysql pid list
func (ose *OSExecutor) GetMySQLPIDList() ([]int, error) {
	cmd := fmt.Sprintf(getMySQLPIDListCommandTemplate, ose.mysqlServer.PortNum, ose.mysqlServer.DataDirBase)
	output, err := ose.Conn.ExecuteCommand(cmd)
	if err != nil {
		return nil, err
	}

	if output == constant.EmptyString {
		return nil, nil
	}

	pidStrList := strings.Split(output, constant.CRLFString)
	pidList := make([]int, len(pidStrList))
	for i, pidStr := range pidStrList {
		pid, err := strconv.Atoi(strings.TrimSpace(pidStr))
		if err != nil {
			return nil, errors.Trace(err)
		}
		pidList[i] = pid
	}

	return pidList, nil
}

// InstallRPM installs the rpm
func (ose *OSExecutor) InstallRPM() error {
	err := ose.Conn.ExecuteCommandWithoutOutput(yumInstallCommand)
	if err != nil {
		return err
	}

	if ose.osVersion.GreaterThanOrEqual(os9Version) {
		pathExists, err := ose.Conn.PathExists(libNCursesPath)
		if err != nil {
			return err
		}
		if !pathExists {
			err = ose.Conn.ExecuteCommandWithoutOutput(lnLibNCursesCommand)
			if err != nil {
				return err
			}
		}
		pathExists, err = ose.Conn.PathExists(libTInfoPath)
		if err != nil {
			return err
		}
		if !pathExists {
			err = ose.Conn.ExecuteCommandWithoutOutput(lnLibTInfoCommand)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// InitUserAndGroup initializes the user and group
func (ose *OSExecutor) InitUserAndGroup() error {
	// init mysql group
	err := ose.Conn.ExecuteCommandWithoutOutput(checkMySQLGroupCommand)
	if err != nil {
		err = ose.Conn.ExecuteCommandWithoutOutput(createMySQLGroupCommand)
		if err != nil {
			return errors.Trace(err)
		}
	}
	// init mysql user
	err = ose.Conn.ExecuteCommandWithoutOutput(checkMySQLUserCommand)
	if err != nil {
		err = ose.Conn.ExecuteCommandWithoutOutput(createMySQLUserCommand)
		if err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

// InitDir initializes the directory
func (ose *OSExecutor) InitDir() error {
	// create directories
	binaryDirParent := filepath.Dir(ose.mysqlServer.BinaryDirBase)
	err := ose.Conn.MkdirAll(binaryDirParent)
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.Conn.MkdirAll(ose.mysqlServer.BackupDir)
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.Conn.MkdirAll(filepath.Join(ose.mysqlServer.DataDirBase, dataDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.Conn.MkdirAll(filepath.Join(ose.mysqlServer.DataDirBase, logDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.Conn.MkdirAll(filepath.Join(ose.mysqlServer.DataDirBase, tmpDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.Conn.MkdirAll(filepath.Join(ose.mysqlServer.DataDirBase, runDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.Conn.MkdirAll(filepath.Join(ose.mysqlServer.LogDirBase, binlogDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.Conn.MkdirAll(filepath.Join(ose.mysqlServer.LogDirBase, relaylogDirName))
	if err != nil {
		return errors.Trace(err)
	}
	// change owner of directories
	err = ose.Conn.Chown(ose.mysqlServer.BackupDir, defaultMySQLUser, defaultMySQLGroup)
	if err != nil {
		return err
	}
	err = ose.Conn.Chown(ose.mysqlServer.DataDirBase, defaultMySQLUser, defaultMySQLGroup)
	if err != nil {
		return err
	}
	err = ose.Conn.Chown(ose.mysqlServer.LogDirBase, defaultMySQLUser, defaultMySQLGroup)
	if err != nil {
		return err
	}

	return nil
}

// InstallMySQLBinary installs the mysql binary
func (ose *OSExecutor) InstallMySQLBinary() error {
	// check mysql binary directory exists
	exists, err := ose.Conn.PathExists(ose.mysqlServer.BinaryDirBase)
	if err != nil {
		return err
	}
	if exists {
		// mysql binary directory exists, maybe just want to add new instance
		return nil
	}
	// copy mysql installation package
	err = ose.copyMySQLServerBinaryPackages()
	if err != nil {
		return err
	}
	// Install mysql binary
	cmd := ose.getDecompressCommand()
	err = ose.Conn.ExecuteCommandWithoutOutput(cmd)
	if err != nil {
		return err
	}

	err = ose.Conn.Move(filepath.Join(
		constant.DefaultTmpDir, ose.getMySQLServerBinaryPackageDecompressedDirName()), ose.mysqlServer.BinaryDirBase)
	if err != nil {
		return err
	}

	return nil
}

// ConfigurePathEnv configures the path environment variable
func (ose *OSExecutor) ConfigurePathEnv() error {
	// check if bash profile exists
	exists, err := ose.Conn.PathExists(bashProfilePath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("ose.ConfigurePathEnv(): file does not exist. file: %s", bashProfilePath)
	}

	// check if mysql binary directory exists in PATH
	cmd := fmt.Sprintf(checkPathEnvCommand, bashProfilePath, ose.mysqlServer.BinaryDirBase)
	output, err := ose.Conn.ExecuteCommand(cmd)
	if err != nil {
		return err
	}

	if output == strconv.Itoa(constant.ZeroInt) {
		// add mysql binary directory to PATH
		cmd = fmt.Sprintf(addPathEnvCommand, ose.mysqlServer.BinaryDirBase, bashProfilePath)
		err = ose.Conn.ExecuteCommandWithoutOutput(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

// copyMySQLServerBinaryPackages copies the mysql server binary packages to the remote host
func (ose *OSExecutor) copyMySQLServerBinaryPackages() error {
	fileName := ose.getMySQLServerBinaryPackageName()
	fileNameSource := ose.getMySQLInstallationPackagePath()
	fileNameDest := filepath.Join(constant.DefaultTmpDir, fileName)
	// copy mysql installation package
	err := ose.Conn.CopySingleFileToRemote(fileNameSource, fileNameDest)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

// getMySQLInstallationPackagePath returns the mysql server binary package path
func (ose *OSExecutor) getMySQLInstallationPackagePath() string {
	return filepath.Join(viper.GetString(config.MySQLInstallationPackageDirKey), ose.getMySQLServerBinaryPackageName())
}

// getMySQLServerBinaryPackageName returns the mysql installation package name
func (ose *OSExecutor) getMySQLServerBinaryPackageName() string {
	packageNameTemplate := mysqlServerBinaryPackageNameTemplateV1
	if ose.arch == constant.AArch64Arch {
		packageNameTemplate = mysqlServerBinaryPackageNameTemplateV2
	}

	return fmt.Sprintf(packageNameTemplate, ose.mysqlVersion.String())
}

// getMySQLServerBinaryPackageDecompressedDirName returns the mysql server binary package decompressed directory name
func (ose *OSExecutor) getMySQLServerBinaryPackageDecompressedDirName() string {
	packageName := ose.getMySQLServerBinaryPackageName()

	if strings.HasSuffix(packageName, tarXZExt) {
		return strings.TrimSuffix(packageName, tarXZExt)
	}

	return strings.TrimSuffix(packageName, tarGZExt)

}

// getDecompressCommand returns the decompress command
func (ose *OSExecutor) getDecompressCommand() string {
	command := decompressCommandV1
	if ose.arch == constant.AArch64Arch {
		command = decompressCommandV2
	}

	return fmt.Sprintf(command, filepath.Join(constant.DefaultTmpDir, ose.getMySQLServerBinaryPackageName()), constant.DefaultTmpDir)
}
