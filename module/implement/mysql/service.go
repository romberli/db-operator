package mysql

import (
	"github.com/romberli/log"
)

type Service struct {
	*DBORepo
	Engine *Engine
}

// NewService returns a new *Service
func NewService(repo *DBORepo, engine *Engine) *Service {
	return newService(repo, engine)
}

// NewServiceWithDefault returns a new *Service with default value
func NewServiceWithDefault(engine *Engine) *Service {
	return newService(NewDBORepoWithDefault(), engine)
}

// newService returns a new *Service
func newService(repo *DBORepo, engine *Engine) *Service {
	return &Service{
		DBORepo: repo,
		Engine:  engine,
	}
}

// Install installs the mysql to the target hosts
func (s *Service) Install() error {
	err := s.DBORepo.GetLock(s.Engine.Addrs)
	if err != nil {
		return err
	}
	defer func() {
		err = s.DBORepo.ReleaseLock(s.Engine.Addrs)
		if err != nil {
			log.Errorf("mysql Service.Install(): release lock failed.\n%+v", err)
		}
	}()

	return s.Engine.Install()
}
