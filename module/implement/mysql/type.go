package mysql

import (
	"github.com/romberli/go-util/constant"
	"time"
)

type OperationInfo struct {
	ID             int       `json:"id" middleware:"id"`
	OperationType  int       `json:"operation_type" middleware:"operation_type"`
	Addrs          string    `json:"addrs" middleware:"addrs"`
	Status         int       `json:"status" middleware:"status"`
	Message        string    `json:"message" middleware:"message"`
	DelFlag        int       `json:"del_flag" middleware:"del_flag"`
	CreateTime     time.Time `json:"create_time" middleware:"create_time"`
	LastUpdateTime time.Time `json:"last_update_time" middleware:"last_update_time"`
}

// NewOperationInfoWithDefault returns a new *OperationInfo with default value
func NewOperationInfoWithDefault() *OperationInfo {
	return &OperationInfo{
		ID:             constant.ZeroInt,
		OperationType:  constant.ZeroInt,
		Addrs:          constant.EmptyString,
		Status:         constant.ZeroInt,
		Message:        constant.EmptyString,
		DelFlag:        constant.ZeroInt,
		CreateTime:     time.Time{},
		LastUpdateTime: time.Time{},
	}
}

type OperationDetail struct {
	ID             int       `json:"id" middleware:"id"`
	OperationID    int       `json:"operation_id" middleware:"operation_id"`
	HostIP         string    `json:"host_ip" middleware:"host_ip"`
	PortNum        int       `json:"port_num" middleware:"port_num"`
	Status         int       `json:"status" middleware:"status"`
	Message        string    `json:"message" middleware:"message"`
	DelFlag        int       `json:"del_flag" middleware:"del_flag"`
	CreateTime     time.Time `json:"create_time" middleware:"create_time"`
	LastUpdateTime time.Time `json:"last_update_time" middleware:"last_update_time"`
}

// NewOperationDetailWithDefault returns a new *OperationDetail with default value
func NewOperationDetailWithDefault() *OperationDetail {
	return &OperationDetail{
		ID:             constant.ZeroInt,
		OperationID:    constant.ZeroInt,
		HostIP:         constant.EmptyString,
		PortNum:        constant.ZeroInt,
		Status:         constant.ZeroInt,
		Message:        constant.EmptyString,
		DelFlag:        constant.ZeroInt,
		CreateTime:     time.Time{},
		LastUpdateTime: time.Time{},
	}
}
