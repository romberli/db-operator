package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/module/implement/mysql/parameter"
	"github.com/romberli/db-operator/pkg/util/ssh"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"github.com/spf13/viper"
)

const (
	pmmClientInstallationPackageNameTemplate = "pmm2-client-%s-6.el7.x86_64.rpm"

	checkPMMClientCommand = "/usr/local/bin/pmm-admin --version"
	pmmAdminNotFound      = "bash: line 1: /usr/local/bin/pmm-admin: No such file or directory"

	pmmClientInstallCommandTemplate            = "/usr/bin/yum localinstall -y %s"
	pmmClientConfigureServerCommandTemplate    = "/usr/local/bin/pmm-admin config --server-insecure-tls --server-url=http://%s:%s@%s"
	pmmClientStartClientCommandTemplate        = "/usr/bin/systemctl start pmm-agent"
	pmmClientCheckConfigurationCommandTemplate = "/usr/local/bin/pmm-admin list"
	pmmClientCheckServiceCommandTemplate       = "/usr/local/bin/pmm-admin list | grep ^MySQL | grep %d | grep -v grep | wc -l"
	pmmClientAddServiceCommandTemplateV1       = "/usr/local/bin/pmm-admin add mysql --host=127.0.0.1 --port=%d --username=%s --password=%s %s"
	pmmClientAddServiceCommandTemplateV2       = "/usr/local/bin/pmm-admin add mysql --host=127.0.0.1 --port=%d --username=%s --password=%s --replication-set=%s %s"

	pmmClientServiceNameTemplate = "%s-%d"
	pmmClientNodeExporterOutput  = "node_exporter"
)

type PMMExecutor struct {
	sshConn   *ssh.Conn
	hostIP    string
	portNum   int
	pmmClient *parameter.PMMClient
}

// NewPMMExecutor returns a new *PMMExecutor
func NewPMMExecutor(sshConn *ssh.Conn, hostIP string, portNum int, pmmClient *parameter.PMMClient) *PMMExecutor {
	return newPMMExecutor(sshConn, hostIP, portNum, pmmClient)
}

// newPMMExecutor returns a new *PMMExecutor
func newPMMExecutor(sshConn *ssh.Conn, hostIP string, portNum int, pmmClient *parameter.PMMClient) *PMMExecutor {
	return &PMMExecutor{
		sshConn:   sshConn,
		hostIP:    hostIP,
		portNum:   portNum,
		pmmClient: pmmClient,
	}
}

// Init initializes pmm client on the host
func (pe *PMMExecutor) Init() error {
	// Install pmm client binary
	installed, err := pe.CheckPMMClient()
	if err != nil {
		return err
	}
	// get arch
	arch, err := pe.sshConn.GetArch()
	if err != nil {
		return err
	}

	var configured bool

	if !installed {
		if arch == constant.X64Arch {
			return errors.Errorf("installing pmm client only supports %s arch, %s is not valid", constant.X64Arch, arch)
		}

		err = pe.Install()
		if err != nil {
			return err
		}
		configured = true
	}

	// check if pmm server is configured
	configured, err = pe.CheckConfiguration()
	if err != nil {
		return err
	}

	if !configured {
		// configure pmm server
		err = pe.ConfigureServer()
		if err != nil {
			return err
		}
		// start pmm client
		err = pe.StartClient()
		if err != nil {
			return err
		}
	}
	// check if the service exists
	exists, err := pe.CheckServiceExists()
	if err != nil {
		return err
	}
	if !exists {
		// add service
		err = pe.AddService()
		if err != nil {
			return err
		}
	}

	return nil
}

// CheckPMMClient checks if pmm client is installed
func (pe *PMMExecutor) CheckPMMClient() (bool, error) {
	_, err := pe.sshConn.ExecuteCommand(checkPMMClientCommand)
	if err != nil {
		e := err.Error()
		if strings.Contains(e, pmmAdminNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Install installs pmm client to the host
func (pe *PMMExecutor) Install() error {
	pmmClientInstallationPackageName := pe.getInstallationPackageName()
	fileSource := viper.GetString(config.PMMClientInstallationPackageDirKey) + pmmClientInstallationPackageName
	fileDest := viper.GetString(config.MySQLInstallationPackageDirKey) + pmmClientInstallationPackageName
	err := pe.sshConn.CopySingleFileToRemote(fileSource, fileDest)
	if err != nil {
		return err
	}

	output, err := pe.sshConn.ExecuteCommand(fmt.Sprintf(pmmClientInstallCommandTemplate, fileDest))
	if err != nil {
		return err
	}
	log.Debugf("PMMExecutor.Install() pmm client output:\n%s", output)

	return nil
}

// CheckConfiguration checks if pmm is configured
func (pe *PMMExecutor) CheckConfiguration() (bool, error) {
	output, err := pe.sshConn.ExecuteCommand(pmmClientCheckConfigurationCommandTemplate)
	if err != nil {
		return false, err
	}

	return strings.Contains(output, pmmClientNodeExporterOutput), nil
}

// ConfigureServer configures pmm server
func (pe *PMMExecutor) ConfigureServer() error {
	sql := fmt.Sprintf(pmmClientConfigureServerCommandTemplate, viper.GetString(config.PMMServerUserKey), viper.GetString(config.PMMServerPassKey), pe.pmmClient.ServerAddr)

	return pe.sshConn.ExecuteCommandWithoutOutput(sql)
}

// StartClient starts pmm client
func (pe *PMMExecutor) StartClient() error {
	return pe.sshConn.ExecuteCommandWithoutOutput(pmmClientStartClientCommandTemplate)
}

// CheckServiceExists checks if the service exists
func (pe *PMMExecutor) CheckServiceExists() (bool, error) {
	output, err := pe.sshConn.ExecuteCommand(fmt.Sprintf(pmmClientCheckServiceCommandTemplate, pe.portNum))
	if err != nil {
		return false, err
	}

	return output == strconv.Itoa(constant.OneInt), nil
}

// AddService adds service to pmm server
func (pe *PMMExecutor) AddService() error {
	// get service name
	serviceName, err := pe.getServiceName()
	if err != nil {
		return err
	}
	// add service
	command := fmt.Sprintf(pmmClientAddServiceCommandTemplateV1,
		pe.portNum,
		viper.GetString(config.MySQLUserMonitorUserKey),
		viper.GetString(config.MySQLUserMonitorPassKey),
		serviceName,
	)
	if pe.pmmClient.ReplicationSetName != constant.EmptyString {
		command = fmt.Sprintf(pmmClientAddServiceCommandTemplateV2,
			pe.portNum,
			viper.GetString(config.MySQLUserMonitorUserKey),
			viper.GetString(config.MySQLUserMonitorPassKey),
			pe.pmmClient.ReplicationSetName,
			serviceName,
		)
	}
	err = pe.sshConn.ExecuteCommandWithoutOutput(command)
	if err != nil {
		return err
	}
	// check if the service exists
	exists, err := pe.CheckServiceExists()
	if err != nil {
		return err
	}
	if !exists {
		return errors.Errorf("pmm client add service failed. host_ip: %s, port_num: %d, service_name: %s", pe.hostIP, pe.portNum, serviceName)
	}

	return nil
}

// getInstallationPackageName returns the installation package name
func (pe *PMMExecutor) getInstallationPackageName() string {
	return fmt.Sprintf(pmmClientInstallationPackageNameTemplate, pe.pmmClient.ClientVersion)
}

// getServiceName gets the service name
func (pe *PMMExecutor) getServiceName() (string, error) {
	hostName, err := pe.sshConn.GetHostName()
	if err != nil {
		return constant.EmptyString, err
	}

	return fmt.Sprintf(pmmClientServiceNameTemplate, hostName, pe.portNum), nil
}
