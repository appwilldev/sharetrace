package models

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/utils"
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

func GetAppuserMoneyTotalByUserid(s *ModelSession, appid string, appuserid string) (int64, int64, int64, error) {
	var (
		total      int64
		total_used int64
		total_left int64
		err        error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	sql := fmt.Sprintf("select sum(money) from appuser_money where appid = '%s' and appuserid = '%s' and (money_type = %d or money_type = %d or money_type = %d) ", appid, appuserid, conf.MONEY_TYPE_CLICK_SHARER, conf.MONEY_TYPE_INSTALL_INSTALLER, conf.MONEY_TYPE_INSTALL_SHARER)
	columnTypes := []reflect.Type{reflect.TypeOf(int64(1))}
	res, err := RawSqlQuery(sql, columnTypes)
	if err != nil {
		return -1, -1, -1, err
	}
	total = res[0][0].(int64) / 100

	sql = fmt.Sprintf("select sum(money) from appuser_money where appid = '%s' and appuserid = '%s' and money_type = %d ", appid, appuserid, conf.MONEY_TYPE_HFCZ)
	columnTypes = []reflect.Type{reflect.TypeOf(int64(1))}
	res, err = RawSqlQuery(sql, columnTypes)
	if err != nil {
		return -1, -1, -1, err
	}
	total_used = res[0][0].(int64) / 100

	total_left = total - total_used

	return total, total_used, total_left, nil
}

func AddClickAwardToAppUser(s *ModelSession, cs_data *ClickSession) error {

	if s == nil {
		s = NewModelSession()
	}
	defer s.Close()

	su_data, _ := GetShareURLById(nil, cs_data.Shareid)
	app, _ := GetAppInfoByAppid(nil, su_data.Appid)
	//if app.Status == 1 && app.Yue > 1000 {
	if app.Status == 1 {
		err := s.Begin()
		if app.ShareClickMoney > 0 {
			id, err := GenerateAppuserMoneyId()
			if err != nil {
				s.Rollback()
				return err
			}
			aum_data := new(AppuserMoney)
			aum_data.Id = id
			aum_data.Appid = app.Appid
			aum_data.Appuserid = su_data.Fromid
			aum_data.ClickSessionID = cs_data.Id
			aum_data.MoneyType = conf.MONEY_TYPE_INSTALL_SHARER
			aum_data.Money = float64(app.ShareClickMoney)
			aum_data.CreatedUTC = utils.GetNowSecond()
			aum_data.Des = "分享链接为App带来了点击:" + cs_data.AgentId
			err = InsertDBModel(s, aum_data)
			if err != nil {
				s.Rollback()
				return err
			}
			//app.Yue = app.Yue - app.ShareInstallMoney
		}

		err = s.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func AddInstallAwardToAppUser(s *ModelSession, app *AppInfo, cs_data *ClickSession) error {

	if s == nil {
		s = NewModelSession()
	}
	defer s.Close()

	su_data, _ := GetShareURLById(nil, cs_data.Shareid)
	//if app.Status == 1 && app.Yue > 1000 {
	if app.Status == 1 {
		err := s.Begin()
		if app.ShareInstallMoney > 0 {
			id, err := GenerateAppuserMoneyId()
			if err != nil {
				s.Rollback()
				return err
			}
			aum_data := new(AppuserMoney)
			aum_data.Id = id
			aum_data.Appid = app.Appid
			aum_data.Appuserid = su_data.Fromid
			aum_data.ClickSessionID = cs_data.Id
			aum_data.MoneyType = conf.MONEY_TYPE_INSTALL_SHARER
			aum_data.Money = float64(app.ShareInstallMoney)
			aum_data.CreatedUTC = utils.GetNowSecond()
			aum_data.Des = "分享链接为App带来了新用户" + cs_data.Installid
			err = InsertDBModel(s, aum_data)
			if err != nil {
				s.Rollback()
				return err
			}
			//app.Yue = app.Yue - app.ShareInstallMoney
		}

		if app.InstallMoney > 0 {
			id, err := GenerateAppuserMoneyId()
			if err != nil {
				s.Rollback()
				return err
			}
			aum_data := new(AppuserMoney)
			aum_data.Id = id
			aum_data.Appid = app.Appid
			aum_data.Appuserid = cs_data.Installid
			aum_data.ClickSessionID = cs_data.Id
			aum_data.MoneyType = conf.MONEY_TYPE_INSTALL_INSTALLER
			aum_data.CreatedUTC = utils.GetNowSecond()
			aum_data.Money = float64(app.InstallMoney)
			aum_data.Des = "通过点击分享链接安装了App"
			err = InsertDBModel(s, aum_data)
			if err != nil {
				s.Rollback()
				return err
			}
			//app.Yue = app.Yue - app.ShareInstallMoney
		}
		//err = UpdateDBModel(s, app)
		//if err != nil {
		//	s.Rollback()
		//	return err
		//}
		//if app.Yue < 0 {
		//	err = fmt.Errorf("Not enough Yue Error")
		//	s.Rollback()
		//	return err
		//}
		err = s.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func AddOrderToAppUser(s *ModelSession, appid string, appuserid string, phoneno string, cardnum string) error {

	if s == nil {
		s = NewModelSession()
	}
	defer s.Close()

	// 下单后，立马扣钱
	// TODO 检查是否有足够的余额,
	auo_id, err := GenerateAppuserOrderId()
	if err != nil {
		return err
	}

	auo_data := new(AppuserOrder)
	auo_data.Id = auo_id
	auo_data.Appid = appid
	auo_data.Appuserid = appuserid
	auo_data.OrderType = conf.ORDER_TYPE_HUAFEI
	auo_data.OrderMoney, _ = strconv.ParseFloat(cardnum, 64)
	auo_data.OrderStatus = conf.ORDER_STATUS_INIT
	auo_data.Phoneno = phoneno
	auo_data.Cardnum = cardnum
	auo_data.CreatedUTC = utils.GetNowSecond()
	auo_data.Des = "使用了账户余额充值话费"
	err = InsertDBModel(s, auo_data)
	if err != nil {
		s.Rollback()
		return err
	}

	aum_id, err := GenerateAppuserMoneyId()
	if err != nil {
		s.Rollback()
		return err
	}
	aum_data := new(AppuserMoney)
	aum_data.Id = aum_id
	aum_data.Appid = appid
	aum_data.Appuserid = appuserid
	aum_data.UserOrderID = auo_id
	aum_data.MoneyType = conf.MONEY_TYPE_HFCZ
	aum_data.Money, _ = strconv.ParseFloat(cardnum, 64)
	aum_data.Money = aum_data.Money * 100.0
	aum_data.CreatedUTC = utils.GetNowSecond()
	aum_data.Des = "使用了账户余额充值话费"
	err = InsertDBModel(s, aum_data)
	if err != nil {
		s.Rollback()
		return err
	}

	err = s.Commit()
	if err != nil {
		return err
	}
	return nil
}
