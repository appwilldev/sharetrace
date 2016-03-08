package caches

import (
	"encoding/json"
	"fmt"
	"github.com/appwilldev/sharetrace/cache"
	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/models"
)

func init() {
	registerJsonTypeInfo(&models.ShareURL{})
}

func GetShareURLModelInfoById(id int64) (*models.ShareURL, error) {
	j, err := getJsonModelInfo(getShareURLInfoCacheKey(id), true, conf.UserExpires)
	if j == nil {
		return nil, err
	}

	v := &models.ShareURL{}
	fillJsonModelInfo(v, j)

	return v, nil
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
	err := cache.Set(conf.DEFAULT_CACHE_DB_NAME, getShareURLInfoCacheKey(data.Id), string(v), conf.UserExpires)
	return err
}

func SetShareURL(data *models.ShareURL) error {
	err := UpdateShareURL(data)
	if err != nil {
		return fmt.Errorf("Failed to cache info %s", err.Error())
	}
	// Set url cache map
	err = cache.Set(conf.DEFAULT_CACHE_DB_NAME, getShareURLIdCacheKeyByURL(data.ShareURL), fmt.Sprintf("%d", data.Id), conf.UserExpires)
	return err
}
