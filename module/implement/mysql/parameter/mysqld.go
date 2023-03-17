package parameter

import (
	"fmt"
	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/module/implement/mysql/parameter/tmpl"
	"github.com/romberli/db-operator/pkg/util/mysql"
)

const (
	mysql80              = "8.0"
	mysqld80TemplateName = "mysqld80"

	defaultTitle = "mysqld"
)

type MySQLD struct {
	Version                         string `json:"version" config:"version"`
	Title                           string `json:"title" config:"title"`
	HostIP                          string `json:"host_ip" config:"host_ip"`
	PortNum                         int    `json:"port_num" config:"port_num"`
	BinaryDirBase                   string `json:"binary_dir_base" config:"binary_dir_base"`
	DataDirBaseName                 string `json:"data_dir_base_name" config:"data_dir_base_name"`
	DataDirBase                     string `json:"data_dir_base" config:"data_dir_base"`
	LogDirBaseName                  string `json:"log_dir_base_name" config:"log_dir_base_name"`
	LogDirBase                      string `json:"log_dir_base" config:"log_dir_base"`
	SemiSyncSourceEnabled           int    `json:"semi_sync_source_enabled" config:"semi_sync_source_enabled"`
	SemiSyncReplicaEnabled          int    `json:"semi_sync_replica_enabled" config:"semi_sync_replica_enabled"`
	SemiSyncSourceTimeout           int    `json:"semi_sync_source_timeout" config:"semi_sync_source_timeout"`
	GroupReplicationConsistency     string `json:"group_replication_consistency" config:"group_replication_consistency"`
	GroupReplicationFlowControlMode string `json:"group_replication_flow_control_mode" config:"group_replication_flow_control_mode"`
	GroupReplicationMemberWeight    int    `json:"group_replication_member_weight" config:"group_replication_member_weight"`
	ServerID                        int    `json:"server_id" config:"server_id"`
	BinlogExpireLogsSeconds         int    `json:"binlog_expire_logs_seconds" config:"binlog_expire_logs_seconds"`
	BinlogExpireLogsDays            int    `json:"binlog_expire_logs_days" config:"binlog_expire_logs_days"`
	BackupDir                       string `json:"backup_dir" config:"backup_dir"`
	MaxConnections                  int    `json:"max_connections" config:"max_connections"`
	InnodbBufferPoolSize            string `json:"innodb_buffer_pool_size" config:"innodb_buffer_pool_size"`
	InnodbIOCapacity                int    `json:"innodb_io_capacity" config:"innodb_io_capacity"`
	InnodbIOCapacityMax             int    `json:"innodb_io_capacity_max" config:"innodb_io_capacity_max"`
}

// NewMySQLD returns a new *MySQLD
func NewMySQLD(version, hostIP string, portNum int, binaryDirBase, dataDirBaseName, logDirBaseName string,
	semiSyncSourceEnabled, semiSyncReplicaEnabled, semiSyncSourceTimeout int,
	groupReplicationConsistency, groupReplicationFlowControlMode string, groupReplicationMemberWeight int,
	serverID, binlogExpireLogsSeconds, binlogExpireLogsDays int, backupDir string,
	maxConnections int, innodbBufferPoolSize string, innodbIOCapacity int) *MySQLD {
	return &MySQLD{
		Version:                         version,
		Title:                           defaultTitle,
		HostIP:                          hostIP,
		PortNum:                         portNum,
		BinaryDirBase:                   binaryDirBase,
		DataDirBaseName:                 dataDirBaseName,
		DataDirBase:                     fmt.Sprintf(dirBaseTemplate, dataDirBaseName, portNum),
		LogDirBaseName:                  logDirBaseName,
		LogDirBase:                      fmt.Sprintf(dirBaseTemplate, logDirBaseName, portNum),
		SemiSyncSourceEnabled:           semiSyncSourceEnabled,
		SemiSyncReplicaEnabled:          semiSyncReplicaEnabled,
		SemiSyncSourceTimeout:           semiSyncSourceTimeout,
		GroupReplicationConsistency:     groupReplicationConsistency,
		GroupReplicationFlowControlMode: groupReplicationFlowControlMode,
		GroupReplicationMemberWeight:    groupReplicationMemberWeight,
		ServerID:                        serverID,
		BinlogExpireLogsSeconds:         binlogExpireLogsSeconds,
		BinlogExpireLogsDays:            binlogExpireLogsDays,
		BackupDir:                       backupDir,
		MaxConnections:                  maxConnections,
		InnodbBufferPoolSize:            innodbBufferPoolSize,
		InnodbIOCapacity:                innodbIOCapacity,
		InnodbIOCapacityMax:             innodbIOCapacity * 2,
	}
}

// SetSemiSyncSourceEnabled sets the semi-sync source enabled of MySQLD
func (md *MySQLD) SetSemiSyncSourceEnabled(semiSyncSourceEnabled int) {
	md.SemiSyncSourceEnabled = semiSyncSourceEnabled
}

// SetSemiSyncReplicaEnabled sets the semi-sync replica enabled of MySQLD
func (md *MySQLD) SetSemiSyncReplicaEnabled(semiSyncReplicaEnabled int) {
	md.SemiSyncReplicaEnabled = semiSyncReplicaEnabled
}

// GetConfig returns the configuration of MySQLD
func (md *MySQLD) GetConfig(v *version.Version, m mode.Mode) ([]byte, error) {
	mysql80Version, err := version.NewVersion(mysql80)
	if err != nil {
		return nil, err
	}

	if v.GreaterThanOrEqual(mysql80Version) {
		return mysql.GetConfig(mysqld80TemplateName, md.configTemplateContent(tmpl.MySQLD80, m), md)
	}

	return nil, errors.Errorf("version must be larger than 8.0, %s is not supported", md.Version)
}

// GetConfigWithTitle returns the configuration of MySQLD with title
func (md *MySQLD) GetConfigWithTitle(title string, v *version.Version, m mode.Mode) ([]byte, error) {
	md.Title = title

	return md.GetConfig(v, m)
}

// configTemplateContent sets the configuration of MySQLD
func (md *MySQLD) configTemplateContent(template string, m mode.Mode) string {
	switch m {
	case mode.AsyncReplication:
	case mode.SemiSyncReplication:
		template = strings.Replace(template, "#plugin_load", "plugin_load", 1)
		template = strings.ReplaceAll(template, "#rpl_semi_sync", "rpl_semi_sync")
	case mode.GroupReplication:
		template = strings.ReplaceAll(template, "#group_replication", "group_replication")
	}

	return template
}
