package parameter

import (
	"fmt"

	"github.com/romberli/db-operator/module/implement/mysql/parameter/tmpl"
	"github.com/romberli/db-operator/pkg/util/mysql"
	"github.com/romberli/go-util/common"
)

const (
	commonTemplateName = "common"

	dirBaseTemplate = "%s/mysql%d"
)

type Common struct {
	PortNum         int    `json:"port_num" config:"port_num"`
	DataDirBaseName string `json:"data_dir_base_name" config:"data_dir_base_name"`
	DataDirBase     string `json:"data_dir_base" config:"data_dir_base"`
	ClientUser      string `json:"client_user" config:"client_user"`
	ClientPass      string `json:"client_pass" config:"client_pass"`
	MySQLDMultiUser string `json:"mysqld_multi_user" config:"mysqld_multi_user"`
	MySQLDMultiPass string `json:"mysqld_multi_pass" config:"mysqld_multi_pass"`
}

// NewCommon returns a new *Common
func NewCommon(portNum int, dataDirBaseName, clientUser, clientPass, mysqldMultiUser, mysqldMultiPass string) *Common {
	return &Common{
		PortNum:         portNum,
		DataDirBaseName: dataDirBaseName,
		DataDirBase:     fmt.Sprintf(dirBaseTemplate, dataDirBaseName, portNum),
		ClientUser:      clientUser,
		ClientPass:      clientPass,
		MySQLDMultiUser: mysqldMultiUser,
		MySQLDMultiPass: mysqldMultiPass,
	}
}

// Set sets Common with given fields, key is the field name and value is the relevant value of the key
func (c *Common) Set(fields map[string]interface{}) error {
	for fieldName, fieldValue := range fields {
		err := common.SetValueOfStruct(c, fieldName, fieldValue)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetConfig() gets the configuration of Common
func (c *Common) GetConfig() ([]byte, error) {
	return mysql.GetConfig(commonTemplateName, tmpl.Common, c)
}
