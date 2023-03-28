package cmd

import (
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-multierror"
	"github.com/romberli/go-util/constant"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"strings"
)

// OverrideConfigByCLI read configuration from command line interface, it will override the config file configuration
func OverrideConfigByCLI() error {
	merr := &multierror.Error{}

	// override config
	overrideConfigByCLI()
	// override daemon
	err := overrideDaemonByCLI()
	if err != nil {
		merr = multierror.Append(merr, err)
	}
	// override log
	err = overrideLogByCLI()
	if err != nil {
		merr = multierror.Append(merr, err)
	}
	// override server
	err = overrideServerByCLI()
	if err != nil {
		merr = multierror.Append(merr, err)
	}
	// override database
	overrideDatabaseByCLI()
	// override mysql
	overrideMySQLByCLI()
	// override pmm
	overridePMMByCLI()

	if merr.ErrorOrNil() != nil {
		return message.NewMessage(message.ErrOverrideConfigByCLI, merr.ErrorOrNil())
	}

	// validate configuration
	err = config.ValidateConfig()
	if err != nil {
		return message.NewMessage(message.ErrValidateConfig, err)
	}

	return nil
}

// overrideConfigByCLI overrides the config section by command line interface
func overrideConfigByCLI() {
	if cfgFile != constant.EmptyString && cfgFile != constant.DefaultRandomString {
		viper.Set(config.ConfKey, cfgFile)
	}
}

// overrideDaemonByCLI overrides the daemon section by command line interface
func overrideDaemonByCLI() error {
	if daemonStr != constant.DefaultRandomString {
		daemon, err := cast.ToBoolE(daemonStr)
		if err != nil {
			return errors.Trace(err)
		}

		viper.Set(config.DaemonKey, daemon)
	}

	return nil
}

// overrideLogByCLI overrides the log section by command line interface
func overrideLogByCLI() error {
	if logFileName != constant.DefaultRandomString {
		viper.Set(config.LogFileNameKey, logFileName)
	}
	if logLevel != constant.DefaultRandomString {
		logLevel = strings.ToLower(logLevel)
		viper.Set(config.LogLevelKey, logLevel)
	}
	if logFormat != constant.DefaultRandomString {
		logLevel = strings.ToLower(logFormat)
		viper.Set(config.LogFormatKey, logFormat)
	}
	if logMaxSize != constant.DefaultRandomInt {
		viper.Set(config.LogMaxSizeKey, logMaxSize)
	}
	if logMaxDays != constant.DefaultRandomInt {
		viper.Set(config.LogMaxDaysKey, logMaxDays)
	}
	if logMaxBackups != constant.DefaultRandomInt {
		viper.Set(config.LogMaxBackupsKey, logMaxBackups)
	}
	if logRotateOnStartupStr != constant.DefaultRandomString {
		rotateOnStartup, err := cast.ToBoolE(logRotateOnStartupStr)
		if err != nil {
			return errors.Trace(err)
		}

		viper.Set(config.LogRotateOnStartupKey, rotateOnStartup)
	}

	return nil
}

// overrideServerByCLI overrides the server section by command line interface
func overrideServerByCLI() error {
	if serverAddr != constant.DefaultRandomString {
		viper.Set(config.ServerAddrKey, serverAddr)
	}
	if serverPidFile != constant.DefaultRandomString {
		viper.Set(config.ServerPidFileKey, serverPidFile)
	}
	if serverReadTimeout != constant.DefaultRandomInt {
		viper.Set(config.ServerReadTimeoutKey, serverReadTimeout)
	}
	if serverWriteTimeout != constant.DefaultRandomInt {
		viper.Set(config.ServerWriteTimeoutKey, serverWriteTimeout)
	}
	if serverPProfEnabledStr != constant.DefaultRandomString {
		pprofEnabled, err := cast.ToBoolE(serverPProfEnabledStr)
		if err != nil {
			return errors.Trace(err)
		}

		viper.Set(config.ServerPProfEnabledKey, pprofEnabled)
	}
	if serverRouterAlternativeBasePath != constant.DefaultRandomString {
		viper.Set(config.ServerRouterAlternativeBasePathKey, serverRouterAlternativeBasePath)
	}
	if serverRouterAlternativeBodyPath != constant.DefaultRandomString {
		viper.Set(config.ServerRouterAlternativeBodyPathKey, serverRouterAlternativeBodyPath)
	}
	if serverRouterHTTPErrorCode != constant.DefaultRandomInt {
		viper.Set(config.ServerRouterHTTPErrorCodeKey, serverRouterHTTPErrorCode)
	}

	return nil
}

// overrideDatabaseByCLI overrides the db section by command line interface
func overrideDatabaseByCLI() {
	if dbDBOMySQLAddr != constant.DefaultRandomString {
		viper.Set(config.DBDBOMySQLAddrKey, dbDBOMySQLAddr)
	}
	if dbDBOMySQLName != constant.DefaultRandomString {
		viper.Set(config.DBDBOMySQLNameKey, dbDBOMySQLName)
	}
	if dbDBOMySQLUser != constant.DefaultRandomString {
		viper.Set(config.DBDBOMySQLUserKey, dbDBOMySQLUser)
	}
	if dbDBOMySQLPass != constant.DefaultRandomString {
		viper.Set(config.DBDBOMySQLPassKey, dbDBOMySQLPass)
	}
	if dbPoolMaxConnections != constant.DefaultRandomInt {
		viper.Set(config.DBPoolMaxConnectionsKey, dbPoolMaxConnections)
	}
	if dbPoolInitConnections != constant.DefaultRandomInt {
		viper.Set(config.DBPoolInitConnectionsKey, dbPoolInitConnections)
	}
	if dbPoolMaxIdleConnections != constant.DefaultRandomInt {
		viper.Set(config.DBPoolMaxIdleConnectionsKey, dbPoolMaxIdleConnections)
	}
	if dbPoolMaxIdleTime != constant.DefaultRandomInt {
		viper.Set(config.DBPoolMaxIdleTimeKey, dbPoolMaxIdleTime)
	}
	if dbPoolMaxWaitTime != constant.DefaultRandomInt {
		viper.Set(config.DBPoolMaxWaitTimeKey, dbPoolMaxWaitTime)
	}
	if dbPoolMaxRetryCount != constant.DefaultRandomInt {
		viper.Set(config.DBPoolMaxRetryCountKey, dbPoolMaxRetryCount)
	}
	if dbPoolKeepAliveInterval != constant.DefaultRandomInt {
		viper.Set(config.DBPoolKeepAliveIntervalKey, dbPoolKeepAliveInterval)
	}
}

// overrideMySQLByCLI overrides the mysql section by command line interface
func overrideMySQLByCLI() {
	if mysqlVersion != constant.DefaultRandomString {
		viper.Set(config.MySQLVersionKey, mysqlVersion)
	}
	if mysqlInstallationPackageDir != constant.DefaultRandomString {
		viper.Set(config.MySQLInstallationPackageDirKey, mysqlInstallationPackageDir)
	}
	if mysqlInstallationTemporaryDir != constant.DefaultRandomString {
		viper.Set(config.MySQLInstallationTemporaryDirKey, mysqlInstallationTemporaryDir)
	}
	if mysqlParameterMaxConnections != constant.DefaultRandomInt {
		viper.Set(config.MySQLParameterMaxConnectionsKey, mysqlParameterMaxConnections)
	}
	if mysqlParameterInnodbBufferPoolSize != constant.DefaultRandomInt {
		viper.Set(config.MySQLParameterInnodbBufferPoolSizeKey, mysqlParameterInnodbBufferPoolSize)
	}
	if mysqlParameterInnodbIOCapacity != constant.DefaultRandomInt {
		viper.Set(config.MySQLParameterInnodbIOCapacityKey, mysqlParameterInnodbIOCapacity)
	}
	if mysqlUserOSUser != constant.DefaultRandomString {
		viper.Set(config.MySQLUserOSUserKey, mysqlUserOSUser)
	}
	if mysqlUserOSPass != constant.DefaultRandomString {
		viper.Set(config.MySQLUserOSPassKey, mysqlUserOSPass)
	}
	if mysqlUserRootPass != constant.DefaultRandomString {
		viper.Set(config.MySQLUserRootPassKey, mysqlUserRootPass)
	}
	if mysqlUserAdminUser != constant.DefaultRandomString {
		viper.Set(config.MySQLUserAdminUserKey, mysqlUserAdminUser)
	}
	if mysqlUserAdminPass != constant.DefaultRandomString {
		viper.Set(config.MySQLUserAdminPassKey, mysqlUserAdminPass)
	}
	if mysqlUserMySQLDMultiUser != constant.DefaultRandomString {
		viper.Set(config.MySQLUserMySQLDMultiUserKey, mysqlUserMySQLDMultiUser)
	}
	if mysqlUserMySQLDMultiPass != constant.DefaultRandomString {
		viper.Set(config.MySQLUserMySQLDMultiPassKey, mysqlUserMySQLDMultiPass)
	}
	if mysqlUserReplicationUser != constant.DefaultRandomString {
		viper.Set(config.MySQLUserReplicationUserKey, mysqlUserReplicationUser)
	}
	if mysqlUserReplicationPass != constant.DefaultRandomString {
		viper.Set(config.MySQLUserReplicationPassKey, mysqlUserReplicationPass)
	}
	if mysqlUserMonitorUser != constant.DefaultRandomString {
		viper.Set(config.MySQLUserMonitorUserKey, mysqlUserMonitorUser)
	}
	if mysqlUserMonitorPass != constant.DefaultRandomString {
		viper.Set(config.MySQLUserMonitorPassKey, mysqlUserMonitorPass)
	}
	if mysqlUserDASUser != constant.DefaultRandomString {
		viper.Set(config.MySQLUserDASUserKey, mysqlUserDASUser)
	}
	if mysqlUserDASPass != constant.DefaultRandomString {
		viper.Set(config.MySQLUserDASPassKey, mysqlUserDASPass)
	}
	if mysqlOperationTimeout != constant.DefaultRandomInt {
		viper.Set(config.MySQLOperationTimeoutKey, mysqlOperationTimeout)
	}
}

// overridePMMByCLI overrides the pmm section by command line interface
func overridePMMByCLI() {
	if pmmServerAddr != constant.DefaultRandomString {
		viper.Set(config.PMMServerAddrKey, pmmServerAddr)
	}
	if pmmServerUser != constant.DefaultRandomString {
		viper.Set(config.PMMServerUserKey, pmmServerUser)
	}
	if pmmServerPass != constant.DefaultRandomString {
		viper.Set(config.PMMServerPassKey, pmmServerPass)
	}
	if pmmClientVersion != constant.DefaultRandomString {
		viper.Set(config.PMMClientVersionKey, pmmClientVersion)
	}
	if pmmClientInstallationPackageDir != constant.DefaultRandomString {
		viper.Set(config.PMMClientInstallationPackageDirKey, pmmClientInstallationPackageDir)
	}
}
