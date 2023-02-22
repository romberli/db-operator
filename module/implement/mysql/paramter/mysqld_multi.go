package paramter

import (
	"github.com/romberli/db-operator/config"
	"github.com/romberli/go-util/constant"
	"github.com/spf13/viper"
)

const (
	DefaultMySQLDMultiLogBaseName = "mysqld_multi/mysqld_multi.log"
	DefaultMySQLDMultiUser        = "mysqld_multi"
	DefaultMySQLDMultiPass        = "mysqld_multi"
)

type MySQLDMulti struct {
	Log  string `json:"log" config:"log"`
	User string `json:"user" config:"user"`
	Pass string `json:"pass" config:"pass"`
}

// NewMySQLDMulti returns a new *MySQLDMulti
func NewMySQLDMulti(log, user, pass string) *MySQLDMulti {
	return &MySQLDMulti{
		Log:  log,
		User: user,
		Pass: pass,
	}
}

// NewMySQLDMultiWithDefault returns a new *MySQLDMulti with default values
func NewMySQLDMultiWithDefault() *MySQLDMulti {
	return &MySQLDMulti{
		Log:  constant.DefaultRandomString,
		User: viper.GetString(config.MySQLUserMySQLDMultiUserKey),
		Pass: viper.GetString(config.MySQLUserMySQLDMultiPassKey),
	}
}

// GetLog returns the log
func (mm *MySQLDMulti) GetLog() string {
	return mm.Log
}

// GetUser returns the user
func (mm *MySQLDMulti) GetUser() string {
	return mm.User
}

// GetPass returns the pass
func (mm *MySQLDMulti) GetPass() string {
	return mm.Pass
}
