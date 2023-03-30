package mysql

import (
	"github.com/romberli/go-util/common"
	"strings"
	"time"

	"github.com/romberli/db-operator/global"
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/log"

	msgMySQL "github.com/romberli/db-operator/pkg/message/mysql"
)

const (
	minRunAgainInterval = 30 * time.Minute

	defaultRunningStatus = 1
	defaultSuccessStatus = 2
	defaultFailedStatus  = 3
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

// GetOperationHistory gets the mysql operation history from the middleware
func (dr *DBORepo) GetOperationHistory(id int) (*OperationInfo, error) {
	sql := `
		SELECT id,
			   operation_type,
			   addrs,
			   status,
			   message,
			   del_flag,
			   create_time,
			   last_update_time
		FROM t_mysql_operation_info
		WHERE del_flag = 0
		  AND id = ?
		ORDER BY id DESC
	`
	log.Debugf("mysql DBORepo.GetOperationHistory() select sql: \n%s\nplaceholders: %d", sql, id)

	result, err := dr.Execute(sql, id)
	if err != nil {
		return nil, err
	}

	operationInfo := NewOperationInfoWithDefault()
	err = result.MapToStructByRowIndex(operationInfo, constant.ZeroInt, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}

	return operationInfo, nil
}

// GetOperationDetail gets the mysql operation detail from the middleware
func (dr *DBORepo) GetOperationDetail(operationID int) ([]*OperationDetail, error) {
	sql := `
		SELECT id,
			   operation_id,
			   host_ip,
			   port_num,
			   status,
			   message,
			   del_flag,
			   create_time,
			   last_update_time
		FROM t_mysql_operation_detail
		WHERE del_flag = 0
		  AND operation_id = ?
		ORDER BY id ASC
	`
	log.Debugf("mysql DBORepo.GetOperationDetail() select sql: \n%s\nplaceholders: %d", sql, operationID)

	result, err := dr.Execute(sql, operationID)
	if err != nil {
		return nil, err
	}

	operationDetailList := make([]*OperationDetail, result.RowNumber())
	for i := constant.ZeroInt; i < result.RowNumber(); i++ {
		operationDetailList[i] = NewOperationDetailWithDefault()
	}

	err = result.MapToStructSlice(operationDetailList, constant.DefaultMiddlewareTag)
	if err != nil {
		return nil, err
	}

	return operationDetailList, nil
}

// GetLock gets the operation lock of the given host info
func (dr *DBORepo) GetLock(operationID int, addrs []string) error {
	// prepare sql
	sql := `INSERT INTO t_mysql_operation_lock(operation_id, addr) VALUES`
	for i := constant.ZeroInt; i < len(addrs); i++ {
		sql = sql + `(?, ?),`
	}
	sql = strings.Trim(sql, constant.CommaString) + constant.SemicolonString
	// prepare placeholders
	placeHolders := make([]interface{}, len(addrs)*constant.TwoInt)
	for i := constant.ZeroInt; i < len(addrs); i++ {
		placeHolders[i*constant.TwoInt] = operationID
		placeHolders[i*constant.TwoInt+constant.OneInt] = addrs[i]
	}
	log.Debugf("mysql DBORepo.GetLock() insert sql: \n%s\noperation_id: %d, addrs: %v",
		sql, operationID, common.ConvertStringSliceToString(addrs, constant.CommaString))

	// execute sql
	_, err := dr.Execute(sql, placeHolders...)
	if err != nil {
		return message.NewMessage(msgMySQL.ErrMySQLRepositoryGetLock, err, addrs)
	}

	return nil
}

// ReleaseLock releases the operation lock of the given host info
func (dr *DBORepo) ReleaseLock(operationID int) error {
	sql := `DELETE FROM t_mysql_operation_lock WHERE operation_id = ? ;`
	log.Debugf("mysql DBORepo.ReleaseLock() delete sql: \n%s\nplaceholders: %d", sql, operationID)

	_, err := dr.Execute(sql)
	if err != nil {
		return message.NewMessage(msgMySQL.ErrMySQLRepositoryReleaseLock, err, operationID)
	}

	return nil
}

// InitOperationHistory initializes the mysql operation history in the middleware
func (dr *DBORepo) InitOperationHistory(operationType int, addrs string) (int, error) {
	sql := `INSERT INTO t_mysql_operation_info(operation_type, addrs, status) VALUES(?, ?, ?) ;`
	log.Debugf("mysql DBORepo.InitOperationHistory() insert sql: \n%s\nplaceholders: %d, %s",
		sql, operationType, addrs, operationStatusRunning)

	result, err := dr.Execute(sql, operationType, addrs, operationStatusRunning)
	if err != nil {
		return constant.ZeroInt, err
	}

	return result.LastInsertID()
}

// UpdateOperationHistory updates the mysql operation history in the middleware
func (dr *DBORepo) UpdateOperationHistory(id int, status int, message string) error {
	sql := `UPDATE t_mysql_operation_info SET status = ?, message = ? WHERE id = ? ;`
	log.Debugf("mysql DBORepo.UpdateOperationHistory() update sql: \n%s\nplaceholders: %d, %d, %s",
		sql, id, status, message)

	_, err := dr.Execute(sql, status, message, id)

	return err
}

// InitOperationDetail initializes the mysql operation detail in the middleware
func (dr *DBORepo) InitOperationDetail(operationID int, hostIP string, portNum int) (int, error) {
	sql := `INSERT INTO t_mysql_operation_detail(operation_id, host_ip, port_num, status) VALUES(?, ?, ?, ?) ;`
	log.Debugf("mysql DBORepo.InitOperationDetail() insert sql: \n%s\nplaceholders: %d, %s, %d",
		sql, operationID, hostIP, portNum, operationStatusRunning)

	result, err := dr.Execute(sql, operationID, hostIP, portNum, operationStatusRunning)
	if err != nil {
		return constant.ZeroInt, err
	}

	return result.LastInsertID()
}

// UpdateOperationDetail updates the mysql operation detail in the middleware
func (dr *DBORepo) UpdateOperationDetail(operationDetailID int, status int, message string) error {
	sql := `UPDATE t_mysql_operation_detail SET status = ?, message = ? WHERE id = ? ;`
	log.Debugf("mysql DBORepo.UpdateOperationDetail() update sql: \n%s\nplaceholders: %d, %d, %s",
		sql, status, message, operationDetailID)

	_, err := dr.Execute(sql, status, message, operationDetailID)

	return err
}
