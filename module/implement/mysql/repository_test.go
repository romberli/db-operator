package mysql

var (
	testDBORepo *DBORepo
)

func init() {
	testInitViper()
	testInitDBOMySQLPool()
	testDBORepo = testInitDBRepo()
}

func testInitDBRepo() *DBORepo {
	return NewDBORepoWithDefault()
}
