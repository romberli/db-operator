package mysql

import (
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
)

const (
	installOperationType = iota + 1
	upgradeOperationType
	removeInstanceOperationType
	removeBinaryOperationType
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
	// init operation id
	addrs := common.ConvertStringSliceToString(s.Engine.Addrs, constant.CommaString)
	operationID, err := s.DBORepo.InitOperationHistory(installOperationType, addrs)
	if err != nil {
		return err
	}
	// get lock
	err = s.DBORepo.GetLock(operationID, s.Engine.Addrs)
	if err != nil {
		return err
	}
	defer func() {
		err = s.DBORepo.ReleaseLock(operationID)
		if err != nil {
			log.Errorf(constant.LogWithStackString, err)
		}
	}()
	// install mysql
	return s.Engine.Install(operationID)
}
