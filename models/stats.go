package models

import (
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
