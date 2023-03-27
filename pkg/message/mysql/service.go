package mysql

import (
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-util/config"
)

func init() {
	initMySQLServiceDebugMessage()
	initMySQLServiceInfoMessage()
	initMySQLServiceErrorMessage()
}

const (
	// debug

	// info
	InfoMySQLServiceInstallMySQL = 202101

	// error
	ErrMySQLServiceInstallMySQL = 402101
)

func initMySQLServiceDebugMessage() {

}

func initMySQLServiceInfoMessage() {
	message.Messages[InfoMySQLServiceInstallMySQL] = config.NewErrMessage(message.DefaultMessageHeader, InfoMySQLServiceInstallMySQL,
		"mysql.Service: install mysql completed. version: %s, mode: %d, addrs: %s")
}

func initMySQLServiceErrorMessage() {

}
