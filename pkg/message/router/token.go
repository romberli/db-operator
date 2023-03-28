package router

import (
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-util/config"
)

func init() {
	initRouterDebugMessage()
	initRouterInfoMessage()
	initRouterErrorMessage()
}

const (
	// debug

	// info

	// error
	ErrRouterGetHandlerFunc = 408001
	ErrRouterValidateToken  = 408002
)

func initRouterDebugMessage() {

}

func initRouterInfoMessage() {

}

func initRouterErrorMessage() {
	message.Messages[ErrRouterGetHandlerFunc] = config.NewErrMessage(message.DefaultMessageHeader, ErrRouterGetHandlerFunc, "router: get token handler func failed")
	message.Messages[ErrRouterValidateToken] = config.NewErrMessage(message.DefaultMessageHeader, ErrRouterValidateToken, "router: validate token failed. token: %s, client ip: %s")
}
