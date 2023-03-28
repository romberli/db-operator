package global

import (
	"github.com/romberli/db-operator/config"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/log"
	"github.com/spf13/viper"
	"time"
)

const (
	purgeInterval = 60 * time.Second
)

type PurgeRepo struct {
	Database middleware.Pool
}

// NewPurgeRepo returns a new *PurgeRepo
func NewPurgeRepo(db middleware.Pool) *PurgeRepo {
	return newPurgeRepo(db)
}

// NewPurgeRepoWithGlobal returns a new *PurgeRepo with global middleware.Pool
func NewPurgeRepoWithGlobal() *PurgeRepo {
	return newPurgeRepo(DBOMySQLPool)
}

// newPurgeRepo returns a new *PurgeRepo
func newPurgeRepo(db middleware.Pool) *PurgeRepo {
	return &PurgeRepo{Database: db}
}

// Execute executes given command and placeholders on the middleware
func (pr *PurgeRepo) Execute(command string, args ...interface{}) (middleware.Result, error) {
	conn, err := pr.Database.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("global PurgeRepo.Execute(): close database connection failed.\n%+v", err)
		}
	}()

	return conn.Execute(command, args...)
}

// PurgeMySQLOperationLock purges the mysql operation lock
func (pr *PurgeRepo) PurgeMySQLOperationLock() error {
	timeout := time.Duration(viper.GetInt(config.MySQLOperationTimeoutKey)) * time.Second
	minTime := time.Now().Add(-timeout).Format(constant.TimeLayoutSecond)

	sql := `DELETE FROM t_mysql_operation_lock WHERE last_update_time < ? ;`
	log.Debugf("global PurgeRepo.PurgeMySQLOperationLock(): sql: %s, args: %s", sql, minTime)

	_, err := pr.Execute(sql, minTime)

	return err
}

type PurgeService struct {
	*PurgeRepo
}

// NewPurgeService returns a new *PurgeService
func NewPurgeService(repo *PurgeRepo) *PurgeService {
	return newPurgeService(repo)
}

// NewPurgeServiceWithDefault returns a new *PurgeService with global middleware.Pool
func NewPurgeServiceWithDefault() *PurgeService {
	return newPurgeService(NewPurgeRepoWithGlobal())
}

// newPurgeService returns a new *PurgeService
func newPurgeService(repo *PurgeRepo) *PurgeService {
	return &PurgeService{PurgeRepo: repo}
}

// Purge purges the mysql operation lock, it will execute periodically
func (ps *PurgeService) PurgeMySQLOperationLock() {
	for {
		err := ps.PurgeRepo.PurgeMySQLOperationLock()
		if err != nil {
			log.Errorf("global PurgeService.PurgeMySQLOperationLock(): purge mysql operation lock failed.\n%+v", err)
		}

		time.Sleep(purgeInterval)
	}
}
