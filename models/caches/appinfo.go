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
	registerJsonTypeInfo(&models.AppInfo{})
}

func GetAppInfoModelById(id int64) (*models.AppInfo, error) {
	j, err := getJsonModelInfo(getAppInfoCacheKey(id), true, 0)
	if j == nil {
		return nil, err
	}

	v := &models.AppInfo{}
	fillJsonModelInfo(v, j)

	return v, nil
}

func GetAppInfoIdByAppid(str string) (string, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getAppIdCacheKeyByAppid(str))
	if err != nil {
		return "", err
	}
	return data, err
}

func UpdateAppInfo(data *models.AppInfo) error {
	jsonStr, _ := json.Marshal(data)
	err := cache.Set(conf.DEFAULT_CACHE_DB_NAME, getAppInfoCacheKey(data.Id), string(jsonStr), 0)
	return err
}

func SetAppInfo(data *models.AppInfo) error {
	err := UpdateAppInfo(data)
	if err != nil {
		logger.ErrorLogger.Error(map[string]interface{}{
			"type":    "cache_data",
			"err_msg": err.Error(),
		})
		return fmt.Errorf("Failed to cache data info %s", err.Error())
	}
	err = cache.Set(conf.DEFAULT_CACHE_DB_NAME, getAppIdCacheKeyByAppid(data.Appid), fmt.Sprintf("%d", data.Id), 0)
	return err
}
