package models

import ()

func GenerateUserInfoId() (int64, error) {
	return generateSequenceValue("user_id")
}

type UserInfo struct {
	Id         int64  `xorm:"id INT PK NOT NULL" json:"id"`
	Email      string `xorm:"email  VARCHAR(256) NOT NULL" json:"email"`
	Passwd     string `xorm:"passwd VARCHAR(256) NOT NULL" json:"passwd"`
	Name       string `xorm:"name VARCHAR(256) DEFAULT NULL" json:"name"`
	Des        string `xorm:"des TEXT  DEFAULT NULL" json:"des"`
	Status     int    `xorm:"status INT NOT NULL" json:"status"`
	CreatedUTC int    `xorm:"created_utc INT NOT NULL" json:"created_utc"`
}

func (*UserInfo) TableName() string {
	return "user_info"
}

func (m *UserInfo) UniqueCond() (string, []interface{}) {
	var res = make([]interface{}, 1)
	res[0] = m.Id

	return "id=?", res
}

func GetUserInfoById(s *ModelSession, userId int64) (*UserInfo, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	user := &UserInfo{}
	has, err := s.Id(userId).Get(user)
	if !has || err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserInfoByEmail(s *ModelSession, email string) (*UserInfo, error) {
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	user := &UserInfo{}
	has, err := s.Where("email=?", email).Get(user)
	if !has || err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserInfoAll(s *ModelSession) ([]*UserInfo, int64, error) {
	var (
		total    int64
		dataList = make([]*UserInfo, 0)
		err      error
	)
	if s == nil {
		s = newAutoCloseModelsSession()
	}

	total, _ = s.Count(new(UserInfo))
	err = s.OrderBy("id desc").Find(&dataList)

	if err != nil {
		return nil, 0, err
	}

	return dataList, total, nil
}
