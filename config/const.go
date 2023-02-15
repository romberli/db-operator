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

// global constant
const (
	DefaultCommandName = "go-template-web"
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
)

// configuration constant
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
)
