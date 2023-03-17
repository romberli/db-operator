package config

import (
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

// InitDerivedConfig initializes the derived configuration
func InitDerivedConfig() error {
	// mysql.versionInt
	versionIntStr := strings.ReplaceAll(viper.GetString(MySQLVersionKey), constant.DotString, constant.EmptyString)
	versionInt, err := strconv.Atoi(versionIntStr)
	if err != nil {
		return errors.Trace(err)
	}
	viper.Set(MySQLVersionIntKey, versionInt)

	return nil
}
