package models

import (
	"fmt"
	"reflect"

	"github.com/go-xorm/xorm"
)

func GetShareClickListByAppid(s *ModelSession, appid string) ([]*ShareClick, error) {
	var (
		dataList = make([]*ShareClick, 0)
		err      error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	var session *xorm.Session
	session = s.Join("INNER", "click_session", "share_url.id=click_session.shareid").OrderBy("click_session.created_utc desc")
	err = session.Where("appid=?", appid).Find(&dataList)

	if err != nil {
		return nil, err
	}

	return dataList, nil
}

func GetShareTotalByAppid(s *ModelSession, appid string, date string) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	data := new(ShareURL)
	total, err = s.Where("appid=?", appid).And("date(to_timestamp(created_utc))=?", date).Count(data)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func GetClickTotalByAppid(s *ModelSession, appid string, date string) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	data := new(ShareClick)
	var session *xorm.Session
	session = s.Join("INNER", "click_session", "share_url.id=click_session.shareid")
	total, err = session.Where("appid=?", appid).And("date(to_timestamp(click_session.created_utc))=?", date).Count(data)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func GetInstallTotalByAppid(s *ModelSession, appid string, date string) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	data := new(ShareClick)
	var session *xorm.Session
	session = s.Join("INNER", "click_session", "share_url.id=click_session.shareid")
	total, err = session.Where("appid=?", appid).And("click_session.status = 1").And("date(to_timestamp(click_session.created_utc))=?", date).Count(data)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func GetButtonTotalByAppid(s *ModelSession, appid string, date string) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	data := new(ShareClick)
	var session *xorm.Session
	session = s.Join("INNER", "click_session", "share_url.id=click_session.shareid")
	total, err = session.Where("appid=?", appid).And("(click_session.buttonid is not null and click_session.status = 2) or (click_session.status = 1)").And("date(to_timestamp(click_session.created_utc))=?", date).Count(data)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func GetClickTotalByHost(s *ModelSession, host string, date string) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	data := new(ClickSession)
	total, err = s.Where("url_host=?", host).And("date(to_timestamp(created_utc))=?", date).Count(data)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func GetButtonTotalByHost(s *ModelSession, host string, date string) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	data := new(ClickSession)
	total, err = s.Where("url_host=?", host).And("(buttonid is not null and status = 2) or (status = 1)").And("date(to_timestamp(created_utc))=?", date).Count(data)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func GetTotalByHostPhone(s *ModelSession, host string, date string, phone string) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	data := new(ClickSession)
	total, err = s.Where("url_host=?", host).And("date(to_timestamp(created_utc))=?", date).And("agent like ?", `%`+phone+`%`).Count(data)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func GetTotalByHostiPhone(s *ModelSession, host string, date string, phone string) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	data := new(ClickSession)
	//total, err = s.Where("url_host=?", host).And("date(to_timestamp(created_utc))=?", date).And("agent like ?", `%iPhone OS%`).And("agent like ?", `%`+phone+`%`).Count(data)
	condition := "url_host=? and  date(to_timestamp(created_utc))=? and (agent like ? or agent like ?) and agent like ?"
	total, err = s.Where(condition, host, date, `%iPhone OS%`, `%iPad%`, phone).Count(data)

	if err != nil {
		return -1, err
	}

	return total, nil
}

func GetMoneyTotalByAppid(s *ModelSession, appid string, date string, money_type int) (int64, error) {
	var (
		total int64
		err   error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	sql := fmt.Sprintf("select sum(money) from appuser_money where appid = '%s' and money_type = %d and date(to_timestamp(created_utc)) = '%s'", appid, money_type, date)
	columnTypes := []reflect.Type{reflect.TypeOf(int64(1))}
	res, err := RawSqlQuery(sql, columnTypes)
	if err != nil {
		return -1, err
	}
	total = res[0][0].(int64) / 100

	return total, nil
}
