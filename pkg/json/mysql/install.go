package mysql

import (
	"encoding/json"

	"github.com/romberli/go-util/constant"

	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"github.com/romberli/db-operator/module/implement/mysql/parameter"
)

type InstallMySQL struct {
	Token            string                 `json:"token"`
	Mode             mode.Mode              `json:"mode"`
	Addrs            []string               `json:"addrs"`
	MySQLServerParam *parameter.MySQLServer `json:"mysql_server_param"`
	PMMClientParam   *parameter.PMMClient   `json:"pmm_client_param"`
}

// NewInstallMySQL returns a new *InstallMySQL
func NewInstallMySQL(token string, mode mode.Mode, addrs []string,
	mysqlServerParam *parameter.MySQLServer, pmmClientParam *parameter.PMMClient) *InstallMySQL {
	return newInstallMySQL(token, mode, addrs, mysqlServerParam, pmmClientParam)
}

// NewInstallMySQLWithDefault returns a new *InstallMySQL with default parameters
func NewInstallMySQLWithDefault() *InstallMySQL {
	return newInstallMySQL(
		constant.EmptyString,
		mode.Standalone,
		[]string{},
		parameter.NewMySQLServerWithDefault(),
		parameter.NewPMMClientWithDefault(),
	)
}

// newInstallMySQL returns a new *InstallMySQL
func newInstallMySQL(token string, mode mode.Mode, addrs []string,
	mysqlServerParam *parameter.MySQLServer, pmmClientParam *parameter.PMMClient) *InstallMySQL {
	return &InstallMySQL{
		Token:            token,
		Mode:             mode,
		Addrs:            addrs,
		MySQLServerParam: mysqlServerParam,
		PMMClientParam:   pmmClientParam,
	}
}

// Unmarshal unmarshals json data to *InstallMySQL
func (im *InstallMySQL) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, im)
	if err != nil {
		return err
	}

	im.MySQLServerParam.SetVersion(im.MySQLServerParam.Version)

	return nil
}
