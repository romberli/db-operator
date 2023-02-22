package pmm

import (
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-util/config"
)

func init() {
	initPMMConfigDebugMessage()
	initPMMConfigInfoMessage()
	initPMMConfigErrorMessage()
}

const (
	// debug

	// info

	// error
	ErrNotValidConfigPMMServerAddr    = 403001
	ErrNotValidConfigPMMServerUser    = 403002
	ErrNotValidConfigPMMClientVersion = 403003
)

func initPMMConfigDebugMessage() {

}

func initPMMConfigInfoMessage() {

}

func initPMMConfigErrorMessage() {
	message.Messages[ErrNotValidConfigPMMServerAddr] = config.NewErrMessage(message.DefaultMessageHeader, ErrNotValidConfigPMMServerAddr,
		"pmm: server addr must be a url, %s is not valid")
	message.Messages[ErrNotValidConfigPMMServerUser] = config.NewErrMessage(message.DefaultMessageHeader, ErrNotValidConfigPMMServerUser,
		"pmm: server user should not be empty")
	message.Messages[ErrNotValidConfigPMMClientVersion] = config.NewErrMessage(message.DefaultMessageHeader, ErrNotValidConfigPMMClientVersion,
		"pmm: client version should be larger than %s, %s is not valid")
}
