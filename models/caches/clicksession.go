package caches

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/appwilldev/sharetrace/cache"
	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/models"
	"github.com/bitly/go-simplejson"
)

func init() {
	registerJsonTypeInfo(&models.ClickSession{})
}

func GetClickSessionModelInfoById(id int64) (*models.ClickSession, error) {
	j, err := getJsonModelInfo(getClickSessionInfoCacheKey(id), true, 0)
	if j == nil {
		return nil, err
	}

	v := &models.ClickSession{}
	fillJsonModelInfo(v, j)

	return v, nil
}

func SaveClickSessionInfo(data *models.ClickSession) {
	jsonStr, _ := json.Marshal(data)
	cache.Set(conf.DEFAULT_CACHE_DB_NAME, getClickSessionInfoCacheKey(data.Id), string(jsonStr), 0)
}

func GetClickSessionJsonModelInfo(id int64) (*simplejson.Json, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getClickSessionInfoCacheKey(id))
	if err != nil {
		return nil, err
	}

	j, err := simplejson.NewFromReader(strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	return j, nil
}

func GetClickSessionIdByCookieid(idStr string) (string, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getClickSessionIdCacheKeyByCookieid(idStr))
	if err != nil {
		return "", err
	}
	return data, err
}

func UpdateClickSession(data *models.ClickSession) error {
	v, _ := json.Marshal(data)
	err := cache.Set(conf.DEFAULT_CACHE_DB_NAME, getClickSessionInfoCacheKey(data.Id), string(v), 0)
	if err != nil && conf.UserExpires > 0 {
		expiresDataCh <- expiresInfo{cacheKey: getClickSessionInfoCacheKey(data.Id), expires: conf.UserExpires}
	}

	return err
}

func NewClickSession(data *models.ClickSession) error {
	err := UpdateClickSession(data)
	if err != nil {
		return fmt.Errorf("Failed to cache link info %s", err.Error())
	}
	err = cache.Set(conf.DEFAULT_CACHE_DB_NAME, getClickSessionIdCacheKeyByCookieid(data.Cookieid), fmt.Sprintf("%d", data.Id), 0)
	if err != nil && conf.UserExpires > 0 {
		expiresDataCh <- expiresInfo{cacheKey: getClickSessionIdCacheKeyByCookieid(data.Cookieid), expires: conf.UserExpires}
	}

	return err
}
