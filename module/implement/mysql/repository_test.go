package mysql

import (
	"fmt"
	"testing"

	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
	"github.com/stretchr/testify/assert"
)

const (
	testOperationID       = 1
	testOperationDetailID = 1
)

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

func testTruncateOperationInfo() error {
	sql := `truncate table t_mysql_operation_info ;`
	_, err := testDBORepo.Execute(sql)
	if err != nil {
		return err
	}

	sql = `truncate table t_mysql_operation_detail ;`
	_, err = testDBORepo.Execute(sql)
	if err != nil {
		return err
	}

	sql = `truncate table t_mysql_operation_lock ;`
	_, err = testDBORepo.Execute(sql)
	if err != nil {
		return err
	}

	return nil
}

func TestDBRepo_All(t *testing.T) {
	TestDBRepo_Execute(t)
	TestDBRepo_GetOperationHistory(t)
	TestDBRepo_GetOperationDetail(t)
	TestDBRepo_GetLock(t)
	TestDBRepo_ReleaseLock(t)
	TestDBRepo_InitOperationHistory(t)
	TestDBRepo_UpdateOperationHistory(t)
	TestDBRepo_InitOperationDetail(t)
	TestDBRepo_UpdateOperationDetail(t)
}

func TestDBRepo_Execute(t *testing.T) {
	asst := assert.New(t)

	result, err := testDBORepo.Execute("select 1")
	asst.Nil(err, "test Execute() failed")
	data, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	asst.Nil(err, "test Execute() failed")
	asst.Equal(constant.OneInt, data, "test Execute() failed")
}

func TestDBRepo_GetOperationHistory(t *testing.T) {
	asst := assert.New(t)
	// init operation history
	operationID, err := testDBORepo.InitOperationHistory(defaultInstallOperation, []string{testAddr1, testAddr2})
	asst.Nil(err, "test GetOperationHistory() failed")
	// get operation history
	operationInfo, err := testDBORepo.GetOperationHistory(operationID)
	asst.Nil(err, "test GetOperationHistory() failed")
	asst.Equal(testOperationID, operationInfo.ID, "test GetOperationHistory() failed")
	// truncate operation info
	err = testTruncateOperationInfo()
	asst.Nil(err, "test GetOperationHistory() failed")
}

func TestDBRepo_GetOperationDetail(t *testing.T) {
	asst := assert.New(t)
	// init operation history
	operationID, err := testDBORepo.InitOperationHistory(defaultInstallOperation, []string{testAddr1, testAddr2})
	asst.Nil(err, "test GetOperationDetail() failed")
	operationDetailID, err := testDBORepo.InitOperationDetail(operationID, testHostIP1, testPortNum1)
	asst.Nil(err, "test GetOperationDetail() failed")
	// get operation detail
	operationDetails, err := testDBORepo.GetOperationDetails(operationID)
	asst.Nil(err, "test GetOperationDetails() failed")
	asst.Equal(operationDetailID, operationDetails[constant.ZeroInt].ID, "test GetOperationDetails() failed")
	// truncate operation info
	err = testTruncateOperationInfo()
	asst.Nil(err, "test GetOperationDetails() failed")
}

func TestDBRepo_GetLock(t *testing.T) {
	asst := assert.New(t)

	err := testDBORepo.GetLock(testOperationID, []string{testAddr1, testAddr2})
	asst.Nil(err, "test GetLock() failed")
	sql := "select count(*) from t_mysql_operation_lock where operation_id = ? and addr in (%s) ;"
	inClause, err := middleware.ConvertSliceToString(testAddr1, testAddr2)
	asst.Nil(err, "test GetLock() failed")
	sql = fmt.Sprintf(sql, inClause)
	result, err := testDBORepo.Execute(sql, testOperationID)
	asst.Nil(err, "test GetLock() failed")
	count, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	asst.Nil(err, "test GetLock() failed")
	asst.Equal(constant.TwoInt, count, "test GetLock() failed")
	// truncate operation info
	err = testTruncateOperationInfo()
	asst.Nil(err, "test GetLock() failed")
}

func TestDBRepo_ReleaseLock(t *testing.T) {
	asst := assert.New(t)
	// get lock
	err := testDBORepo.GetLock(testOperationID, []string{testAddr1, testAddr2})
	asst.Nil(err, "test ReleaseLock() failed")
	sql := "select count(*) from t_mysql_operation_lock where operation_id = ? and addr in (%s) ;"
	inClause, err := middleware.ConvertSliceToString(testAddr1, testAddr2)
	asst.Nil(err, "test ReleaseLock() failed")
	sql = fmt.Sprintf(sql, inClause)
	result, err := testDBORepo.Execute(sql, testOperationID)
	asst.Nil(err, "test ReleaseLock() failed")
	count, err := result.GetInt(constant.ZeroInt, constant.ZeroInt)
	asst.Nil(err, "test ReleaseLock() failed")
	asst.Equal(constant.TwoInt, count, "test ReleaseLock() failed")
	// release lock
	err = testDBORepo.ReleaseLock(testOperationID)
	asst.Nil(err, "test ReleaseLock() failed")
	// truncate operation info
	err = testTruncateOperationInfo()
	asst.Nil(err, "test ReleaseLock() failed")
}

func TestDBRepo_InitOperationHistory(t *testing.T) {
	asst := assert.New(t)

	operationID, err := testDBORepo.InitOperationHistory(defaultInstallOperation, []string{testAddr1, testAddr2})
	asst.Nil(err, "test InitOperationHistory() failed")
	asst.Equal(testOperationID, operationID, "test InitOperationHistory() failed")
	// truncate operation info
	err = testTruncateOperationInfo()
	asst.Nil(err, "test InitOperationHistory() failed")
}

func TestDBRepo_UpdateOperationHistory(t *testing.T) {
	asst := assert.New(t)
	// init operation history
	operationID, err := testDBORepo.InitOperationHistory(defaultInstallOperation, []string{testAddr1, testAddr2})
	asst.Nil(err, "test UpdateOperationHistory() failed")
	// update operation history
	err = testDBORepo.UpdateOperationHistory(operationID, defaultSuccessStatus, constant.EmptyString)
	asst.Nil(err, "test UpdateOperationHistory() failed")
	// get operation history
	operationInfo, err := testDBORepo.GetOperationHistory(operationID)
	asst.Nil(err, "test UpdateOperationHistory() failed")
	asst.Equal(defaultSuccessStatus, operationInfo.Status, "test UpdateOperationHistory() failed")
	// truncate operation info
	err = testTruncateOperationInfo()
	asst.Nil(err, "test UpdateOperationHistory() failed")
}

func TestDBRepo_InitOperationDetail(t *testing.T) {
	asst := assert.New(t)

	operationDetailID, err := testDBORepo.InitOperationDetail(testOperationID, testHostIP1, testPortNum1)
	asst.Nil(err, "test InitOperationHistory() failed")
	asst.Equal(testOperationDetailID, operationDetailID, "test InitOperationDetail() failed")
	// truncate operation info
	err = testTruncateOperationInfo()
	asst.Nil(err, "test InitOperationDetail() failed")
}

func TestDBRepo_UpdateOperationDetail(t *testing.T) {
	asst := assert.New(t)

	operationDetailID, err := testDBORepo.InitOperationDetail(testOperationID, testHostIP1, testPortNum1)
	asst.Nil(err, "test UpdateOperationDetail() failed")
	err = testDBORepo.UpdateOperationDetail(operationDetailID, defaultSuccessStatus, constant.EmptyString)
	asst.Nil(err, "test UpdateOperationDetail() failed")
	operationDetails, err := testDBORepo.GetOperationDetails(testOperationID)
	asst.Nil(err, "test UpdateOperationDetail() failed")
	asst.Equal(defaultSuccessStatus, operationDetails[constant.ZeroInt].Status, "test UpdateOperationDetail() failed")
	// truncate operation info
	err = testTruncateOperationInfo()
	asst.Nil(err, "test UpdateOperationDetail() failed")
}
