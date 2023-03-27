package mysql

import (
	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"github.com/romberli/db-operator/module/implement/mysql/parameter"
	"github.com/romberli/go-util/constant"
)

type InstallMySQL struct {
	Token            string
	Mode             mode.Mode
	Addrs            []string
	MySQLServerParam *parameter.MySQLServer
	PMMClientParam   *parameter.PMMClient
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
