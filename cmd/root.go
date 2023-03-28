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

package cmd

import (
	"fmt"
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const defaultConfigFileType = "yaml"

var (
	// common
	baseDir string
	cfgFile string
	// daemon
	daemon    bool
	daemonStr string
	// log
	logFileName           string
	logLevel              string
	logFormat             string
	logMaxSize            int
	logMaxDays            int
	logMaxBackups         int
	logRotateOnStartupStr string
	// server
	serverAddr                      string
	serverPid                       int
	serverPidFile                   string
	serverReadTimeout               int
	serverWriteTimeout              int
	serverPProfEnabledStr           string
	serverRouterAlternativeBasePath string
	serverRouterAlternativeBodyPath string
	serverRouterHTTPErrorCode       int
	// database
	dbDBOMySQLAddr           string
	dbDBOMySQLName           string
	dbDBOMySQLUser           string
	dbDBOMySQLPass           string
	dbPoolMaxConnections     int
	dbPoolInitConnections    int
	dbPoolMaxIdleConnections int
	dbPoolMaxIdleTime        int
	dbPoolMaxWaitTime        int
	dbPoolMaxRetryCount      int
	dbPoolKeepAliveInterval  int
	// mysql
	mysqlVersion                       string
	mysqlInstallationPackageDir        string
	mysqlInstallationTemporaryDir      string
	mysqlParameterMaxConnections       int
	mysqlParameterInnodbBufferPoolSize int
	mysqlParameterInnodbIOCapacity     int
	mysqlUserOSUser                    string
	mysqlUserOSPass                    string
	mysqlUserRootPass                  string
	mysqlUserAdminUser                 string
	mysqlUserAdminPass                 string
	mysqlUserMySQLDMultiUser           string
	mysqlUserMySQLDMultiPass           string
	mysqlUserReplicationUser           string
	mysqlUserReplicationPass           string
	mysqlUserMonitorUser               string
	mysqlUserMonitorPass               string
	mysqlUserDASUser                   string
	mysqlUserDASPass                   string
	mysqlOperationTimeout              int
	// pmm
	pmmServerAddr                   string
	pmmServerUser                   string
	pmmServerPass                   string
	pmmClientVersion                string
	pmmClientInstallationPackageDir string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "db-operator",
	Short: "db-operator",
	Long:  `db-operator is a tmpl of golang web server`,
	Run: func(cmd *cobra.Command, args []string) {
		// if no subcommand is set, it will print help information.
		if len(args) == constant.ZeroInt {
			err := cmd.Help()
			if err != nil {
				fmt.Println(fmt.Sprintf(constant.LogWithStackString, message.NewMessage(message.ErrPrintHelpInfo, errors.Trace(err))))
				os.Exit(constant.DefaultAbnormalExitCode)
			}

			os.Exit(constant.DefaultNormalExitCode)
		}

		// init config
		err := initConfig()
		if err != nil {
			fmt.Println(fmt.Sprintf(constant.LogWithStackString, message.NewMessage(message.ErrInitConfig, err)))
			os.Exit(constant.DefaultAbnormalExitCode)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(fmt.Sprintf(constant.LogWithStackString, errors.Trace(err)))
		os.Exit(constant.DefaultAbnormalExitCode)
	}
}

func init() {
	// set usage tmpl
	rootCmd.SetUsageTemplate(UsageTemplateWithoutDefault())

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// config
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", constant.DefaultRandomString, "config file path")
	// daemon
	rootCmd.PersistentFlags().StringVar(&daemonStr, "daemon", constant.DefaultRandomString, fmt.Sprintf("whether run in background as a daemon(default: %s)", constant.FalseString))
	// log
	rootCmd.PersistentFlags().StringVar(&logFileName, "log-file", constant.DefaultRandomString, fmt.Sprintf("specify the log file name(default: %s)", filepath.Join(config.DefaultLogDir, log.DefaultLogFileName)))
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", constant.DefaultRandomString, fmt.Sprintf("specify the log level(default: %s)", log.DefaultLogLevel))
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", constant.DefaultRandomString, fmt.Sprintf("specify the log format(default: %s)", log.DefaultLogFormat))
	rootCmd.PersistentFlags().IntVar(&logMaxSize, "log-max-size", constant.DefaultRandomInt, fmt.Sprintf("specify the log file max size(default: %d)", log.DefaultLogMaxSize))
	rootCmd.PersistentFlags().IntVar(&logMaxDays, "log-max-days", constant.DefaultRandomInt, fmt.Sprintf("specify the log file max days(default: %d)", log.DefaultLogMaxDays))
	rootCmd.PersistentFlags().IntVar(&logMaxBackups, "log-max-backups", constant.DefaultRandomInt, fmt.Sprintf("specify the log file max backups(default: %d)", log.DefaultLogMaxBackups))
	rootCmd.PersistentFlags().StringVar(&logRotateOnStartupStr, "log-rotate-on-startup", constant.DefaultRandomString, fmt.Sprintf("specify if rotating the log file on startup(default: %s)", constant.FalseString))
	// server
	rootCmd.PersistentFlags().StringVar(&serverAddr, "server-addr", constant.DefaultRandomString, fmt.Sprintf("specify the server addr(default: %s)", config.DefaultServerAddr))
	rootCmd.PersistentFlags().StringVar(&serverPidFile, "server-pid-file", constant.DefaultRandomString, fmt.Sprintf("specify the server pid file path(default: %s)", filepath.Join(config.DefaultBaseDir, fmt.Sprintf("%s.pid", config.DefaultCommandName))))
	rootCmd.PersistentFlags().IntVar(&serverReadTimeout, "server-read-timeout", constant.DefaultRandomInt, fmt.Sprintf("specify the read timeout in seconds of http request(default: %d)", config.DefaultServerReadTimeout))
	rootCmd.PersistentFlags().IntVar(&serverWriteTimeout, "server-write-timeout", constant.DefaultRandomInt, fmt.Sprintf("specify the write timeout in seconds of http request(default: %d)", config.DefaultServerWriteTimeout))
	rootCmd.PersistentFlags().StringVar(&serverPProfEnabledStr, "server-pprof-enabled", constant.DefaultRandomString, fmt.Sprintf("specify if enable the pprof(default: %s)", constant.FalseString))
	rootCmd.PersistentFlags().StringVar(&serverRouterAlternativeBasePath, "server-router-alternative-base-path", constant.DefaultRandomString, fmt.Sprintf("specify the alternative base path(default: %s)", config.DefaultServerRouterAlternativeBasePath))
	rootCmd.PersistentFlags().StringVar(&serverRouterAlternativeBodyPath, "server-router-alternative-body-path", constant.DefaultRandomString, fmt.Sprintf("specify the alternative body path of the json body of the http request(default: %s)", config.DefaultServerRouterAlternativeBodyPath))
	rootCmd.PersistentFlags().IntVar(&serverRouterHTTPErrorCode, "server-router-http-error-code", constant.DefaultRandomInt, fmt.Sprintf("specify the http return code when the server encountered an error(default: %d)", config.DefaultServerRouterHTTPErrorCode))
	//  database
	rootCmd.PersistentFlags().StringVar(&dbDBOMySQLAddr, "db-dbo-mysql-addr", constant.DefaultRandomString, fmt.Sprintf("specify dbo database address(format: host:port)(default: %s)", constant.DefaultMySQLAddr))
	rootCmd.PersistentFlags().StringVar(&dbDBOMySQLName, "db-dbo-mysql-name", constant.DefaultRandomString, fmt.Sprintf("specify dbo database name(default: %s)", config.DefaultDBName))
	rootCmd.PersistentFlags().StringVar(&dbDBOMySQLUser, "db-dbo-mysql-user", constant.DefaultRandomString, fmt.Sprintf("specify dbo database user name(default: %s)", config.DefaultDBUser))
	rootCmd.PersistentFlags().StringVar(&dbDBOMySQLPass, "db-dbo-mysql-pass", constant.DefaultRandomString, fmt.Sprintf("specify dbo database user password(default: %s)", config.DefaultDBPass))
	rootCmd.PersistentFlags().IntVar(&dbPoolMaxConnections, "db-pool-max-connections", constant.DefaultRandomInt, fmt.Sprintf("specify max connections of the connection pool(default: %d)", mysql.DefaultMaxConnections))
	rootCmd.PersistentFlags().IntVar(&dbPoolInitConnections, "db-pool-init-connections", constant.DefaultRandomInt, fmt.Sprintf("specify initial connections of the connection pool(default: %d)", mysql.DefaultMaxIdleConnections))
	rootCmd.PersistentFlags().IntVar(&dbPoolMaxIdleConnections, "db-pool-max-idle-connections", constant.DefaultRandomInt, fmt.Sprintf("specify max idle connections of the connection pool(default: %d)", mysql.DefaultMaxIdleConnections))
	rootCmd.PersistentFlags().IntVar(&dbPoolMaxIdleTime, "db-pool-max-idle-time", constant.DefaultRandomInt, fmt.Sprintf("specify max idle time of connections of the connection pool, (default: %d, unit: seconds)", mysql.DefaultMaxIdleTime))
	rootCmd.PersistentFlags().IntVar(&dbPoolMaxWaitTime, "db-pool-max-wait-time", constant.DefaultRandomInt, fmt.Sprintf("specify max wait time of getting a the connection from pool, (default: %d, unit: seconds)", mysql.DefaultMaxWaitTime))
	rootCmd.PersistentFlags().IntVar(&dbPoolMaxRetryCount, "db-pool-max-retry-count", constant.DefaultRandomInt, fmt.Sprintf("specify max retry count of getting a the connection from pool, (default: %d)", mysql.DefaultMaxRetryCount))
	rootCmd.PersistentFlags().IntVar(&dbPoolKeepAliveInterval, "db-pool-keep-alive-interval", constant.DefaultRandomInt, fmt.Sprintf("specify keep alive interval of connections of the connection pool(default: %d, unit: seconds)", mysql.DefaultKeepAliveInterval))
	// mysql
	rootCmd.PersistentFlags().StringVar(&mysqlVersion, "mysql-version", constant.DefaultRandomString, fmt.Sprintf("specify the default mysql version(default: %s)", config.DefaultMySQLVersion))
	rootCmd.PersistentFlags().StringVar(&mysqlInstallationPackageDir, "mysql-installation-package-dir", constant.DefaultRandomString, fmt.Sprintf("specify the mysql binary installation package directory(default: %s)", config.DefaultMySQLInstallationPackageDir))
	rootCmd.PersistentFlags().StringVar(&mysqlInstallationTemporaryDir, "mysql-installation-temporary-dir", constant.DefaultRandomString, fmt.Sprintf("specify the temporary directory for mysql installation(default: %s)", config.DefaultMySQLInstallationTemporaryDir))
	rootCmd.PersistentFlags().IntVar(&mysqlParameterMaxConnections, "mysql-parameter-max-connections", constant.DefaultRandomInt, fmt.Sprintf("specify the default max connections(default: %d)", config.DefaultMySQLParameterMaxConnections))
	rootCmd.PersistentFlags().IntVar(&mysqlParameterInnodbBufferPoolSize, "mysql-parameter-innodb-buffer-pool-size", constant.DefaultRandomInt, fmt.Sprintf("specify the default innodb buffer pool size(default: %d)", config.DefaultMySQLParameterInnodbBufferPoolSize))
	rootCmd.PersistentFlags().IntVar(&mysqlParameterInnodbIOCapacity, "mysql-parameter-innodb-io-capacity", constant.DefaultRandomInt, fmt.Sprintf("specify the default innodb io capacity(default: %d)", config.DefaultMySQLParameterInnodbIOCapacity))
	rootCmd.PersistentFlags().StringVar(&mysqlUserOSUser, "mysql-user-os-user", constant.DefaultRandomString, fmt.Sprintf("specify the default os user(default: %s)", config.DefaultMySQLUserOSUser))
	rootCmd.PersistentFlags().StringVar(&mysqlUserOSPass, "mysql-user-os-pass", constant.DefaultRandomString, fmt.Sprintf("specify the default os password(default: %s)", config.DefaultMySQLUserOSPass))
	rootCmd.PersistentFlags().StringVar(&mysqlUserRootPass, "mysql-user-root-pass", constant.DefaultRandomString, fmt.Sprintf("specify the default root password(default: %s)", config.DefaultMySQLUserRootPass))
	rootCmd.PersistentFlags().StringVar(&mysqlUserAdminUser, "mysql:user-admin-user", constant.DefaultRandomString, fmt.Sprintf("specify the default admin user(default: %s)", config.DefaultMySQLUserAdminUser))
	rootCmd.PersistentFlags().StringVar(&mysqlUserAdminPass, "mysql:user-admin-pass", constant.DefaultRandomString, fmt.Sprintf("specify the default admin password(default: %s)", config.DefaultMySQLUserAdminPass))
	rootCmd.PersistentFlags().StringVar(&mysqlUserMySQLDMultiUser, "mysql:user-mysqld-multi-user", constant.DefaultRandomString, fmt.Sprintf("specify the default mysqld multi user(default: %s)", config.DefaultMySQLUserMySQLDMultiUser))
	rootCmd.PersistentFlags().StringVar(&mysqlUserMySQLDMultiPass, "mysql:user-mysqld-multi-pass", constant.DefaultRandomString, fmt.Sprintf("specify the default mysqld multi password(default: %s)", config.DefaultMySQLUserMySQLDMultiPass))
	rootCmd.PersistentFlags().StringVar(&mysqlUserReplicationUser, "mysql:user-replication-user", constant.DefaultRandomString, fmt.Sprintf("specify the default replication user(default: %s)", config.DefaultMySQLUserReplicationUser))
	rootCmd.PersistentFlags().StringVar(&mysqlUserReplicationPass, "mysql:user-replication-pass", constant.DefaultRandomString, fmt.Sprintf("specify the default replication password(default: %s)", config.DefaultMySQLUserReplicationPass))
	rootCmd.PersistentFlags().StringVar(&mysqlUserMonitorUser, "mysql:user-monitor-user", constant.DefaultRandomString, fmt.Sprintf("specify the default monitor user(default: %s)", config.DefaultMySQLUserMonitorUser))
	rootCmd.PersistentFlags().StringVar(&mysqlUserMonitorPass, "mysql:user-monitor-pass", constant.DefaultRandomString, fmt.Sprintf("specify the default monitor password(default: %s)", config.DefaultMySQLUserMonitorPass))
	rootCmd.PersistentFlags().StringVar(&mysqlUserDASUser, "mysql:user-das-user", constant.DefaultRandomString, fmt.Sprintf("specify the default das user(default: %s)", config.DefaultMySQLUserDASUser))
	rootCmd.PersistentFlags().StringVar(&mysqlUserDASPass, "mysql:user-das-pass", constant.DefaultRandomString, fmt.Sprintf("specify the default das password(default: %s)", config.DefaultMySQLUserDASPass))
	rootCmd.PersistentFlags().IntVar(&mysqlOperationTimeout, "mysql-operation-timeout", constant.DefaultRandomInt, fmt.Sprintf("specify the default mysql operation timeout(default: %d, unit: seconds)", config.DefaultMySQLOperationTimeout))
	// pmm
	rootCmd.PersistentFlags().StringVar(&pmmServerAddr, "pmm-server-addr", constant.DefaultRandomString, fmt.Sprintf("specify the pmm server address(default: %s)", config.DefaultPMMServerAddr))
	rootCmd.PersistentFlags().StringVar(&pmmServerUser, "pmm-server-user", constant.DefaultRandomString, fmt.Sprintf("specify the pmm server user(default: %s)", config.DefaultPMMServerUser))
	rootCmd.PersistentFlags().StringVar(&pmmServerPass, "pmm-server-pass", constant.DefaultRandomString, fmt.Sprintf("specify the pmm server password(default: %s)", config.DefaultPMMServerPass))
	rootCmd.PersistentFlags().StringVar(&pmmClientVersion, "pmm-client-version", constant.DefaultRandomString, fmt.Sprintf("specify the pmm client version(default: %s)", config.DefaultPMMClientVersion))
	rootCmd.PersistentFlags().StringVar(&pmmClientInstallationPackageDir, "pmm-client-installation-package-dir", constant.DefaultRandomString, fmt.Sprintf("specify the pmm client binary installation package dir(default: %s)", config.DefaultPMMClientInstallationPackageDir))
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() error {
	var err error

	// init default config
	err = initDefaultConfig()
	if err != nil {
		return message.NewMessage(message.ErrInitDefaultConfig, err.Error())
	}

	// read config with config file
	err = ReadConfigFile()
	if err != nil {
		return message.NewMessage(message.ErrInitDefaultConfig, err)
	}

	// override config with command line arguments
	err = OverrideConfigByCLI()
	if err != nil {
		return message.NewMessage(message.ErrOverrideCommandLineArgs, err)
	}

	// init derived config
	err = config.InitDerivedConfig()
	if err != nil {
		return message.NewMessage(message.ErrInitDerivedConfig, err)
	}

	// init log
	fileName := viper.GetString(config.LogFileNameKey)
	level := viper.GetString(config.LogLevelKey)
	format := viper.GetString(config.LogFormatKey)
	maxSize := viper.GetInt(config.LogMaxSizeKey)
	maxDays := viper.GetInt(config.LogMaxDaysKey)
	maxBackups := viper.GetInt(config.LogMaxBackupsKey)

	fileNameAbs := fileName
	isAbs := filepath.IsAbs(fileName)
	if !isAbs {
		fileNameAbs, err = filepath.Abs(fileName)
		if err != nil {
			return message.NewMessage(message.ErrAbsoluteLogFilePath, errors.Trace(err), fileName)
		}
	}
	_, _, err = log.InitFileLogger(fileNameAbs, level, format, maxSize, maxDays, maxBackups)
	if err != nil {
		return message.NewMessage(message.ErrInitLogger, err)
	}

	log.SetDisableDoubleQuotes(true)
	log.SetDisableEscape(true)

	return nil
}

// initDefaultConfig initiate default configuration
func initDefaultConfig() (err error) {
	// get base dir
	baseDir, err = filepath.Abs(config.DefaultBaseDir)
	if err != nil {
		return message.NewMessage(message.ErrBaseDir, errors.Trace(err), config.DefaultCommandName)
	}
	// set default config value
	config.SetDefaultConfig(baseDir)
	err = config.ValidateConfig()
	if err != nil {
		return err
	}

	return nil
}

// ReadConfigFile read configuration from config file, it will override the init configuration
func ReadConfigFile() (err error) {
	if cfgFile != constant.EmptyString && cfgFile != constant.DefaultRandomString {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType(defaultConfigFileType)
		err = viper.ReadInConfig()
		if err != nil {
			return errors.Trace(err)
		}
		err = config.ValidateConfig()
		if err != nil {
			return message.NewMessage(message.ErrValidateConfig, err)
		}
	}

	return nil
}

// UsageTemplateWithoutDefault returns a usage tmpl which does not contain default part
func UsageTemplateWithoutDefault() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsagesWithoutDefault | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsagesWithoutDefault | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}
