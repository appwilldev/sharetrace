package models

import (
	"log"
)

func GenerateAppuserOrderId() (int64, error) {
	return generateSequenceValue("appuser_order_id")
}

type AppuserOrder struct {
	Id        int64  `xorm:"id INT PK NOT NULL" json:"id"`
	Appid     string `xorm:"appid VARCHAR(256) NOT NULL" json:"appid"`
	Appuserid string `xorm:"appuserid VARCHAR(256) NOT NULL" json:"appuserid"`

	OrderType   int     `xorm:"order_type INT DEFAULT 0" json:"order_type"`
	OrderMoney  float64 `xorm:"order_money INT DEFAULT 0" json:"order_money"`
	OrderStatus int     `xorm:"order_status INT DEFAULT 0" json:"order_status"`

	SponderId    string `xorm:"sponder_id VARCHAR(256)" json:"sponder_id"`
	Phoneno      string `xorm:"phoneno VARCHAR(256) NOT NULL" json:"phoneno"`
	Cardnum      string `xorm:"cardnum VARCHAR(256) NOT NULL" json:"cardnum"`
	OrderRetInfo string `xorm:"order_ret_info TEXT NOT NULL" json:"order_ret_info"`

	Des        string `xorm:"des TEXT  DEFAULT NULL" json:"des"`
	Status     int    `xorm:"status INT NOT NULL" json:"status"`
	CreatedUTC int    `xorm:"created_utc INT NOT NULL" json:"created_utc"`
}

func (*AppuserOrder) TableName() string {
	return "appuser_order"
}

func (m *AppuserOrder) UniqueCond() (string, []interface{}) {
	var res = make([]interface{}, 1)
	res[0] = m.Id

	return "id=?", res
}

func GetAppuserOrderById(s *ModelSession, id int64) (*AppuserOrder, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	app := &AppuserOrder{}
	has, err := s.Id(id).Get(app)
	if !has || err != nil {
		return nil, err
	}

	return app, nil
}

func GetAppuserOrderListByUserid(s *ModelSession, appid string, appuserid string) ([]*AppuserOrder, int64, error) {
	var (
		total    int64
		dataList = make([]*AppuserOrder, 0)
		err      error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	total, _ = s.Count(new(AppuserOrder))
	err = s.Where("appid=? and appuserid=? ", appid, appuserid).OrderBy("id desc").Find(&dataList)

	if err != nil {
		return nil, 0, err
	}

	return dataList, total, nil
}

func GetAppuserOrderListByOrderStatus(s *ModelSession, status int) ([]*AppuserOrder, int64, error) {
	var (
		total    int64
		dataList = make([]*AppuserOrder, 0)
		err      error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	total, _ = s.Count(new(AppuserOrder))
	err = s.Where("order_status=?", status).OrderBy("id desc").Find(&dataList)

	if err != nil {
		return nil, 0, err
	}

	return dataList, total, nil
}

func AddConsumeToAppUser(s *ModelSession, app *AppInfo, appuserid string, phoneno string, cardnum string) {
	if s == nil {
		s = NewModelSession()
	}
	defer s.Close()

	err := s.Commit()
	if err != nil {
		log.Println("!!!Commit Error:", err.Error())
		return
	}
}
