package mysql

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/module/implement/mysql"
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/db-operator/pkg/resp"
	"github.com/romberli/go-util/linux"

	jsonmysql "github.com/romberli/db-operator/pkg/json/mysql"
	msgMySQL "github.com/romberli/db-operator/pkg/message/mysql"
)

const (
	installMySQLMessage = `{"version": "%s",  "mode": %d, "addrs": "%s", "message": "install mysql server completed"}`
)

// @Tags health
// @Summary install mysql server
// @Accept	application/json
// @Param	token	 			body string 			   true "token"
// @Param 	mode 				body int  				   true "mode"
// @Param   addrs 				body []string 			   true "addrs"
// @Param   mysqlServerParam	body parameter.MySQLServer true "mysql_server_param"
// @Param	pmmClientParam		body parameter.PMMClient   true "pmm_client_param"
// @Produce application/json
// @Success 200 {string} string "0"
// @Router	/api/v1/mysql/install [get]
func Install(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		resp.ResponseNOK(c, message.ErrGetRawData, errors.Trace(err))
		return
	}

	installMySQL := jsonmysql.NewInstallMySQLWithDefault()
	err = json.Unmarshal(data, &installMySQL)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, errors.Trace(err))
		return
	}
	mysqlVersion, err := version.NewVersion(installMySQL.MySQLServerParam.Version)
	if err != nil {
		resp.ResponseNOK(c, msgMySQL.ErrMySQLNotValidConfigMySQLVersion, errors.Trace(err))
		return
	}
	err = linux.SortAddrs(installMySQL.Addrs)
	if err != nil {
		resp.ResponseNOK(c, message.ErrSortAddrs, err, installMySQL.Addrs)
		return
	}

	e := mysql.NewEngineWithDefault(
		mysqlVersion,
		installMySQL.Mode,
		installMySQL.Addrs,
		installMySQL.MySQLServerParam,
		installMySQL.PMMClientParam,
	)

	jsonBytes, err := json.Marshal(installMySQL.Addrs)
	if err != nil {
		resp.ResponseNOK(c, message.ErrMarshalData, errors.Trace(err))
		return
	}
	jsonStr := string(jsonBytes)

	s := mysql.NewServiceWithDefault(e)
	err = s.Install()
	if err != nil {
		resp.ResponseNOK(c, msgMySQL.ErrMySQLServiceInstallMySQL, err,
			installMySQL.MySQLServerParam.Version, installMySQL.Mode, jsonStr)
		return
	}

	resp.ResponseOK(c, fmt.Sprintf(installMySQLMessage, installMySQL.MySQLServerParam.Version, installMySQL.Mode, jsonStr),
		msgMySQL.InfoMySQLServiceInstallMySQL, installMySQL.MySQLServerParam.Version, installMySQL.Mode, jsonStr)
}
