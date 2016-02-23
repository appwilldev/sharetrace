package caches

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/appwilldev/sharetrace/cache"
	"github.com/appwilldev/sharetrace/conf"
	"github.com/appwilldev/sharetrace/utils"
	"github.com/bitly/go-simplejson"
)

type expiresInfo struct {
	cacheKey string
	expires  int
}

var (
	invalidJsonFormatError = fmt.Errorf("not_json_format")
	expiresDataCh          = make(chan expiresInfo, 10240)
)

const (
	secondBase    int     = 1438600000
	oneDaySeconds float64 = 24 * 3600
)

func dumpExpiresTask(data map[string]int) {
	for cacheKey, expires := range data {
		cache.Expire(conf.DEFAULT_CACHE_DB_NAME, cacheKey, expires)
	}
}

func handleExpiresTask() {
	expiresData := make(map[string]int)
	dumpExpiresUTC := utils.GetNowSecond()
	expiresCount := 0

	for {
		expires := <-expiresDataCh
		expiresData[expires.cacheKey] = expires.expires
		expiresCount++

		if expiresCount%1000 == 0 && utils.GetNowSecond()-dumpExpiresUTC > conf.DumpExpiresDuration {
			go dumpExpiresTask(expiresData)
			expiresData = make(map[string]int)
			expiresCount = 0
			dumpExpiresUTC = utils.GetNowSecond()
		}
	}
}

func init() {
	go handleExpiresTask()
}

func randomFloat() float64 {
	s := rand.Float64()
	ss := strconv.FormatFloat(s, 'f', 5, 64)
	s, _ = strconv.ParseFloat(ss, 64)

	return s
}

func getJsonModelInfo(key string, deleteIfNotJson bool, expires int) (jsonRes *simplejson.Json, err error) {
	data, err := cache.Get(conf.DEFAULT_CACHE_DB_NAME, key)
	if err != nil {
		if err == cache.ValueNilError {
			return nil, cache.ValueNilError
		} else {
			log.Println("get_cache_info error, key:%s, err_msg:%s", key, err.Error())
			return
		}
	}

	r := strings.NewReader(data)
	jsonRes, err = simplejson.NewFromReader(r)
	if err != nil {
		log.Println("get_cache_info error, key:%s, err_msg:%s", key, err.Error())

		if deleteIfNotJson {
			deleteKey(key)
		}

		return nil, invalidJsonFormatError
	}

	if expires > 0 {
		expiresDataCh <- expiresInfo{cacheKey: key, expires: expires}
	}

	return
}

func fillJsonModelList(idList []string, keyGen func(string) string, deleteIfNotJson bool,
	listKey string, removeKeyFunc func(string, string) error,
	restoreFunc func(int64) (*simplejson.Json, error), elemExpires int) []*simplejson.Json {
	res := make([]*simplejson.Json, 0)

	for _, id := range idList {
		j, err := getJsonModelInfo(keyGen(id), deleteIfNotJson, elemExpires)
		if err == nil {
			res = append(res, j)
			continue
		}

		if err != invalidJsonFormatError && err != cache.ValueNilError {
			// 其他错误，不用从listKey中移除
			continue
		}

		if restoreFunc != nil {
			idNum, err := strconv.ParseInt(id, 10, 64)
			if err == nil {
				j, err = restoreFunc(idNum)
				if err != nil {
					//数据库错误，不能判断是不是过期，因此不能从列表中删除
					log.Println("restore_info_from_db error, key:%s, err_msg:%s", keyGen(id), err.Error())
					continue
				}
			} else {
				j = nil
			}
		}

		if j == nil {
			if restoreFunc == nil || removeKeyFunc == nil {
				// 不提供restore函数，不从列表中删除
				continue
			}
			// 数据库中也未找到，可以从列表中删除
			if _err := removeKeyFunc(listKey, id); _err != nil {
				log.Println("remove_id_from_list error, key:%s, err_msg:%s", listKey, _err.Error())
			}
		} else {
			res = append(res, j)
		}
	}

	return res
}

func getJsonModelListFromList(listKey string, start, end int, keyGen func(string) string, deleteIfNotJson bool,
	restoreFunc func(int64) (*simplejson.Json, error), elemExpires int) ([]*simplejson.Json, error) {
	idList, err := cache.ListRange(conf.DEFAULT_CACHE_DB_NAME, listKey, start, end)
	if err != nil {
		return nil, err
	}

	if deleteIfNotJson {
		return fillJsonModelList(idList, keyGen, true, listKey, removeListElement, restoreFunc, elemExpires), nil
	}

	return fillJsonModelList(idList, keyGen, false, "", nil, restoreFunc, elemExpires), nil
}

func pushElementList(key string, val string) error {
	return cache.ListPush(conf.DEFAULT_CACHE_DB_NAME, key, val)
}

func listCount(key string) int {
	n, err := cache.ListLen(conf.DEFAULT_CACHE_DB_NAME, key)
	if err != nil {
		//TODO:
		return 0
	}

	return n
}

func removeListElement(key string, val string) error {
	return cache.ListElementDelete(conf.DEFAULT_CACHE_DB_NAME, key, val)
}

func addSortedSetElement(key string, val string, score float64) error {
	return cache.SortedSetAdd(conf.DEFAULT_CACHE_DB_NAME, key, val, score)
}

func incrSortedSetElementScore(key string, val string, score float64) (float64, error) {
	return cache.SortedSetIncrScore(conf.DEFAULT_CACHE_DB_NAME, key, val, score)
}

func sortedSetElementScore(key string, val string) (float64, error) {
	return cache.SortedSetScore(conf.DEFAULT_CACHE_DB_NAME, key, val)
}

func removeSortedSetElement(key string, val string) error {
	return cache.SortedSetElementDelete(conf.DEFAULT_CACHE_DB_NAME, key, val)
}

func deleteKey(key string) error {
	return cache.Delete(conf.DEFAULT_CACHE_DB_NAME, key)
}

func sortedSetCountBigInt(key string) int64 {
	n, err := cache.SortedSetCardsBigInt(conf.DEFAULT_CACHE_DB_NAME, key)
	if err != nil {
		//TODO:
		return 0
	}

	return n
}

func sortedSetCount(key string) int {
	n, err := cache.SortedSetCards(conf.DEFAULT_CACHE_DB_NAME, key)
	if err != nil {
		//TODO:
		return 0
	}

	return n
}

// 顺序从sorted set中获取list: zrevrange
func getJsonModelListFromSortedSetSeq(listKey string, start, end int, jsonModelKeyGen func(string) string,
	deleteIfNotJson bool, restoreFunc func(int64) (*simplejson.Json, error), elemExpires int) ([]*simplejson.Json, error) {
	idList, err := cache.SortedSetRevRange(conf.DEFAULT_CACHE_DB_NAME, listKey, start, end)
	if err != nil {
		return nil, err
	}

	if deleteIfNotJson {
		return fillJsonModelList(idList, jsonModelKeyGen, true, listKey, removeSortedSetElement, restoreFunc, elemExpires), nil
	}

	return fillJsonModelList(idList, jsonModelKeyGen, false, "", nil, restoreFunc, elemExpires), nil
}

// 按照score从sorted set中获取list: zrevrange ... withscores
func getJsonModelListFromSortedSetByScore(listKey, maxScore, minScore string, offset, count int,
	keyGen func(string) string, deleteIfNotJson bool,
	restoreFunc func(int64) (*simplejson.Json, error), elemExpires int) ([]*simplejson.Json, error) {
	idList, _, err := cache.SortedSetRevRangeByScore(conf.DEFAULT_CACHE_DB_NAME, listKey, maxScore, minScore, count)
	if err != nil {
		return nil, err
	}

	if deleteIfNotJson {
		return fillJsonModelList(idList, keyGen, true, listKey, removeSortedSetElement, restoreFunc, elemExpires), nil
	}

	return fillJsonModelList(idList, keyGen, false, "", nil, restoreFunc, elemExpires), nil
}
