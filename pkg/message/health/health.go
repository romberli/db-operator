package health

import (
	"github.com/romberli/go-template-web/pkg/message"
	"github.com/romberli/go-util/config"
)

func init() {
	initHealthDebugMessage()
	initHealthInfoMessage()
	initHealthErrorMessage()
}

const (
	// info
	InfoHealthStatus = 209001
	InfoHealthPing   = 200102
)

func initHealthDebugMessage() {

}

func initHealthInfoMessage() {
	message.Messages[InfoHealthStatus] = config.NewErrMessage(message.DefaultMessageHeader, InfoHealthStatus, "health: check status completed")
	message.Messages[InfoHealthPing] = config.NewErrMessage(message.DefaultMessageHeader, InfoHealthPing, "health: ping completed")
}

func initHealthErrorMessage() {

}
