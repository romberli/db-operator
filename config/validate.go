package config

import (
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/spf13/cast"
	"github.com/spf13/viper"

	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-multierror"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"

	msgMySQL "github.com/romberli/db-operator/pkg/message/mysql"
	msgPMM "github.com/romberli/db-operator/pkg/message/pmm"
)

const (
	minMySQLVersion     = "5.7.0"
	minPMMClientVersion = "2.0.0"
)

// ValidateConfig validates if the configuration is valid
func ValidateConfig() (err error) {
	merr := &multierror.Error{}

	// validate daemon section
	err = ValidateDaemon()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	// validate log section
	err = ValidateLog()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	// validate server section
	err = ValidateServer()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	// validate mysql section
	err = ValidateMySQL()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	// validate pmm section
	err = ValidatePMM()
	if err != nil {
		merr = multierror.Append(merr, err)
	}

	return errors.Trace(merr.ErrorOrNil())
}

// ValidateDaemon validates if daemon section is valid
func ValidateDaemon() error {
	_, err := cast.ToBoolE(viper.Get(DaemonKey))

	return errors.Trace(err)
}

// ValidateLog validates if log section is valid.
func ValidateLog() error {
	var valid bool

	merr := &multierror.Error{}

	// validate log.FileName
	logFileName, err := cast.ToStringE(viper.Get(LogFileNameKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	logFileName = strings.TrimSpace(logFileName)
	if logFileName == constant.EmptyString {
		merr = multierror.Append(merr, message.NewMessage(message.ErrEmptyLogFileName))
	}
	isAbs := filepath.IsAbs(logFileName)
	if !isAbs {
		logFileName, err = filepath.Abs(logFileName)
		if err != nil {
			merr = multierror.Append(merr, errors.Trace(err))
		}
	}
	valid, _ = govalidator.IsFilePath(logFileName)
	if !valid {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogFileName, logFileName))
	}

	// validate log.level
	logLevel, err := cast.ToStringE(viper.Get(LogLevelKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	valid, err = common.ElementInSlice(ValidLogLevels, logLevel)
	if err != nil {
		merr = multierror.Append(merr, err)
	}
	if !valid {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogLevel, logLevel))
	}

	// validate log.format
	logFormat, err := cast.ToStringE(viper.Get(LogFormatKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	valid, err = common.ElementInSlice(ValidLogFormats, logFormat)
	if err != nil {
		merr = multierror.Append(merr, err)
	}
	if !valid {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogFormat, logFormat))
	}

	// validate log.maxSize
	logMaxSize, err := cast.ToIntE(viper.Get(LogMaxSizeKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	if logMaxSize < MinLogMaxSize || logMaxSize > MaxLogMaxSize {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogMaxSize, MinLogMaxSize, MaxLogMaxSize, logMaxSize))
	}

	// validate log.maxDays
	logMaxDays, err := cast.ToIntE(viper.Get(LogMaxDaysKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	if logMaxDays < MinLogMaxDays || logMaxDays > MaxLogMaxDays {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogMaxDays, MinLogMaxDays, MaxLogMaxDays, logMaxDays))
	}

	// validate log.maxBackups
	logMaxBackups, err := cast.ToIntE(viper.Get(LogMaxBackupsKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	if logMaxBackups < MinLogMaxDays || logMaxBackups > MaxLogMaxDays {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidLogMaxBackups, MinLogMaxBackups, MaxLogMaxBackups, logMaxBackups))
	}

	// validate log.rotateOnStartup
	_, err = cast.ToBoolE(viper.Get(LogRotateOnStartupKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	return merr.ErrorOrNil()
}

// ValidateServer validates if server section is valid
func ValidateServer() error {
	merr := &multierror.Error{}

	// validate server.addr
	serverAddr, err := cast.ToStringE(viper.Get(ServerAddrKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	serverAddrList := strings.Split(serverAddr, constant.ColonString)

	switch len(serverAddrList) {
	case 2:
		port := serverAddrList[1]
		if !govalidator.IsPort(port) {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidServerPort, constant.MinPort, constant.MaxPort, port))
		}
	default:
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidServerAddr, serverAddr))
	}

	// validate server.pidFile
	serverPidFile, err := cast.ToStringE(viper.Get(ServerPidFileKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	isAbs := filepath.IsAbs(serverPidFile)
	if !isAbs {
		serverPidFile, err = filepath.Abs(serverPidFile)
		if err != nil {
			merr = multierror.Append(merr, errors.Trace(err))
		}
	}
	ok, _ := govalidator.IsFilePath(serverPidFile)
	if !ok {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidPidFile, serverPidFile))
	}

	// validate server.readTimeout
	serverReadTimeout, err := cast.ToIntE(viper.Get(ServerReadTimeoutKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	if serverReadTimeout < MinServerReadTimeout || serverReadTimeout > MaxServerReadTimeout {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidServerReadTimeout, MinServerReadTimeout, MaxServerWriteTimeout, serverReadTimeout))
	}

	// validate server.writeTimeout
	serverWriteTimeout, err := cast.ToIntE(viper.Get(ServerWriteTimeoutKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}
	if serverWriteTimeout < MinServerWriteTimeout || serverWriteTimeout > MaxServerWriteTimeout {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidServerWriteTimeout, MinServerWriteTimeout, MaxServerWriteTimeout, serverWriteTimeout))
	}

	// validate server.pprof.enabled
	_, err = cast.ToBoolE(viper.Get(ServerPProfEnabledKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate server.router.alternativeBaseURL
	serverRouterAlternativeBasePath, err := cast.ToStringE(viper.Get(ServerRouterAlternativeBasePathKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else if serverRouterAlternativeBasePath != constant.EmptyString && !strings.HasPrefix(serverRouterAlternativeBasePath, constant.SlashString) {
		merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidServerRouterAlternativeBasePath, serverRouterAlternativeBasePath))
	}

	// validate server.router.bodyPath
	_, err = cast.ToStringE(viper.Get(ServerRouterAlternativeBodyPathKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate server.router.httpErrorCode
	serverRouterHttpErrorCode, err := cast.ToIntE(viper.Get(ServerRouterHTTPErrorCodeKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		valid, err := common.ElementInSlice(ValidServerRouterHTTPErrorCodes, serverRouterHttpErrorCode)
		if err != nil {
			merr = multierror.Append(merr, err)
		}
		if !valid {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidServerRouterHTTPErrorCode, serverRouterHttpErrorCode))
		}
	}

	return merr.ErrorOrNil()
}

// ValidateMySQL validates if MySQL section is valid
func ValidateMySQL() error {
	merr := &multierror.Error{}

	// validate mysql.version
	mysqlVersion, err := cast.ToStringE(viper.Get(MySQLVersionKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		v, err := version.NewVersion(mysqlVersion)
		if err != nil {
			merr = multierror.Append(merr, errors.Trace(err))
		} else {
			minVersion, err := version.NewVersion(minMySQLVersion)
			if err != nil {
				merr = multierror.Append(merr, errors.Trace(err))
			} else {
				if v.LessThan(minVersion) {
					merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLVersion, minMySQLVersion, mysqlVersion))
				}
			}
		}
	}

	// validate mysql.installationPackageDir
	mysqlInstallationPackageDir, err := cast.ToStringE(viper.Get(MySQLInstallationPackageDirKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		ok, _ := govalidator.IsFilePath(mysqlInstallationPackageDir)
		if !ok {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidFilePath, mysqlInstallationPackageDir))
		}
	}

	// validate mysql.parameter.maxConnections
	maxConnections, err := cast.ToIntE(viper.Get(MySQLParameterMaxConnectionsKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if maxConnections < MinMySQLParameterMaxConnections || maxConnections > MaxMySQLParameterMaxConnections {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLParameterMaxConnections, MinMySQLParameterMaxConnections, MaxMySQLParameterMaxConnections, maxConnections))
		}
	}

	// validate mysql.parameter.innodbBufferPoolSize
	innodbBufferPoolSize, err := cast.ToIntE(viper.Get(MySQLParameterInnodbBufferPoolSizeKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if innodbBufferPoolSize < MinMySQLParameterInnodbBufferPoolSize || innodbBufferPoolSize > MaxMySQLParameterInnodbBufferPoolSize {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLParameterInnodbBufferPoolSize, MinMySQLParameterInnodbBufferPoolSize, MaxMySQLParameterInnodbBufferPoolSize, innodbBufferPoolSize))
		}
	}

	// validate mysql.parameter.innodbIOCapacity
	innodbIOCapacity, err := cast.ToIntE(viper.Get(MySQLParameterInnodbIOCapacityKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if innodbIOCapacity < MinMySQLParameterInnodbIOCapacity || innodbIOCapacity > MaxMySQLParameterInnodbIOCapacity {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLParameterInnodbIOCapacity, MinMySQLParameterInnodbIOCapacity, MaxMySQLParameterInnodbIOCapacity, innodbIOCapacity))
		}
	}

	// validate mysql.user.osUser
	osUser, err := cast.ToStringE(viper.Get(MySQLUserOSUserKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if osUser == constant.EmptyString {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLUser, MySQLUserOSUserKey))
		}
	}

	// validate mysql.user.osPass
	_, err = cast.ToStringE(viper.Get(MySQLUserOSPassKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate mysql.user.rootPass
	_, err = cast.ToStringE(viper.Get(MySQLUserRootPassKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate mysql.user.adminUser
	adminUser, err := cast.ToStringE(viper.Get(MySQLUserAdminUserKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if adminUser == constant.EmptyString {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLUser, MySQLUserAdminUserKey))
		}
	}

	// validate mysql.user.adminPass
	_, err = cast.ToStringE(viper.Get(MySQLUserAdminPassKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate mysql.user.mysqldMultiUser
	mysqldMultiUser, err := cast.ToStringE(viper.Get(MySQLUserMySQLDMultiUserKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if mysqldMultiUser == constant.EmptyString {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLUser, MySQLUserMonitorUserKey))
		}
	}

	// validate mysql.user.mysqldMultiPass
	_, err = cast.ToStringE(viper.Get(MySQLUserMySQLDMultiPassKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate mysql.user.replicationUser
	replicationUser, err := cast.ToStringE(viper.Get(MySQLUserReplicationUserKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if replicationUser == constant.EmptyString {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLUser, MySQLUserReplicationUserKey))
		}
	}

	// validate mysql.user.replicationPass
	_, err = cast.ToStringE(viper.Get(MySQLUserReplicationPassKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate mysql.user.monitorUser
	monitorUser, err := cast.ToStringE(viper.Get(MySQLUserMonitorUserKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if monitorUser == constant.EmptyString {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLUser, MySQLUserMonitorUserKey))
		}
	}

	// validate mysql.user.monitorPass
	_, err = cast.ToStringE(viper.Get(MySQLUserMonitorPassKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate mysql.user.dasUser
	dasUser, err := cast.ToStringE(viper.Get(MySQLUserDASUserKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if dasUser == constant.EmptyString {
			merr = multierror.Append(merr, message.NewMessage(msgMySQL.ErrNotValidConfigMySQLUser, MySQLUserDASUserKey))
		}
	}

	// validate mysql.user.dasPass
	_, err = cast.ToStringE(viper.Get(MySQLUserDASPassKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	return merr.ErrorOrNil()
}

// ValidatePMM validates if pmm section is valid
func ValidatePMM() error {
	merr := &multierror.Error{}

	// validate pmm.server.addr
	addr, err := cast.ToStringE(viper.Get(PMMServerAddrKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if addr == constant.EmptyString {
			merr = multierror.Append(merr, message.NewMessage(msgPMM.ErrNotValidConfigPMMServerAddr, addr))
		}
	}

	// validate pmm.server.user
	user, err := cast.ToStringE(viper.Get(PMMServerUserKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		if user == constant.EmptyString {
			merr = multierror.Append(merr, message.NewMessage(msgPMM.ErrNotValidConfigPMMServerUser, user))
		}
	}

	// validate pmm.server.pass
	_, err = cast.ToStringE(viper.Get(PMMServerPassKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	}

	// validate pmm.client.version
	clientVersion, err := cast.ToStringE(viper.Get(PMMClientVersionKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		v, err := version.NewVersion(minPMMClientVersion)
		if err != nil {
			merr = multierror.Append(merr, errors.Trace(err))
		} else {
			minClientVersion, err := version.NewVersion(minPMMClientVersion)
			if err != nil {
				merr = multierror.Append(merr, errors.Trace(err))
			} else {
				if v.LessThan(minClientVersion) {
					merr = multierror.Append(merr, message.NewMessage(msgPMM.ErrNotValidConfigPMMClientVersion, minPMMClientVersion, clientVersion))
				}
			}
		}
	}

	// validate pmm.client.installationPackageDir
	clientInstallationPackageDir, err := cast.ToStringE(viper.Get(PMMClientInstallationPackageDirKey))
	if err != nil {
		merr = multierror.Append(merr, errors.Trace(err))
	} else {
		ok, _ := govalidator.IsFilePath(clientInstallationPackageDir)
		if !ok {
			merr = multierror.Append(merr, message.NewMessage(message.ErrNotValidFilePath, clientInstallationPackageDir))
		}
	}

	return merr.ErrorOrNil()
}
