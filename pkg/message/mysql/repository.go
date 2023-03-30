package mysql

import (
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-util/config"
)

func init() {
	initMySQLRepositoryDebugMessage()
	initMySQLRepositoryInfoMessage()
	initMySQLRepositoryErrorMessage()
}

const (
	// debug

	// info

	// error
	ErrMySQLRepositoryGetLock     = 402301
	ErrMySQLRepositoryReleaseLock = 402302
)

func initMySQLRepositoryDebugMessage() {

}

func initMySQLRepositoryInfoMessage() {

}

func initMySQLRepositoryErrorMessage() {
	message.Messages[ErrMySQLRepositoryGetLock] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLRepositoryGetLock,
		"mysql.Repository: get lock failed. operation_id: %d, addrs: %s")
	message.Messages[ErrMySQLRepositoryReleaseLock] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLRepositoryReleaseLock,
		"mysql.Repository: release lock failed. operation_id: %d")
}
