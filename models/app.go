package models

import ()

func GenerateAppInfoId() (int64, error) {
	return generateSequenceValue("app_id")
}

type AppInfo struct {
	Id         int64  `xorm:"id INT PK NOT NULL" json:"id"`
	Appid      string `xorm:"appid VARCHAR(256) NOT NULL" json:"appid"`
	AppName    string `xorm:"appname VARCHAR(256) NOT NULL" json:"appname"`
	AppSchema  string `xorm:"appschema VARCHAR(256) NOT NULL" json:"appschema"`
	AppHost    string `xorm:"apphost VARCHAR(256) DEFAULT NULL" json:"apphost"`
	AppIcon    string `xorm:"appicon VARCHAR(2048) DEFAULT NULL" json:"appicon"`
	Userid     int64  `xorm:"userid VARCHAR(256) DEFAULT NULL" json:"userid"`
	Des        string `xorm:"des TEXT  DEFAULT NULL" json:"des"`
	Status     int    `xorm:"status INT NOT NULL" json:"status"`
	CreatedUTC int    `xorm:"created_utc INT NOT NULL" json:"created_utc"`
}

func (*AppInfo) TableName() string {
	return "app_info"
}

func (m *AppInfo) UniqueCond() (string, []interface{}) {
	var res = make([]interface{}, 1)
	res[0] = m.Id

	return "id=?", res
}

func GetAppInfoById(s *ModelSession, id int64) (*AppInfo, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	app := &AppInfo{}
	has, err := s.Id(id).Get(app)
	if !has || err != nil {
		return nil, err
	}

	return app, nil
}

func GetAppInfoByAppid(s *ModelSession, appid string) (*AppInfo, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	app := &AppInfo{}
	has, err := s.Where("appid=?", appid).Get(app)
	if !has || err != nil {
		return nil, err
	}

	return app, nil
}

func GetAppInfoAll(s *ModelSession) ([]*AppInfo, int64, error) {
	var (
		total    int64
		dataList = make([]*AppInfo, 0)
		err      error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	total, _ = s.Count(new(AppInfo))
	err = s.OrderBy("id desc").Find(&dataList)

	if err != nil {
		return nil, 0, err
	}

	return dataList, total, nil
}

func GetAppInfoListByUserid(s *ModelSession, userid int64) ([]*AppInfo, int64, error) {
	var (
		total    int64
		dataList = make([]*AppInfo, 0)
		err      error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	total, _ = s.Count(new(AppInfo))
	err = s.Where("userid=?", userid).OrderBy("id desc").Find(&dataList)

	if err != nil {
		return nil, 0, err
	}

	return dataList, total, nil
}
