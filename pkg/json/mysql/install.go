package mysql

import (
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

func NewInstallMySQL(token string, mode mode.Mode, addrs []string,
	mysqlServerParam *parameter.MySQLServer, pmmClientParam *parameter.PMMClient) *InstallMySQL {
	return newInstallMySQL(token, mode, addrs, mysqlServerParam, pmmClientParam)
}

func NewInstallMySQLWithDefault() *InstallMySQL {
	return newInstallMySQL(
		constant.EmptyString,
		mode.Standalone,
		[]string{},
		parameter.NewMySQLServerWithDefault(),
		parameter.NewPMMClientWithDefault(),
	)
}

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
