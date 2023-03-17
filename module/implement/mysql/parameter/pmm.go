package parameter

import (
	"github.com/romberli/db-operator/config"
	"github.com/romberli/go-util/constant"
	"github.com/spf13/viper"
)

type PMMClient struct {
	ServerAddr         string `json:"server_addr"`
	ServiceName        string `json:"service_name"`
	ClientVersion      string `json:"client_version"`
	ReplicationSetName string `json:"replication_set_name"`
}

// NewPMMClient returns a new *PMMClient
func NewPMMClient(serverAddr, serviceName, clientVersion, replicationSetName string) *PMMClient {
	return newPMMClient(serverAddr, serviceName, clientVersion, replicationSetName)
}

// NewPMMClientWithDefault returns a new *PMMClient with default values
func NewPMMClientWithDefault() *PMMClient {
	return newPMMClient(
		viper.GetString(config.PMMServerAddrKey),
		constant.EmptyString,
		viper.GetString(config.PMMClientVersionKey),
		constant.EmptyString,
	)
}

// newPMMClient returns a new *PMMClient
func newPMMClient(serverAddr, serviceName, clientVersion, replicationSetName string) *PMMClient {
	return &PMMClient{
		ServerAddr:         serverAddr,
		ServiceName:        serviceName,
		ClientVersion:      clientVersion,
		ReplicationSetName: replicationSetName,
	}
}

// SetServiceName sets service name
func (p *PMMClient) SetServiceName(serviceName string) {
	p.ServiceName = serviceName
}

// SetReplicationSetName sets replication set name
func (p *PMMClient) SetReplicationSetName(replicationSetName string) {
	p.ReplicationSetName = replicationSetName
}
