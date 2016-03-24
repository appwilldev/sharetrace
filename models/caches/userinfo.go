package caches

import (
	"encoding/json"
	"fmt"

	"github.com/appwilldev/sharetrace/cache"
	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/logger"
	"github.com/appwilldev/sharetrace/models"
)

func init() {
	registerJsonTypeInfo(&models.UserInfo{})
}

func GetUserModelInfoById(id int64) (*models.UserInfo, error) {
	j, err := getJsonModelInfo(getUserInfoCacheKey(id), true, 0)
	if j == nil {
		return nil, err
	}

	v := &models.UserInfo{}
	fillJsonModelInfo(v, j)

	return v, nil
}

func GetUserInfoIdByEmail(str string) (string, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getUserIdCacheKeyByEmail(str))
	if err != nil {
		return "", err
	}
	return data, err
}

func UpdateUserInfo(data *models.UserInfo) error {
	jsonStr, _ := json.Marshal(data)
	err := cache.Set(conf.DEFAULT_CACHE_DB_NAME, getUserInfoCacheKey(data.Id), string(jsonStr), 0)
	return err
}

func SetUserInfo(data *models.UserInfo) error {
	err := UpdateUserInfo(data)
	if err != nil {
		logger.ErrorLogger.Error(map[string]interface{}{
			"type":    "cache_data",
			"err_msg": err.Error(),
		})
		return fmt.Errorf("Failed to cache data info %s", err.Error())
	}
	err = cache.Set(conf.DEFAULT_CACHE_DB_NAME, getUserIdCacheKeyByEmail(data.Email), fmt.Sprintf("%d", data.Id), 0)
	return err
}
