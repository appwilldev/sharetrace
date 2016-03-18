package models

import (
	"fmt"
	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/utils"
	"log"
)

func GenerateAppuserMoneyId() (int64, error) {
	return generateSequenceValue("appuser_money_id")
}

type AppuserMoney struct {
	Id        int64  `xorm:"id INT PK NOT NULL" json:"id"`
	Appid     string `xorm:"appid VARCHAR(256) NOT NULL" json:"appid"`
	Appuserid string `xorm:"appuserid VARCHAR(256) NOT NULL" json:"appuserid"`

	ClickSessionID int64 `xorm:"click_session_id BIGINT DEFAULT 0" json:"click_session_id"`
	UserOrderID    int64 `xorm:"user_order_id BIGINT DEFAULT 0" json:"user_order_id"`

	MoneyType   int     `xorm:"money_type INT DEFAULT 0" json:"money_type"`
	Money       float64 `xorm:"money INT DEFAULT 0" json:"money"`
	MoneyStatus int     `xorm:"money_status INT DEFAULT 0" json:"money_status"`

	Des        string `xorm:"des TEXT  DEFAULT NULL" json:"des"`
	Status     int    `xorm:"status INT NOT NULL" json:"status"`
	CreatedUTC int    `xorm:"created_utc INT NOT NULL" json:"created_utc"`
}

func (*AppuserMoney) TableName() string {
	return "appuser_money"
}

func (m *AppuserMoney) UniqueCond() (string, []interface{}) {
	var res = make([]interface{}, 1)
	res[0] = m.Id

	return "id=?", res
}

func GetAppuserMoneyById(s *ModelSession, id int64) (*AppuserMoney, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	app := &AppuserMoney{}
	has, err := s.Id(id).Get(app)
	if !has || err != nil {
		return nil, err
	}

	return app, nil
}

func GetAppuserMoneyListByUserid(s *ModelSession, appid string, appuserid string) ([]*AppuserMoney, int64, error) {
	var (
		total    int64
		dataList = make([]*AppuserMoney, 0)
		err      error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	total, _ = s.Count(new(AppuserMoney))
	err = s.Where("appid=? and appuserid=? ", appid, appuserid).OrderBy("id desc").Find(&dataList)

	if err != nil {
		return nil, 0, err
	}

	return dataList, total, nil
}

func AddAwardToAppUser(s *ModelSession, app *AppInfo, cs_data *ClickSession) error {

	if s == nil {
		s = NewModelSession()
	}
	defer s.Close()

	su_data, _ := GetShareURLById(nil, cs_data.Shareid)
	if app.Status == 1 && app.Yue > 1000 {
		err := s.Begin()
		if app.ShareInstallMoney > 0 {
			id, err := GenerateAppuserMoneyId()
			if err != nil {
				log.Println(err.Error())
				s.Rollback()
				return err
			}
			apm_data := new(AppuserMoney)
			apm_data.Id = id
			apm_data.Appid = app.Appid
			apm_data.Appuserid = su_data.Fromid
			apm_data.ClickSessionID = cs_data.Id
			apm_data.MoneyType = conf.MONEY_TYPE_INSTALL_SHARER
			apm_data.Money = float64(app.ShareInstallMoney)
			apm_data.CreatedUTC = utils.GetNowSecond()
			apm_data.Des = "分享链接吸引用户" + cs_data.Installid + "安装了App"
			err = InsertDBModel(s, apm_data)
			if err != nil {
				log.Println(err.Error())
				s.Rollback()
				return err
			}
			app.Yue = app.Yue - app.ShareInstallMoney
		}

		if app.InstallMoney > 0 {
			id, err := GenerateAppuserMoneyId()
			if err != nil {
				log.Println(err.Error())
				s.Rollback()
				return err
			}
			apm_data := new(AppuserMoney)
			apm_data.Id = id
			apm_data.Appid = app.Appid
			apm_data.Appuserid = cs_data.Installid
			apm_data.ClickSessionID = cs_data.Id
			apm_data.MoneyType = conf.MONEY_TYPE_INSTALL_INSTALLER
			apm_data.CreatedUTC = utils.GetNowSecond()
			apm_data.Money = float64(app.InstallMoney)
			apm_data.Des = "通过点击分享链接安装了App"
			err = InsertDBModel(s, apm_data)
			if err != nil {
				log.Println(err.Error())
				s.Rollback()
				return err
			}
			app.Yue = app.Yue - app.ShareInstallMoney
		}
		err = UpdateDBModel(s, app)
		if err != nil {
			log.Println(err.Error())
			s.Rollback()
			return err
		}
		if app.Yue < 0 {
			err = fmt.Errorf("Not enough Yue Error")
			log.Println(err.Error())
			s.Rollback()
			return err
		}
		err = s.Commit()
		if err != nil {
			log.Println("!!!Commit Error:", err.Error())
			return err
		}
	}
	return nil
}