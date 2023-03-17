package tmpl

const (
	InitUserScript = `
		alter user root@'localhost' identified by '{{.RootPass}}' ;
		create user root@'%' identified by '{{.RootPass}}' ;
		grant all on *.* to root@'%' with grant option ;
		create user {{.AdminUser}}@'%' identified by '{{.AdminPass}}' ;
		grant all on *.* to {{.AdminUser}}@'%' with grant option ;
	    create user {{.MySQLDMultiUser}}@'localhost' identified by '{{.MySQLDMultiPass}}' ;
		grant shutdown on *.* to {{.MySQLDMultiUser}}@'localhost' ;
	    create user {{.ReplicationUser}}@'%' identified by '{{.ReplicationPass}}' ;
		grant replication client, replication slave on *.* to {{.ReplicationUser}}@'%' ;
		create user {{.MonitorUser}}@'localhost' identified by '{{.MonitorPass}}' ;
		grant select, reload, process, super, replication client on *.* to {{.MonitorUser}}@'localhost' ;
		create user {{.MonitorUser}}@'127.0.0.1' identified by '{{.MonitorPass}}' ;
		grant select, reload, process, super, replication client on *.* to {{.MonitorUser}}@'127.0.0.1' ;
		create user {{.DASUser}}@'%' identified by '{{.DASPass}}' ;
		grant select, reload, process, super, replication client, replication slave on *.* to {{.DASUser}}@'%' ;
	`
)
