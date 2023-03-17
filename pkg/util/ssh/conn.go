package ssh

import (
	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/linux"
	"strings"
)

const (
	getOSVersionCommand = "/usr/bin/cat /etc/redhat-release"
	getArchCommand      = "/usr/bin/uname -m"

	centOS     = "CentOS Linux"
	almaLinux  = "AlmaLinux"
	rockyLinux = "Rocky Linux"
)

type Conn struct {
	*linux.SSHConn
}

// NewConn returns a new *Conn
func NewConn(conn *linux.SSHConn) *Conn {
	return newConn(conn)
}

// newConn returns a new *Conn
func newConn(conn *linux.SSHConn) *Conn {
	return &Conn{
		SSHConn: conn,
	}
}

// GetOSVersion returns the os version of the host
func (c *Conn) GetOSVersion() (*version.Version, error) {
	output, err := c.ExecuteCommand(getOSVersionCommand)
	if err != nil {
		return nil, err
	}
	if len(output) == constant.ZeroInt {
		return nil, errors.New("get os version failed")
	}

	versionList := strings.Split(output, constant.SpaceString)
	if len(versionList) < constant.FourInt {
		return nil, errors.Errorf("os version is not valid. os_version: %s", output)
	}

	var osVersionStr string
	if strings.Contains(output, centOS) || strings.Contains(output, rockyLinux) {
		osVersionStr = versionList[constant.ThreeInt]
	} else if strings.Contains(output, almaLinux) {
		osVersionStr = versionList[constant.TwoInt]
	} else {
		return nil, errors.Errorf("os version must be one of [CentOS, AlmaLinux, Rocky], %s is not valid", output)
	}

	osVersion, err := version.NewVersion(osVersionStr)
	if err != nil {
		return nil, err
	}

	return osVersion, nil

}

// GetArch returns the arch of the host
func (c *Conn) GetArch() (string, error) {
	output, err := c.ExecuteCommand(getArchCommand)
	if err != nil {
		return "", err
	}
	if len(output) == constant.ZeroInt {
		return "", errors.New("get arch failed")
	}

	return output, nil
}
