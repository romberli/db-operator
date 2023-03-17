package mysql

import (
	"encoding/json"
	"github.com/romberli/db-operator/module/implement/mysql"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/db-operator/pkg/resp"
)

const (
	tokenJSON     = "token"
	parameterJSON = "parameter"
)

// @Tags health
// @Summary install single instance mysql
// @Accept	application/json
// @Param	token	 	body string 			 true "token"
// @Param 	mode 		body int				 true "mode"
// @Param   addrs 		body string 			 true "addrs"
// @Param   parameter	body parameter.MySQLServer true "parameter"
// @Produce application/json
// @Success 200 {string} string "0"
// @Router	/api/v1/mysql/install/ [get]
func Install(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		resp.ResponseNOK(c, message.ErrGetRawData, errors.Trace(err))
	}

	data = jsonparser.Delete(data, tokenJSON)

	e := mysql.NewEngineWithDefault()
	err = json.Unmarshal(data, &e)
	if err != nil {
		resp.ResponseNOK(c, message.ErrUnmarshalRawData, errors.Trace(err))
	}

	s := mysql.NewService(e)
	err = s.Install()

}
