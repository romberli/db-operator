package mysql

import (
	"github.com/romberli/db-operator/global"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/log"
)

type DBORepo struct {
	Database middleware.Pool
}

// NewDBORepo returns a new *DBORepo
func NewDBORepo(db middleware.Pool) *DBORepo {
	return newDBORepo(db)
}

// NewDBORepoWithDefault returns a new *DBORepo with default middleware.Pool
func NewDBORepoWithDefault() *DBORepo {
	return newDBORepo(global.DBOMySQLPool)
}

// newDBORepo returns a new *DBORepo
func newDBORepo(db middleware.Pool) *DBORepo {
	return &DBORepo{
		Database: db,
	}
}

// Execute executes given command and placeholders on the middleware
func (dr *DBORepo) Execute(command string, args ...interface{}) (middleware.Result, error) {
	conn, err := dr.Database.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("mysql.Install.DBORepo.Execute(): close database connection failed.\n%+v", err)
		}
	}()

	return conn.Execute(command, args...)
}

// Transaction returns a middleware.Transaction that could execute multiple commands as a transaction
func (dr *DBORepo) Transaction() (middleware.Transaction, error) {
	return dr.Database.Transaction()
}
