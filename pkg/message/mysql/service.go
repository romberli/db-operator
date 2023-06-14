package mysql

import (
	"github.com/romberli/go-util/config"

	"github.com/romberli/db-operator/pkg/message"
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
	ErrMySQLServiceInstallMySQL           = 402101
	ErrMySQLServiceUpdateOperationHistory = 402102
)

func initMySQLServiceDebugMessage() {

}

func initMySQLServiceInfoMessage() {
	message.Messages[InfoMySQLServiceInstallMySQL] = config.NewErrMessage(message.DefaultMessageHeader, InfoMySQLServiceInstallMySQL,
		"mysql.Service: install mysql completed. version: %s, mode: %d, addrs: %s")
}

func initMySQLServiceErrorMessage() {
	message.Messages[ErrMySQLServiceInstallMySQL] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLServiceInstallMySQL,
		"mysql.Service: install mysql failed. version: %s, mode: %d, addrs: %s")
	message.Messages[ErrMySQLServiceUpdateOperationHistory] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLServiceUpdateOperationHistory,
		"mysql.Service: update operation history failed. operationID: %d, status: %d")
}
