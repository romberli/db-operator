package mysql

import (
	"github.com/romberli/go-util/config"

	"github.com/romberli/db-operator/pkg/message"
)

func init() {
	initMySQLConfigDebugMessage()
	initMySQLConfigInfoMessage()
	initMySQLConfigErrorMessage()
}

const (
	// debug

	// info

	// error
	ErrMySQLNotValidConfigMySQLVersion                       = 402001
	ErrMySQLNotValidConfigMySQLParameterMaxConnections       = 402002
	ErrMySQLNotValidConfigMySQLParameterInnodbBufferPoolSize = 402003
	ErrMySQLNotValidConfigMySQLParameterInnodbIOCapacity     = 402004
	ErrMySQLNotValidConfigMySQLUser                          = 402005
	ErrMySQLNotValidConfigMySQLOperationTimeout              = 402006
)

func initMySQLConfigDebugMessage() {

}

func initMySQLConfigInfoMessage() {

}

func initMySQLConfigErrorMessage() {
	message.Messages[ErrMySQLNotValidConfigMySQLVersion] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLNotValidConfigMySQLVersion,
		"mysql.Config: version should be formatted as X.Y.Z and larger than %s, %s is not valid")
	message.Messages[ErrMySQLNotValidConfigMySQLParameterMaxConnections] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLNotValidConfigMySQLParameterMaxConnections,
		"mysql.Config: default max_connections should be in the range [%d, %d], %d is not valid")
	message.Messages[ErrMySQLNotValidConfigMySQLParameterInnodbBufferPoolSize] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLNotValidConfigMySQLParameterInnodbBufferPoolSize,
		"mysql.Config: default innodb_buffer_pool_size should be in the range [%d, %d], %d is not valid")
	message.Messages[ErrMySQLNotValidConfigMySQLParameterInnodbIOCapacity] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLNotValidConfigMySQLParameterInnodbIOCapacity,
		"mysql.Config: default innodb_io_capacity should be in the range [%d, %d], %d is not valid")
	message.Messages[ErrMySQLNotValidConfigMySQLUser] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLNotValidConfigMySQLUser,
		"mysql.Config: %s should not be empty")
	message.Messages[ErrMySQLNotValidConfigMySQLOperationTimeout] = config.NewErrMessage(message.DefaultMessageHeader, ErrMySQLNotValidConfigMySQLOperationTimeout,
		"mysql.Config: operation timeout should be in the range [%d, %d], %d is not valid")
}
