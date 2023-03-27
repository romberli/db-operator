package tmpl

const (
	MySQLD80 = `
[{{.Title}}]
port={{.PortNum}}
mysqlx_port={{.PortNum}}0
admin_port={{.PortNum}}2
basedir={{.BinaryDirBase}}
datadir={{.DataDirBase}}/data
tmpdir={{.DataDirBase}}/tmp
socket={{.DataDirBase}}/run/mysql.sock
mysqlx_socket={{.DataDirBase}}/run/mysqlx.sock
pid-file={{.DataDirBase}}/run/mysql.pid
log-error={{.DataDirBase}}/log/mysql.err
#mysqld={{.BinaryDirBase}}/bin/mysqld_safe
#mysqladmin={{.BinaryDirBase}}/bin/mysqladmin
default-time-zone='+08:00'
character-set-server=utf8mb4
thread_cache_size=512
sql_mode=STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION,PIPES_AS_CONCAT,ONLY_FULL_GROUP_BY,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO
#tls_version=''

#plugin_load="rpl_semi_sync_source=semisync_source.so;rpl_semi_sync_replica=semisync_replica.so"
#rpl_semi_sync_source_wait_point=after_sync
#rpl_semi_sync_source_enabled={{.SemiSyncSourceEnabled}}
#rpl_semi_sync_replica_enabled={{.SemiSyncReplicaEnabled}}
#rpl_semi_sync_source_timeout={{.SemiSyncSourceTimeout}}
#rpl_semi_sync_source_wait_for_replica_count=1
#rpl_semi_sync_source_wait_no_replica=1

#group_replication_single_primary_mode=on
#group_replication_consistency={{.GroupReplicationConsistency}}
#group_replication_flow_control_mode={{.GroupReplicationFlowControlMode}}
#group_replication_member_weight={{.GroupReplicationMemberWeight}}

server-id={{.ServerID}}
gtid_mode=on
enforce_gtid_consistency=1
binlog_gtid_simple_recovery=1
sync_binlog=1
log-bin={{.LogDirBase}}/binlog/mysql-bin
binlog_format=row
binlog_row_image=full
max_binlog_size=1G
binlog_cache_size=1M
binlog_error_action=ABORT_SERVER
binlog_expire_logs_seconds={{.BinlogExpireLogsSeconds}}
binlog_cache_size=4m
log_replica_updates=1
relay_log={{.LogDirBase}}/relaylog/mysql-relay
max_relay_log_size=1G
relay_log_purge=1
relay_log_recovery=1
report_host={{.HostIP}}
report_port={{.PortNum}}
replica_parallel_workers=16
replica_preserve_commit_order=1
replica_transaction_retries=128
binlog_transaction_dependency_tracking=writeset
binlog_transaction_dependency_history_size=25000

secure_file_priv={{.BackupDir}}
max_connections={{.MaxConnections}}
transaction-isolation=READ-COMMITTED
table_open_cache=2048
lower_case_table_names=1
max_allowed_packet=64M
tmp_table_size=64M
max_heap_table_size=64M
sort_buffer_size=4M
join_buffer_size=4M
read_buffer_size=8M
read_rnd_buffer_size=4M
key_buffer_size=32M
bulk_insert_buffer_size=64M
innodb_flush_log_at_trx_commit=1
innodb_log_file_size=1G
innodb_log_files_in_group=4
innodb_log_group_home_dir={{.LogDirBase}}/data
innodb_data_file_path=ibdata1:1024M:autoextend
innodb_autoextend_increment=16
innodb_buffer_pool_instances=8
innodb_buffer_pool_size={{.InnodbBufferPoolSize}}
innodb_sort_buffer_size=4M
innodb_log_buffer_size=32M
innodb_read_io_threads=16
innodb_write_io_threads=16
innodb_io_capacity={{.InnodbIOCapacity}}
innodb_io_capacity_max={{.InnodbIOCapacityMax}}
innodb_page_cleaners=16
innodb_flush_method=O_DIRECT
innodb_monitor_enable=ALL
innodb_print_all_deadlocks=1
innodb_numa_interleave=1

general_log=OFF
general_log_file={{.DataDirBase}}/log/general.log
slow_query_log=ON
slow_query_log_file={{.DataDirBase}}/log/mysql-slow.log
long_query_time=0.1
log_output=file
performance_schema=ON

`
)
