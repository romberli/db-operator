/*
Copyright © 2020 Romber Li <romber2001@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/log"
	"path/filepath"
	"strings"

	"github.com/romberli/go-util/constant"
	"github.com/spf13/viper"
)

var (
	ValidLogLevels                  = []string{"debug", "info", "warn", "warning", "error", "fatal"}
	ValidLogFormats                 = []string{"text", "json"}
	ValidServerRouterHTTPErrorCodes = []int{200, 500}
)

// SetDefaultConfig set default configuration, it is the lowest priority
func SetDefaultConfig(baseDir string) {
	// daemon
	SetDefaultDaemon()
	// log
	SetDefaultLog(baseDir)
	// server
	SetDefaultServer(baseDir)
	// database
	SetDefaultDB()
	// mysql
	SetDefaultMySQL()
	// pmm
	SetDefaultPMM()
}

// SetDefaultDaemon sets the default value of daemon
func SetDefaultDaemon() {
	viper.SetDefault(DaemonKey, DefaultDaemon)
}

// SetDefaultLog sets the default value of log
func SetDefaultLog(baseDir string) {
	defaultLogFile := filepath.Join(baseDir, DefaultLogDir, log.DefaultLogFileName)
	viper.SetDefault(LogFileNameKey, defaultLogFile)
	viper.SetDefault(LogLevelKey, log.DefaultLogLevel)
	viper.SetDefault(LogFormatKey, log.DefaultLogFormat)
	viper.SetDefault(LogMaxSizeKey, log.DefaultLogMaxSize)
	viper.SetDefault(LogMaxDaysKey, log.DefaultLogMaxDays)
	viper.SetDefault(LogMaxBackupsKey, log.DefaultLogMaxBackups)
	viper.SetDefault(LogRotateOnStartupKey, DefaultRotateOnStartup)
}

// SetDefaultServer sets the default value of server
func SetDefaultServer(baseDir string) {
	viper.SetDefault(ServerAddrKey, DefaultServerAddr)
	defaultPidFile := filepath.Join(baseDir, fmt.Sprintf("%s.pid", DefaultCommandName))
	viper.SetDefault(ServerPidFileKey, defaultPidFile)
	viper.SetDefault(ServerReadTimeoutKey, DefaultServerReadTimeout)
	viper.SetDefault(ServerWriteTimeoutKey, DefaultServerWriteTimeout)
	viper.SetDefault(ServerPProfEnabledKey, DefaultServerPProfEnabled)
	viper.SetDefault(ServerRouterAlternativeBasePathKey, DefaultServerRouterAlternativeBasePath)
	viper.SetDefault(ServerRouterAlternativeBodyPathKey, DefaultServerRouterAlternativeBodyPath)
	viper.SetDefault(ServerRouterHTTPErrorCodeKey, DefaultServerRouterHTTPErrorCode)
}

// SetDefaultDB sets the default value of db
func SetDefaultDB() {
	viper.SetDefault(DBDBOMySQLAddrKey, constant.DefaultMySQLAddr)
	viper.SetDefault(DBDBOMySQLNameKey, DefaultDBName)
	viper.SetDefault(DBDBOMySQLUserKey, DefaultDBUser)
	viper.SetDefault(DBDBOMySQLPassKey, DefaultDBPass)
	viper.SetDefault(DBPoolMaxConnectionsKey, mysql.DefaultMaxConnections)
	viper.SetDefault(DBPoolInitConnectionsKey, mysql.DefaultInitConnections)
	viper.SetDefault(DBPoolMaxIdleConnectionsKey, mysql.DefaultMaxIdleConnections)
	viper.SetDefault(DBPoolMaxIdleTimeKey, mysql.DefaultMaxIdleTime)
	viper.SetDefault(DBPoolMaxWaitTimeKey, mysql.DefaultMaxWaitTime)
	viper.SetDefault(DBPoolMaxRetryCountKey, mysql.DefaultMaxRetryCount)
	viper.SetDefault(DBPoolKeepAliveIntervalKey, mysql.DefaultKeepAliveInterval)
}

// SetDefaultMySQL sets the default value of mysql
func SetDefaultMySQL() {
	viper.SetDefault(MySQLVersionKey, DefaultMySQLVersion)
	viper.SetDefault(MySQLVersionIntKey, DefaultMySQLVersionInt)
	viper.SetDefault(MySQLInstallationPackageDirKey, DefaultMySQLInstallationPackageDir)
	viper.SetDefault(MySQLInstallationTemporaryDirKey, DefaultMySQLInstallationTemporaryDir)
	viper.SetDefault(MySQLParameterMaxConnectionsKey, DefaultMySQLParameterMaxConnections)
	viper.SetDefault(MySQLParameterInnodbBufferPoolSizeKey, DefaultMySQLParameterInnodbBufferPoolSize)
	viper.SetDefault(MySQLParameterInnodbIOCapacityKey, DefaultMySQLParameterInnodbIOCapacity)
	viper.SetDefault(MySQLUserOSUserKey, DefaultMySQLUserOSUser)
	viper.SetDefault(MySQLUserOSPassKey, DefaultMySQLUserOSPass)
	viper.SetDefault(MySQLUserRootPassKey, DefaultMySQLUserRootPass)
	viper.SetDefault(MySQLUserAdminUserKey, DefaultMySQLUserAdminUser)
	viper.SetDefault(MySQLUserAdminPassKey, DefaultMySQLUserAdminPass)
	viper.SetDefault(MySQLUserMySQLDMultiUserKey, DefaultMySQLUserMySQLDMultiUser)
	viper.SetDefault(MySQLUserMySQLDMultiPassKey, DefaultMySQLUserMySQLDMultiPass)
	viper.SetDefault(MySQLUserReplicationUserKey, DefaultMySQLUserReplicationUser)
	viper.SetDefault(MySQLUserReplicationPassKey, DefaultMySQLUserReplicationPass)
	viper.SetDefault(MySQLUserMonitorUserKey, DefaultMySQLUserMonitorUser)
	viper.SetDefault(MySQLUserMonitorPassKey, DefaultMySQLUserMonitorPass)
	viper.SetDefault(MySQLUserDASUserKey, DefaultMySQLUserDASUser)
	viper.SetDefault(MySQLUserDASPassKey, DefaultMySQLUserDASPass)
}

// SetDefaultPMM sets the default value of pmm
func SetDefaultPMM() {
	viper.SetDefault(PMMServerAddrKey, DefaultPMMServerAddr)
	viper.SetDefault(PMMServerUserKey, DefaultPMMServerUser)
	viper.SetDefault(PMMServerPassKey, DefaultPMMServerPass)
	viper.SetDefault(PMMClientVersionKey, DefaultPMMClientVersion)
	viper.SetDefault(PMMClientInstallationPackageDirKey, DefaultPMMClientInstallationPackageDir)
}

// TrimSpaceOfArg trims spaces of given argument
func TrimSpaceOfArg(arg string) string {
	args := strings.SplitN(arg, constant.EqualString, 2)

	switch len(args) {
	case 1:
		return strings.TrimSpace(args[constant.ZeroInt])
	case 2:
		argName := strings.TrimSpace(args[constant.ZeroInt])
		argValue := strings.TrimSpace(args[1])
		return fmt.Sprintf("%s=%s", argName, argValue)
	default:
		return arg
	}
}
