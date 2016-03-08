package caches

import (
	"encoding/json"
	"fmt"
	"github.com/appwilldev/sharetrace/cache"
	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/models"
)

func init() {
	registerJsonTypeInfo(&models.ClickSession{})
}

func GetClickSessionModelInfoById(id int64) (*models.ClickSession, error) {
	j, err := getJsonModelInfo(getClickSessionInfoCacheKey(id), true, conf.UserExpires)
	if j == nil {
		return nil, err
	}

	v := &models.ClickSession{}
	fillJsonModelInfo(v, j)

	return v, nil
}

func GetClickSessionId(clicktype int, parastr string) (string, error) {
	if clicktype == conf.CLICK_TYPE_COOKIE {
		return GetClickSessionIdByCookieid(parastr)
	} else if clicktype == conf.CLICK_TYPE_IP {
		return GetClickSessionIdByIP(parastr)
	}
	return "", fmt.Errorf("ClickType Error")
}

func GetClickSessionIdByCookieid(str string) (string, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getClickSessionIdCacheKeyByCookieid(str))
	if err != nil {
		return "", err
	}
	return data, err
}

// IP can not make sure unique, the latest will overwride the older
func GetClickSessionIdByIP(str string) (string, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getClickSessionIdCacheKeyByIP(str))
	if err != nil {
		return "", err
	}
	return data, err
}

func GetClickSessionIdByAgentId(str string) (string, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getClickSessionIdCacheKeyByAgentId(str))
	if err != nil {
		return "", err
	}
	return data, err
}

func UpdateClickSession(data *models.ClickSession) error {
	v, _ := json.Marshal(data)
	err := cache.Set(conf.DEFAULT_CACHE_DB_NAME, getClickSessionInfoCacheKey(data.Id), string(v), conf.UserExpires)

	return err
}

func SetClickSession(data *models.ClickSession) error {
	err := UpdateClickSession(data)
	if err != nil {
		return fmt.Errorf("Failed to cache clicksession info %s", err.Error())
	}
	if data.Cookieid != "" {
		err = cache.Set(conf.DEFAULT_CACHE_DB_NAME, getClickSessionIdCacheKeyByCookieid(data.Cookieid), fmt.Sprintf("%d", data.Id), conf.UserExpires)
	}
	if data.AgentIP != "" {
		err = cache.Set(conf.DEFAULT_CACHE_DB_NAME, getClickSessionIdCacheKeyByIP(data.AgentIP), fmt.Sprintf("%d", data.Id), conf.UserExpires)
	}
	if data.AgentId != "" {
		err = cache.Set(conf.DEFAULT_CACHE_DB_NAME, getClickSessionIdCacheKeyByAgentId(data.AgentId), fmt.Sprintf("%d", data.Id), conf.UserExpires)
	}

	return err
}
