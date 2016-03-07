package models

import (
//"fmt"
//"github.com/appwilldev/sharetrace/conf"
//"github.com/go-xorm/xorm"
)

//////////////////////// Click_Trace
func GenerateClickTraceId() (int64, error) {
	return generateSequenceValue("click_trace_id")
}

type ClickTrace struct {
	Id       int64  `xorm:"id BIGINT PK NOT NULL" json:"id"`
	ClickURL string `xorm:"click_url VARCHAR(1024) DEFAULT NULL" json:"click_url"`
	URLHost  string `xorm:"url_host VARCHAR(256) DEFAULT NULL" json:"url_host"`
	Agent    string `xorm:"agent VARCHAR(1024) DEFAULT NULL" json:"agent"`
	AgentIP  string `xorm:"agentip VARCHAR(256) DEFAULT NULL" json:"agentip"`
	AgentId  string `xorm:"agentid VARCHAR(256) DEFAULT NULL" json:"agentid"`
	Des      string `xorm:"des TEXT DEFAULT NULL" json:"des"`

	Status     int `xorm:"status INT DEFAULT 0" json:"status"`
	CreatedUTC int `xorm:"created_utc INT" json:"created_utc"`
}

func (ClickTrace) TableName() string {
	return "click_trace"
}

func (m *ClickTrace) UniqueCond() (string, []interface{}) {
	var res = make([]interface{}, 1)
	res[0] = m.Id
	return "id=?", res
}

func GetClickTraceById(s *ModelSession, Id int64) (*ClickTrace, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}
	data := &ClickTrace{}
	has, err := s.Id(Id).Get(data)
	if !has || err != nil {
		return nil, err
	}
	return data, nil
}

func GetClickTraceByAgentId(s *ModelSession, IdStr string) (*ClickTrace, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}
	data := &ClickTrace{}
	has, err := s.Where("agentid=?", IdStr).Get(data)
	if !has || err != nil {
		return nil, err
	}
	return data, nil
}
