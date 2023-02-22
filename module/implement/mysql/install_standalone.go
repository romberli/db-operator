package mysql

type Standalone struct {
	Version        string            `json:"version"`
	PackageType    int               `json:"package_type"`
	InstallDirType int               `json:"install_type"`
	HostIP         string            `json:"host_ip"`
	PortNum        int               `json:"port_num"`
	Params         map[string]string `json:"params"`
}
