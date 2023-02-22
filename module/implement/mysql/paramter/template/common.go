package template

const (
	Common = `
[client]
socket={{.Socket}}
user={{.User}}
password={{.Pass}}

[mysql]
prompt=[\\u@\\h:\\p][\\d]>
default-character-set=utf8mb4

[mysqld_multi]
log=/data/mysql/data/mysqld_multi/mysqld_multi.log
user=mysqld_multi
pass=mysqld_multi

`
)
