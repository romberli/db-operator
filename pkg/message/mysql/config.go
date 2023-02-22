package mysql

import (
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-util/config"
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
	ErrNotValidConfigMySQLVersion                       = 402001
	ErrNotValidConfigMySQLParameterMaxConnections       = 402002
	ErrNotValidConfigMySQLParameterInnodbBufferPoolSize = 402003
	ErrNotValidConfigMySQLParameterInnodbIOCapacity     = 402004
	ErrNotValidConfigMySQLUser                          = 402005
)

func initMySQLConfigDebugMessage() {

}

func initMySQLConfigInfoMessage() {

}

func initMySQLConfigErrorMessage() {
	message.Messages[ErrNotValidConfigMySQLVersion] = config.NewErrMessage(message.DefaultMessageHeader, ErrNotValidConfigMySQLVersion,
		"mysql: version should be larger than %s, %s is not valid")
	message.Messages[ErrNotValidConfigMySQLParameterMaxConnections] = config.NewErrMessage(message.DefaultMessageHeader, ErrNotValidConfigMySQLParameterMaxConnections,
		"mysql: default max_connections should be in the range [%d, %d], %d is not valid")
	message.Messages[ErrNotValidConfigMySQLParameterInnodbBufferPoolSize] = config.NewErrMessage(message.DefaultMessageHeader, ErrNotValidConfigMySQLParameterInnodbBufferPoolSize,
		"mysql: default innodb_buffer_pool_size should be in the range [%d, %d], %d is not valid")
	message.Messages[ErrNotValidConfigMySQLParameterInnodbIOCapacity] = config.NewErrMessage(message.DefaultMessageHeader, ErrNotValidConfigMySQLParameterInnodbIOCapacity,
		"mysql: default innodb_io_capacity should be in the range [%d, %d], %d is not valid")
	message.Messages[ErrNotValidConfigMySQLUser] = config.NewErrMessage(message.DefaultMessageHeader, ErrNotValidConfigMySQLUser,
		"mysql: %s should not be empty")
}
