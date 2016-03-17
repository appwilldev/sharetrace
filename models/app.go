package models

import (
	"fmt"
)

func GenerateAppInfoId() (int64, error) {
	return generateSequenceValue("app_id")
}

type AppInfo struct {
	Id        int64  `xorm:"id INT PK NOT NULL" json:"id"`
	Appid     string `xorm:"appid VARCHAR(256) NOT NULL" json:"appid"`
	AppName   string `xorm:"appname VARCHAR(256) NOT NULL" json:"appname"`
	AppSchema string `xorm:"appschema VARCHAR(256) NOT NULL" json:"appschema"`
	AppHost   string `xorm:"apphost VARCHAR(256) DEFAULT NULL" json:"apphost"`
	AppIcon   string `xorm:"appicon VARCHAR(2048) DEFAULT NULL" json:"appicon"`
	Userid    int64  `xorm:"userid VARCHAR(256) DEFAULT NULL" json:"userid"`

	Yue               int `xorm:"yue INT DEFAULT 0" json:"yue"`
	ShareClickMoney   int `xorm:"share_click_money INT DEFAULT 0" json:"share_click_money"`
	ShareInstallMoney int `xorm:"share_install_money INT DEFAULT 0" json:"share_install_money"`
	InstallMoney      int `xorm:"install_money INT DEFAULT 0" json:"install_money"`

	Des        string `xorm:"des TEXT  DEFAULT NULL" json:"des"`
	Status     int    `xorm:"status INT NOT NULL" json:"status"`
	CreatedUTC int    `xorm:"created_utc INT NOT NULL" json:"created_utc"`
}

type AppInfoRet struct {
	Id        int64  `xorm:"id INT PK NOT NULL" json:"id"`
	Appid     string `xorm:"appid VARCHAR(256) NOT NULL" json:"appid"`
	AppName   string `xorm:"appname VARCHAR(256) NOT NULL" json:"appname"`
	AppSchema string `xorm:"appschema VARCHAR(256) NOT NULL" json:"appschema"`
	AppHost   string `xorm:"apphost VARCHAR(256) DEFAULT NULL" json:"apphost"`
	AppIcon   string `xorm:"appicon VARCHAR(2048) DEFAULT NULL" json:"appicon"`
	Userid    int64  `xorm:"userid VARCHAR(256) DEFAULT NULL" json:"userid"`

	Yue               string `json:"yue"`
	ShareClickMoney   string `json:"share_click_money"`
	ShareInstallMoney string `json:"share_install_money"`
	InstallMoney      string `json:"install_money"`

	Des        string `xorm:"des TEXT  DEFAULT NULL" json:"des"`
	Status     int    `xorm:"status INT NOT NULL" json:"status"`
	CreatedUTC int    `xorm:"created_utc INT NOT NULL" json:"created_utc"`
}

func GetAppInfoRet(data *AppInfo) *AppInfoRet {
	ret := new(AppInfoRet)
	ret.Id = data.Id
	ret.Appid = data.Appid
	ret.AppName = data.AppName
	ret.AppSchema = data.AppSchema
	ret.AppHost = data.AppHost
	ret.AppIcon = data.AppIcon
	ret.Userid = data.Userid
	ret.Yue = fmt.Sprintf("%.2f", float64(data.Yue)/100.0)
	ret.ShareClickMoney = fmt.Sprintf("%.2f", float64(data.ShareClickMoney)/100.0)
	ret.ShareInstallMoney = fmt.Sprintf("%.2f", float64(data.ShareInstallMoney)/100.0)
	ret.InstallMoney = fmt.Sprintf("%.2f", float64(data.InstallMoney)/100.0)
	ret.Des = data.Des
	ret.Status = data.Status
	ret.CreatedUTC = data.CreatedUTC
	return ret
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
