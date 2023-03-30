package mysql

import (
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-util/config"
)

func init() {
	initDefaultEngineDebugMessage()
	initDefaultEngineInfoMessage()
	initDefaultEngineErrorMessage()
}

const (
	// debug

	// info
	InfoMySQLEngineInitInstance = 202201

	// error
	ErrMySQLEngineUpdateOperationDetail = 402201
)

func initDefaultEngineDebugMessage() {

}

func initDefaultEngineInfoMessage() {
	message.Messages[InfoMySQLEngineInitInstance] = config.NewErrMessage(message.DefaultMessageHeader, InfoMySQLEngineInitInstance,
		"mysql.Engine: init instance completed. operationID: %d, operationDetailID: %d, hostIP: %s, portNum: %d")
}

func initDefaultEngineErrorMessage() {
	message.Messages[ErrMySQLEngineUpdateOperationDetail] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLEngineUpdateOperationDetail,
		"mysql.Engine: update operation detail failed. operationID: %d, operationDetailID: %d, hostIP: %s, portNum: %d, status: %d")
}
