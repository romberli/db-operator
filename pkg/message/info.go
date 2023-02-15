package message

import (
	"github.com/romberli/go-util/config"
)

const (
	// server
	InfoServerStart      = 200001
	InfoServerStop       = 200002
	InfoServerIsRunning  = 200003
	InfoServerNotRunning = 200004
)

func initInfoMessage() {
	// server
	Messages[InfoServerStart] = config.NewErrMessage(DefaultMessageHeader, InfoServerStart, "db-operator started. addr: %s, pid: %d, pid file: %s")
	Messages[InfoServerStop] = config.NewErrMessage(DefaultMessageHeader, InfoServerStop, "db-operator stopped successfully. pid: %d, pid file: %s")
	Messages[InfoServerIsRunning] = config.NewErrMessage(DefaultMessageHeader, InfoServerIsRunning, "db-operator is running. pid: %d")
	Messages[InfoServerNotRunning] = config.NewErrMessage(DefaultMessageHeader, InfoServerNotRunning, "db-operator is not running. pid: %d")
}
