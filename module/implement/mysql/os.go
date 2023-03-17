package mysql

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/module/implement/mysql/parameter"
	"github.com/romberli/db-operator/pkg/util/ssh"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"github.com/spf13/viper"
)

const (
	os9VersionStr           = "9.0.0"
	mysql8026VersionStr     = "8.0.26"
	mysql8032VersionStr     = "8.0.32"
	minX64MySQLVersionStr   = mysql8026VersionStr
	minAArchMySQLVersionStr = mysql8032VersionStr
	// mysql 8.0
	mysqlInstallationPackageNameTemplateV1 = "mysql-%s-linux-glibc2.12-x86_64.tar.xz"
	mysqlInstallationPackageNameTemplateV2 = "mysql-%s-linux-glibc2.17-aarch64.tar.gz"
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

	yumInstallCommand       = "/usr/bin/yum Install -y ncurses-c++-libs ncurses-libs"
	libNCursesPath          = "/usr/lib64/libncurses.so.5"
	libTInfoPath            = "/usr/lib64/libtinfo.so.5"
	lnLibNCursesCommand     = "/usr/bin/ln -s /usr/lib64/libncurses.so.6.2 /usr/lib64/libncurses.so.5"
	lnLibTInfoCommand       = "/usr/bin/ln -s /usr/lib64/libtinfo.so.6.2 /usr/lib64/libtinfo.so.5"
	checkMySQLGroupCommand  = "/usr/bin/id -g mysql"
	checkMySQLUserCommand   = "/usr/bin/id -u mysql"
	createMySQLGroupCommand = "/usr/sbin/groupadd -g 1001 mysql"
	createMySQLUserCommand  = "/usr/sbin/useradd -u 1001 -g mysql mysql"
)

var (
	minAArchMySQLVersion = version.Must(version.NewVersion(minAArchMySQLVersionStr))
	minX64MySQLVersion   = version.Must(version.NewVersion(minX64MySQLVersionStr))
	os9Version           = version.Must(version.NewVersion(os9VersionStr))
)

type OSExecutor struct {
	arch         string
	osVersion    *version.Version
	mysqlVersion *version.Version
	sshConn      *ssh.Conn
	mysqlServer  *parameter.MySQLServer
}

// NewOSExecutor returns a new *OSExecutor
func NewOSExecutor(mysqlVersion *version.Version, sshConn *ssh.Conn, mysqlServer *parameter.MySQLServer) *OSExecutor {
	return newOSExecutor(mysqlVersion, sshConn, mysqlServer)
}

// newOSExecutor returns a new *OSExecutor
func newOSExecutor(mysqlVersion *version.Version, sshConn *ssh.Conn, mysqlServer *parameter.MySQLServer) *OSExecutor {
	return &OSExecutor{
		mysqlVersion: mysqlVersion,
		sshConn:      sshConn,
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

	return nil
}

// InitExecutor initializes the os executor
func (ose *OSExecutor) InitExecutor() error {
	var err error
	// get os version
	ose.osVersion, err = ose.sshConn.GetOSVersion()
	if err != nil {
		return err
	}
	// get arch
	ose.arch, err = ose.sshConn.GetArch()
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
	// check if mysql installation package exists
	installationPackagePath := filepath.Join(viper.GetString(config.MySQLInstallationPackageDirKey), ose.getMySQLInstallationPackageName())
	exists, err := linux.PathExists(installationPackagePath)
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("mysql installation package does not exist. installation package path: %s", installationPackagePath)
	}
	// check if the mysql data directory exists
	dataDir := filepath.Join(ose.mysqlServer.DataDirBase, dataDirName)
	output, err := ose.sshConn.ListPath(dataDir)
	if err == nil && len(output) > constant.ZeroInt {
		return errors.Errorf("mysql data directory exists and is not empty, installation aborted. data directory: %s", dataDir)
	}
	// check if the mysql binlog directory exists
	binlogDir := filepath.Join(ose.mysqlServer.LogDirBase, binlogDirName)
	output, err = ose.sshConn.ListPath(binlogDir)
	if err == nil && len(output) > constant.ZeroInt {
		return errors.Errorf("mysql binlog directory exists and is not empty, installation aborted. binlog directory: %s", dataDir)
	}
	// check if the mysql relaylog directory exists
	relaylogDir := filepath.Join(ose.mysqlServer.LogDirBase, relaylogDirName)
	output, err = ose.sshConn.ListPath(relaylogDir)
	if err == nil && len(output) > constant.ZeroInt {
		return errors.Errorf("mysql relaylog directory exists and is not empty, installation aborted. relaylog directory: %s", dataDir)
	}

	return nil
}

// InstallRPM installs the rpm
func (ose *OSExecutor) InstallRPM() error {
	err := ose.sshConn.ExecuteCommandWithoutOutput(yumInstallCommand)
	if err != nil {
		return err
	}

	if ose.osVersion.GreaterThanOrEqual(os9Version) {
		pathExists, err := ose.sshConn.PathExists(libNCursesPath)
		if err != nil {
			return err
		}
		if !pathExists {
			err = ose.sshConn.ExecuteCommandWithoutOutput(lnLibNCursesCommand)
			if err != nil {
				return err
			}
		}
		pathExists, err = ose.sshConn.PathExists(libTInfoPath)
		if err != nil {
			return err
		}
		if !pathExists {
			err = ose.sshConn.ExecuteCommandWithoutOutput(lnLibTInfoCommand)
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
	err := ose.sshConn.ExecuteCommandWithoutOutput(checkMySQLGroupCommand)
	if err != nil {
		err = ose.sshConn.ExecuteCommandWithoutOutput(createMySQLGroupCommand)
		if err != nil {
			return errors.Trace(err)
		}
	}
	// init mysql user
	err = ose.sshConn.ExecuteCommandWithoutOutput(checkMySQLUserCommand)
	if err != nil {
		err = ose.sshConn.ExecuteCommandWithoutOutput(createMySQLUserCommand)
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
	err := ose.sshConn.MkdirAll(binaryDirParent)
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.sshConn.MkdirAll(ose.mysqlServer.BackupDir)
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.sshConn.MkdirAll(filepath.Join(ose.mysqlServer.DataDirBase, dataDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.sshConn.MkdirAll(filepath.Join(ose.mysqlServer.DataDirBase, logDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.sshConn.MkdirAll(filepath.Join(ose.mysqlServer.DataDirBase, tmpDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.sshConn.MkdirAll(filepath.Join(ose.mysqlServer.DataDirBase, runDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.sshConn.MkdirAll(filepath.Join(ose.mysqlServer.LogDirBase, binlogDirName))
	if err != nil {
		return errors.Trace(err)
	}
	err = ose.sshConn.MkdirAll(filepath.Join(ose.mysqlServer.LogDirBase, relaylogDirName))
	if err != nil {
		return errors.Trace(err)
	}
	// change owner of directories
	err = ose.sshConn.Chown(ose.mysqlServer.BackupDir, defaultMySQLUser, defaultMySQLGroup)
	if err != nil {
		return err
	}
	err = ose.sshConn.Chown(ose.mysqlServer.DataDirBase, defaultMySQLUser, defaultMySQLGroup)
	if err != nil {
		return err
	}
	err = ose.sshConn.Chown(ose.mysqlServer.LogDirBase, defaultMySQLUser, defaultMySQLGroup)
	if err != nil {
		return err
	}

	return nil
}

// InstallMySQLBinary installs the mysql binary
func (ose *OSExecutor) InstallMySQLBinary() error {
	// copy mysql installation package
	err := ose.copyInstallationPackages()
	if err != nil {
		return err
	}
	// Install mysql binary
	cmd := ose.getDecompressCommand()
	err = ose.sshConn.ExecuteCommandWithoutOutput(cmd)
	if err != nil {
		return err
	}

	err = ose.sshConn.Move(filepath.Join(
		constant.DefaultTmpDir, ose.getMySQLInstallationPackageDecompressedDirName()), ose.mysqlServer.BinaryDirBase)
	if err != nil {
		return err
	}

	return nil
}

// copyInstallationPackages copies the installation packages to the remote host
func (ose *OSExecutor) copyInstallationPackages() error {
	fileName := ose.getMySQLInstallationPackageName()
	fileNameSource := ose.getMySQLInstallationPackagePath()
	fileNameDest := filepath.Join(constant.DefaultTmpDir, fileName)
	// copy mysql installation package
	err := ose.sshConn.CopySingleFileToRemote(fileNameSource, fileNameDest)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

// getMySQLInstallationPackagePath returns the mysql installation package path
func (ose *OSExecutor) getMySQLInstallationPackagePath() string {
	return filepath.Join(viper.GetString(config.MySQLInstallationPackageDirKey), ose.getMySQLInstallationPackageName())
}

// getMySQLInstallationPackageName returns the mysql installation package name
func (ose *OSExecutor) getMySQLInstallationPackageName() string {
	packageNameTemplate := mysqlInstallationPackageNameTemplateV1
	if ose.arch == constant.AArch64Arch {
		packageNameTemplate = mysqlInstallationPackageNameTemplateV2
	}

	return fmt.Sprintf(packageNameTemplate, ose.mysqlVersion.String())
}

// getMySQLInstallationPackageDecompressedDirName returns the mysql installation package decompressed directory name
func (ose *OSExecutor) getMySQLInstallationPackageDecompressedDirName() string {
	packageName := ose.getMySQLInstallationPackageName()

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

	return fmt.Sprintf(command, filepath.Join(constant.DefaultTmpDir, ose.getMySQLInstallationPackageName()), constant.DefaultTmpDir)
}
