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
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/pingcap/errors"
	"github.com/romberli/go-multierror"
	"github.com/romberli/go-template-web/pkg/message"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/spf13/cast"
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
	viper.SetDefault(DaemonKey, DefaultDaemon)
	// log
	defaultLogFile := filepath.Join(baseDir, DefaultLogDir, log.DefaultLogFileName)
	viper.SetDefault(LogFileNameKey, defaultLogFile)
	viper.SetDefault(LogLevelKey, log.DefaultLogLevel)
	viper.SetDefault(LogFormatKey, log.DefaultLogFormat)
	viper.SetDefault(LogMaxSizeKey, log.DefaultLogMaxSize)
	viper.SetDefault(LogMaxDaysKey, log.DefaultLogMaxDays)
	viper.SetDefault(LogMaxBackupsKey, log.DefaultLogMaxBackups)
	viper.SetDefault(LogRotateOnStartupKey, DefaultRotateOnStartup)
	// server
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
