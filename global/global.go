package global

import (
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/go-util/middleware/mysql"
	"github.com/romberli/log"
	"github.com/spf13/viper"
)

const (
	defaultMaxIdleConns        = 200
	defaultMaxIdleConnsPerHost = 50
)

var (
	DBOMySQLPool *mysql.Pool
)

// InitDBOMySQLPool initializes the global DBOMySQLPool
func InitDBOMySQLPool() (err error) {
	dbAddr := viper.GetString(config.DBDBOMySQLAddrKey)
	dbName := viper.GetString(config.DBDBOMySQLNameKey)
	dbUser := viper.GetString(config.DBDBOMySQLUserKey)
	dbPass := viper.GetString(config.DBDBOMySQLPassKey)
	maxConnections := viper.GetInt(config.DBPoolMaxConnectionsKey)
	initConnections := viper.GetInt(config.DBPoolInitConnectionsKey)
	maxIdleConnections := viper.GetInt(config.DBPoolMaxIdleConnectionsKey)
	maxIdleTime := viper.GetInt(config.DBPoolMaxIdleTimeKey)
	maxWaitTime := viper.GetInt(config.DBPoolMaxWaitTimeKey)
	maxRetryCount := viper.GetInt(config.DBPoolMaxRetryCountKey)
	keepAliveInterval := viper.GetInt(config.DBPoolKeepAliveIntervalKey)

	cfg := mysql.NewConfig(dbAddr, dbName, dbUser, dbPass)
	poolConfig := mysql.NewPoolConfigWithConfig(cfg, maxConnections, initConnections, maxIdleConnections, maxIdleTime,
		maxWaitTime, maxRetryCount, keepAliveInterval)
	log.Debugf("pool config: %v", poolConfig)
	DBOMySQLPool, err = mysql.NewPoolWithPoolConfig(poolConfig)
	if err != nil {
		return errors.Errorf("create dbo mysql pool failed. addr: %s, db: %s, user: %s. error:\n%s",
			dbAddr, dbName, dbUser, err.Error())
	}

	return nil
}
