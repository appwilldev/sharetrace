package caches

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/appwilldev/sharetrace/cache"
	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/models"
	"github.com/bitly/go-simplejson"
)

func init() {
	registerJsonTypeInfo(&models.ShareURL{})
}

func GetShareURLModelInfoById(id int64) (*models.ShareURL, error) {
	j, err := getJsonModelInfo(getShareURLInfoCacheKey(id), true, 0)
	if j == nil {
		return nil, err
	}

	v := &models.ShareURL{}
	fillJsonModelInfo(v, j)

	return v, nil
}

func SaveShareURLInfo(data *models.ShareURL) {
	jsonStr, _ := json.Marshal(data)
	cache.Set(conf.DEFAULT_CACHE_DB_NAME, getShareURLInfoCacheKey(data.Id), string(jsonStr), 0)
}

func GetShareURLJsonModelInfo(id int64) (*simplejson.Json, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getShareURLInfoCacheKey(id))
	if err != nil {
		return nil, err
	}

	j, err := simplejson.NewFromReader(strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	return j, nil
}

func GetShareURLIdByUrl(url string) (string, error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, getShareURLIdCacheKeyByURL(url))
	if err != nil {
		return "", err
	}
	return data, err
}

func UpdateShareURL(data *models.ShareURL) error {
	v, _ := json.Marshal(data)
	err := cache.Set(conf.DEFAULT_CACHE_DB_NAME, getShareURLInfoCacheKey(data.Id), string(v), 0)
	if err != nil && conf.UserExpires > 0 {
		expiresDataCh <- expiresInfo{cacheKey: getShareURLInfoCacheKey(data.Id), expires: conf.UserExpires}
	}

	return err
}

func NewShareURL(data *models.ShareURL) error {
	err := UpdateShareURL(data)
	if err != nil {
		log.Println("Failed to cache link info %s", err.Error())
		return fmt.Errorf("Failed to cache link info %s", err.Error())
	}
	// Set link url cache map
	err = cache.Set(conf.DEFAULT_CACHE_DB_NAME, getShareURLIdCacheKeyByURL(data.ShareURL), fmt.Sprintf("%d", data.Id), 0)
	if err != nil && conf.UserExpires > 0 {
		expiresDataCh <- expiresInfo{cacheKey: getShareURLIdCacheKeyByURL(data.ShareURL), expires: conf.UserExpires}
	}

	return err
}
