package parameter

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/pingcap/errors"
	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/module/implement/mysql/mode"
	"github.com/romberli/db-operator/module/implement/mysql/parameter/tmpl"
	"github.com/romberli/db-operator/pkg/util/mysql"
	"github.com/romberli/go-util/constant"
	"github.com/spf13/viper"
)

const (
	DefaultServerID = 3306000001

	DefaultDataDirBaseName = "/data/mysql/data"
	DefaultLogDirBaseName  = "/data/mysql/data"

	DefaultSemiSyncSourceEnabled           = 0
	DefaultSemiSyncReplicaEnabled          = 1
	DefaultSemiSyncSourceTimeout           = 10000
	DefaultGroupReplicationConsistency     = "eventual"
	DefaultGroupReplicationFlowControlMode = "disabled"
	DefaultGroupReplicationMemberWeight    = 50

	DefaultBinaryDirBaseTemplate   = "/data/mysql/mysql%s"
	DefaultBinlogExpireLogsSeconds = 604800
	DefaultBinlogExpireLogsDays    = 7
	DefaultBackupDir               = "/data/backup"
	DefaultMaxConnections          = 2000
	DefaultInnodbIOCapacity        = 1000
	DefaultInnodbIOCapacityMax     = 2000
	DefaultServerIDTemplate        = "%d%03s%03s"

	initUserScriptTemplateName = "InitUserScript"
)

type MySQLServer struct {
	Version                         string `json:"version" config:"version"`
	HostIP                          string `json:"host_ip" config:"host_ip"`
	PortNum                         int    `json:"port_num" config:"port_num"`
	RootPass                        string `json:"root_pass" config:"root_pass"`
	AdminUser                       string `json:"admin_user" config:"admin_user"`
	AdminPass                       string `json:"admin_pass" config:"admin_pass"`
	ClientUser                      string `json:"client_user" config:"client_user"`
	ClientPass                      string `json:"client_pass" config:"client_pass"`
	MySQLDMultiUser                 string `json:"mysqld_multi_user" config:"mysqld_multi_user"`
	MySQLDMultiPass                 string `json:"mysqld_multi_pass" config:"mysqld_multi_pass"`
	ReplicationUser                 string `json:"replication_user" config:"replication_user"`
	ReplicationPass                 string `json:"replication_pass" config:"replication_pass"`
	MonitorUser                     string `json:"monitor_user" config:"monitor_user"`
	MonitorPass                     string `json:"monitor_pass" config:"monitor_pass"`
	DASUser                         string `json:"das_user" config:"das_user"`
	DASPass                         string `json:"das_pass" config:"das_pass"`
	Title                           string `json:"title" config:"title"`
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

// NewMySQLServer returns a new *MySQLServer
func NewMySQLServer(version, hostIP string, portNum int, rootPass, adminUser, adminPass, clientUser, clientPass,
	mysqldMultiUser, mysqldMultiPass, replicationUser, replicationPass, monitorUser, monitorPass, dasUser, dasPaas,
	title, binaryDirBase, dataDirBaseName, logDirBaseName string,
	semiSyncSourceEnabled, semiSyncReplicaEnabled, semiSyncSourceTimeout int,
	groupReplicationConsistency, groupReplicationFlowControlMode string, groupReplicationMemberWeight, serverID,
	binlogExpireLogsSeconds, binlogExpireLogsDays int, backupDir string,
	maxConnections int, innodbBufferPoolSize string, innodbIOCapacity int) *MySQLServer {
	return newMySQLServer(
		version,
		hostIP,
		portNum,
		rootPass,
		adminUser,
		adminPass,
		clientUser,
		clientPass,
		mysqldMultiUser,
		mysqldMultiPass,
		replicationUser,
		replicationPass,
		monitorUser,
		monitorPass,
		dasUser,
		dasPaas,
		title,
		binaryDirBase,
		dataDirBaseName,
		logDirBaseName,
		semiSyncSourceEnabled,
		semiSyncReplicaEnabled,
		semiSyncSourceTimeout,
		groupReplicationConsistency,
		groupReplicationFlowControlMode,
		groupReplicationMemberWeight,
		serverID,
		binlogExpireLogsSeconds,
		binlogExpireLogsDays,
		backupDir,
		maxConnections,
		innodbBufferPoolSize,
		innodbIOCapacity,
		innodbIOCapacity*2,
	)
}

// NewMySQLServerWithDefault returns a new *MySQLServer with default values
func NewMySQLServerWithDefault() *MySQLServer {
	return newMySQLServer(
		viper.GetString(config.MySQLVersionKey),
		constant.DefaultLocalHostIP,
		constant.DefaultMySQLPort,
		viper.GetString(config.MySQLUserRootPassKey),
		viper.GetString(config.MySQLUserAdminUserKey),
		viper.GetString(config.MySQLUserAdminPassKey),
		constant.DefaultRootUserName,
		viper.GetString(config.MySQLUserRootPassKey),
		viper.GetString(config.MySQLUserMySQLDMultiUserKey),
		viper.GetString(config.MySQLUserMySQLDMultiPassKey),
		viper.GetString(config.MySQLUserReplicationUserKey),
		viper.GetString(config.MySQLUserReplicationPassKey),
		viper.GetString(config.MySQLUserMonitorUserKey),
		viper.GetString(config.MySQLUserMonitorPassKey),
		viper.GetString(config.MySQLUserDASUserKey),
		viper.GetString(config.MySQLUserDASPassKey),
		defaultTitle,
		fmt.Sprintf(DefaultBinaryDirBaseTemplate, viper.GetString(config.MySQLVersionKey)),
		DefaultDataDirBaseName,
		DefaultLogDirBaseName,
		DefaultSemiSyncSourceEnabled,
		DefaultSemiSyncReplicaEnabled,
		DefaultSemiSyncSourceTimeout,
		DefaultGroupReplicationConsistency,
		DefaultGroupReplicationFlowControlMode,
		DefaultGroupReplicationMemberWeight,
		DefaultServerID,
		DefaultBinlogExpireLogsSeconds,
		DefaultBinlogExpireLogsDays,
		DefaultBackupDir,
		DefaultMaxConnections,
		viper.GetString(config.MySQLParameterInnodbBufferPoolSizeKey),
		DefaultInnodbIOCapacity,
		DefaultInnodbIOCapacityMax,
	)
}

// newMySQLServer returns a new *MySQLServer
func newMySQLServer(version, hostIP string, portNum int, rootPass, adminUser, adminPass, clientUser, clientPass,
	mysqldMultiUser, mysqldMultiPass, replicationUser, replicationPass, monitorUser, monitorPass, dasUser, dasPaas,
	title, binaryDirBase, dataDirBaseName, logDirBaseName string,
	semiSyncSourceEnabled, semiSyncReplicaEnabled, semiSyncSourceTimeout int,
	groupReplicationConsistency, groupReplicationFlowControlMode string,
	groupReplicationMemberWeight, serverID, binlogExpireLogsSeconds, binlogExpireLogsDays int, backupDir string,
	maxConnections int, innodbBufferPoolSize string, innodbIOCapacity, innodbIOCapacityMax int) *MySQLServer {
	return &MySQLServer{
		Version:                         version,
		HostIP:                          hostIP,
		PortNum:                         portNum,
		RootPass:                        rootPass,
		AdminUser:                       adminUser,
		AdminPass:                       adminPass,
		ClientUser:                      clientUser,
		ClientPass:                      clientPass,
		MySQLDMultiUser:                 mysqldMultiUser,
		MySQLDMultiPass:                 mysqldMultiPass,
		ReplicationUser:                 replicationUser,
		ReplicationPass:                 replicationPass,
		MonitorUser:                     monitorUser,
		MonitorPass:                     monitorPass,
		DASUser:                         dasUser,
		DASPass:                         dasPaas,
		Title:                           title,
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
		InnodbIOCapacityMax:             innodbIOCapacityMax,
	}
}

// GetCommon() gets the Common of MySQLServer
func (ms *MySQLServer) GetCommon() *Common {
	return NewCommon(
		ms.PortNum,
		ms.DataDirBaseName,
		ms.ClientUser,
		ms.ClientPass,
		ms.MySQLDMultiUser,
		ms.MySQLDMultiPass,
	)
}

// GetMySQLD() gets the MySQLD of MySQLServer
func (ms *MySQLServer) GetMySQLD() *MySQLD {
	return NewMySQLD(
		ms.Version,
		ms.HostIP,
		ms.PortNum,
		ms.BinaryDirBase,
		ms.DataDirBaseName,
		ms.LogDirBaseName,
		ms.SemiSyncSourceEnabled,
		ms.SemiSyncReplicaEnabled,
		ms.SemiSyncSourceTimeout,
		ms.GroupReplicationConsistency,
		ms.GroupReplicationFlowControlMode,
		ms.GroupReplicationMemberWeight,
		ms.ServerID,
		ms.BinlogExpireLogsSeconds,
		ms.BinlogExpireLogsDays,
		ms.BackupDir,
		ms.MaxConnections,
		ms.InnodbBufferPoolSize,
		ms.InnodbIOCapacity,
	)
}

// SetSemiSyncSourceEnabled sets the semi-sync source enabled
func (ms *MySQLServer) SetSemiSyncSourceEnabled(semiSyncSourceEnabled int) {
	ms.SemiSyncSourceEnabled = semiSyncSourceEnabled
}

// SetSemiSyncReplicaEnabled sets the semi-sync replica enabled
func (ms *MySQLServer) SetSemiSyncReplicaEnabled(semiSyncReplicaEnabled int) {
	ms.SemiSyncReplicaEnabled = semiSyncReplicaEnabled
}

// InitWithHostInfo resets the MySQLServer with given host info
func (ms *MySQLServer) InitWithHostInfo(hostIP string, portNum int) error {
	ms.HostIP = hostIP
	ms.PortNum = portNum

	ms.DataDirBase = fmt.Sprintf(dirBaseTemplate, ms.DataDirBaseName, portNum)
	ms.LogDirBase = fmt.Sprintf(dirBaseTemplate, ms.LogDirBaseName, portNum)

	ipList := strings.Split(hostIP, constant.DotString)
	serverIDStr := fmt.Sprintf(DefaultServerIDTemplate, portNum, ipList[constant.TwoInt], ipList[constant.ThreeInt])

	var err error

	ms.ServerID, err = strconv.Atoi(serverIDStr)

	return errors.Trace(err)
}

func (ms *MySQLServer) GetCommonConfig() ([]byte, error) {
	return ms.GetCommon().GetConfig()
}

func (ms *MySQLServer) GetMySQLDConfig(v *version.Version, m mode.Mode) ([]byte, error) {
	return ms.GetMySQLD().GetConfig(v, m)
}

// GetMySQLDConfigWithTitle gets the MySQLD configuration with given title
func (ms *MySQLServer) GetMySQLDConfigWithTitle(title string, v *version.Version, m mode.Mode) ([]byte, error) {
	return ms.GetMySQLD().GetConfigWithTitle(title, v, m)
}

// GetConfig() gets the configuration of MySQLServer
func (ms *MySQLServer) GetConfig(v *version.Version, m mode.Mode) ([]byte, error) {
	commonConfig, err := ms.GetCommonConfig()
	if err != nil {
		return nil, err
	}

	mysqldConfig, err := ms.GetMySQLDConfig(v, m)
	if err != nil {
		return nil, err
	}

	return append(commonConfig, mysqldConfig...), nil
}

// GetConfigWithTitle gets the configuration of MySQLServer with given title
func (ms *MySQLServer) GetConfigWithTitle(title string, v *version.Version, m mode.Mode) ([]byte, error) {
	commonConfig, err := ms.GetCommonConfig()
	if err != nil {
		return nil, err
	}

	mysqldConfig, err := ms.GetMySQLDConfigWithTitle(title, v, m)
	if err != nil {
		return nil, err
	}

	return append(commonConfig, mysqldConfig...), nil
}

// WriteConfig writes the configuration to the specified file
func (ms *MySQLServer) WriteConfig(configPath string, data []byte) error {
	return os.WriteFile(configPath, data, constant.DefaultFileMode)
}

// GetInitUserSQL gets the SQL to initialize the user
func (ms *MySQLServer) GetInitUserSQL() (string, error) {
	sqlBytes, err := mysql.GetConfig(initUserScriptTemplateName, tmpl.InitUserScript, ms)
	if err != nil {
		return constant.EmptyString, err
	}

	return string(sqlBytes), nil
}

// Marshal marshals the MySQLServer to json bytes
func (ms *MySQLServer) Marshal() ([]byte, error) {
	return json.Marshal(ms)
}

// Unmarshal unmarshals the json bytes to MySQLServer
func (ms *MySQLServer) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, &ms)
	if err != nil {
		return errors.Trace(err)
	}

	serverID, err := mysql.GetServerID(ms.HostIP, ms.PortNum)
	if err != nil {
		return err
	}

	ms.ServerID = serverID

	return nil
}
