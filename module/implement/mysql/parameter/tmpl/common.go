package tmpl

const (
	Common = `[client]
socket={{.DataDirBase}}/mysql.sock
user={{.ClientUser}}
password={{.ClientPass}}

[mysql]
prompt=[\\u@\\h:\\p][\\d]>
default-character-set=utf8mb4

[mysqld_multi]
log={{.DataDirBaseName}}/mysqld_multi/mysqld_multi.log
user={{.MySQLDMultiUser}}
pass={{.MySQLDMultiPass}}
`
)
