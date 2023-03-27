package mysql

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
)

const (
	defaultIPLength         = 4
	defaultServerIDTemplate = "%d%03s%03s"
)

// GetConfig gets the configuration with templateName, templateContent and data
func GetConfig(templateName, templateContent string, data any) ([]byte, error) {
	t, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var buffer bytes.Buffer

	err = t.Execute(&buffer, data)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return buffer.Bytes(), nil
}

// GetServerID gets the server id from host ip and port number
func GetServerID(hostIP string, portNum int) (int, error) {
	// get server id
	ipList := strings.Split(hostIP, constant.DotString)
	if len(ipList) != defaultIPLength {
		return constant.ZeroInt, errors.Errorf("invalid host ip: %s", hostIP)
	}
	serverIDStr := fmt.Sprintf(defaultServerIDTemplate, portNum, ipList[2], ipList[3])
	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		return constant.ZeroInt, errors.Trace(err)
	}

	return serverID, nil
}
