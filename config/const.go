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
	"github.com/romberli/go-util/constant"
	"net/http"
)

// global varariable
const (
	DefaultCommandName = "db-operator"
	DefaultBaseDir     = constant.CurrentDir
	// daemon
	DefaultDaemon  = false
	DaemonArg      = "--daemon"
	DaemonArgTrue  = "--daemon=true"
	DaemonArgFalse = "--daemon=false"
	// log
	DefaultLogDir          = "./log"
	MinLogMaxSize          = 1
	MaxLogMaxSize          = constant.MaxInt
	MinLogMaxDays          = 1
	MaxLogMaxDays          = constant.MaxInt
	MinLogMaxBackups       = 1
	MaxLogMaxBackups       = constant.MaxInt
	DefaultRotateOnStartup = false
	// server
	DefaultServerAddr                      = "0.0.0.0:80"
	DefaultServerReadTimeout               = 5
	DefaultServerWriteTimeout              = 10
	MinServerReadTimeout                   = 0
	MaxServerReadTimeout                   = 60
	MinServerWriteTimeout                  = 1
	MaxServerWriteTimeout                  = 60
	DefaultServerPProfEnabled              = false
	DefaultServerRouterAlternativeBasePath = constant.EmptyString
	DefaultServerRouterAlternativeBodyPath = constant.EmptyString
	DefaultServerRouterHTTPErrorCode       = http.StatusInternalServerError
	// mysql
	DefaultMySQLVersion                       = "8.0.32"
	DefaultMySQLInstallationPackageDir        = "/data/software/mysql"
	MinMySQLParameterMaxConnections           = 1
	MaxMySQLParameterMaxConnections           = 10000
	DefaultMySQLParameterMaxConnections       = 2000
	MinMySQLParameterInnodbBufferPoolSize     = 1
	MaxMySQLParameterInnodbBufferPoolSize     = 1024 * 1024 * 1024
	DefaultMySQLParameterInnodbBufferPoolSize = 1024
	MinMySQLParameterInnodbIOCapacity         = 1
	MaxMySQLParameterInnodbIOCapacity         = 10000000
	DefaultMySQLParameterInnodbIOCapacity     = 1000
	DefaultMySQLUserOSUser                    = "root"
	DefaultMySQLUserOSPass                    = "root"
	DefaultMySQLUserRootPass                  = "root"
	DefaultMySQLUserAdminUser                 = "admin"
	DefaultMySQLUserAdminPass                 = "admin"
	DefaultMySQLUserMySQLDMultiUser           = "mysqld_multi"
	DefaultMySQLUserMySQLDMultiPass           = "mysqld_multi"
	DefaultMySQLUserReplicationUser           = "replication"
	DefaultMySQLUserReplicationPass           = "replication"
	DefaultMySQLUserMonitorUser               = "pmm"
	DefaultMySQLUserMonitorPass               = "pmm"
	DefaultMySQLUserDASUser                   = "das"
	DefaultMySQLUserDASPass                   = "das"
	// pmm
	DefaultPMMServerAddr                   = "127.0.0.1:443"
	DefaultPMMServerUser                   = "admin"
	DefaultPMMServerPass                   = "admin"
	DefaultPMMClientVersion                = "2.34.0"
	DefaultPMMClientInstallationPackageDir = "/data/software/mysql"
)

// configuration varariable
const (
	// config
	ConfKey = "config"
	// daemon
	DaemonKey = "daemon"
	// log
	LogFileNameKey        = "log.fileName"
	LogLevelKey           = "log.level"
	LogFormatKey          = "log.format"
	LogMaxSizeKey         = "log.maxSize"
	LogMaxDaysKey         = "log.maxDays"
	LogMaxBackupsKey      = "log.maxBackups"
	LogRotateOnStartupKey = "log.rotateOnStartup"
	// server
	ServerAddrKey                      = "server.addr"
	ServerPidFileKey                   = "server.pidFile"
	ServerReadTimeoutKey               = "server.readTimeout"
	ServerWriteTimeoutKey              = "server.writeTimeout"
	ServerPProfEnabledKey              = "server.pprof.enabled"
	ServerRouterAlternativeBasePathKey = "server.router.alternativeBasePath"
	ServerRouterAlternativeBodyPathKey = "server.router.alternativeBodyPath"
	ServerRouterHTTPErrorCodeKey       = "server.router.httpErrorCode"
	// mysql
	MySQLVersionKey                       = "mysql.version"
	MySQLInstallationPackageDirKey        = "mysql.installationPackageDir"
	MySQLParameterMaxConnectionsKey       = "mysql.parameter.maxConnections"
	MySQLParameterInnodbBufferPoolSizeKey = "mysql.parameter.innodbBufferPoolSize"
	MySQLParameterInnodbIOCapacityKey     = "mysql.parameter.innodbIOCapacity"
	MySQLUserOSUserKey                    = "mysql.user.osUser"
	MySQLUserOSPassKey                    = "mysql.user.osPass"
	MySQLUserRootPassKey                  = "mysql.user.rootPass"
	MySQLUserAdminUserKey                 = "mysql.user.adminUser"
	MySQLUserAdminPassKey                 = "mysql.user.adminPass"
	MySQLUserMySQLDMultiUserKey           = "mysql.user.mysqldMultiUser"
	MySQLUserMySQLDMultiPassKey           = "mysql.user.mysqldMultiPass"
	MySQLUserReplicationUserKey           = "mysql.user.replicationUser"
	MySQLUserReplicationPassKey           = "mysql.user.replicationPass"
	MySQLUserMonitorUserKey               = "mysql.user.monitorUser"
	MySQLUserMonitorPassKey               = "mysql.user.monitorPass"
	MySQLUserDASUserKey                   = "mysql.user.dasUser"
	MySQLUserDASPassKey                   = "mysql.user.dasPass"
	// pmm
	PMMServerAddrKey                   = "pmm.server.addr"
	PMMServerUserKey                   = "pmm.server.user"
	PMMServerPassKey                   = "pmm.server.pass"
	PMMClientVersionKey                = "pmm.client.version"
	PMMClientInstallationPackageDirKey = "pmm.client.installationPackageDir"
)
