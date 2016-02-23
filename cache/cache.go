package cache

import (
	"fmt"
	"log"

	"strconv"

	"github.com/appwilldev/sharetrace/conf"
	redisPool "github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

var (
	ValueNilError = fmt.Errorf("redis_nil_value")
	Pools         = make(map[string]*redisPool.Pool)
)

func init() {
	var err error

	for name, config := range conf.RedConfig {
		Pools[name], err =
			redisPool.New(
				"tcp",
				fmt.Sprintf("%s:%d", config.Host, config.Port),
				256)
		if err != nil {
			log.Fatalf("Failed to init redis pool: " + err.Error())
		}
	}
}

func getConnFromPool(dbName string) (*redis.Client, error) {
	pool := Pools[dbName]
	if pool == nil {
		log.Println("No cache instance found for %s", dbName)
		return nil, fmt.Errorf("No cache instance found for %s", dbName)
	}

	return pool.Get()
}

func putConnToPool(conn *redis.Client, resp *redis.Resp, dbName string) {
	if resp.IsType(redis.IOErr) {
		log.Println("get redis connect timeout")
		Pools[dbName].Empty()
		return
	}

	Pools[dbName].Put(conn)
}

func Set(dbName string, key string, value string, expireTime int) error {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return err
	}

	var res *redis.Resp
	if expireTime > 0 {
		res = conn.Cmd("set", key, value, "EX", expireTime)
	} else {
		res = conn.Cmd("set", key, value)
	}
	defer putConnToPool(conn, res, dbName)

	if res.Err != nil {
		return res.Err
	}

	return res.Err
}

func Get(dbName string, key string) (string, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return "", err
	}

	res := conn.Cmd("get", key)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return "", res.Err
	}
	if res.IsType(redis.Nil) {
		return "", ValueNilError
	}

	return res.Str()
}

func IncrByCount(dbName string, key string, count int64) (int, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return 0, err
	}

	res := conn.Cmd("incrby", key, count)
	defer putConnToPool(conn, res, dbName)

	return res.Int()
}

func Expire(dbName string, key string, expireTime int) error {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return err
	}

	res := conn.Cmd("expire", key, expireTime)
	defer putConnToPool(conn, res, dbName)

	return res.Err
}

func Delete(dbName string, key string) error {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return err
	}

	res := conn.Cmd("del", key)
	defer putConnToPool(conn, res, dbName)

	return res.Err
}

func Exists(dbName string, key string) (bool, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return false, err
	}

	res := conn.Cmd("exists", key)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return false, res.Err
	}

	r, err := res.Int()
	if r == 0 || err != nil {
		return false, err
	}

	return true, nil
}

func ListPush(dbName string, key string, value string) error {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return err
	}

	res := conn.Cmd("lpush", key, value)
	defer putConnToPool(conn, res, dbName)

	return res.Err
}

func ListPop(dbName string, key string) error {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return err
	}

	res := conn.Cmd("rpop", key)
	defer putConnToPool(conn, res, dbName)

	return res.Err
}

func ListLen(dbName string, key string) (int, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return 0, err
	}

	res := conn.Cmd("llen", key)
	defer putConnToPool(conn, res, dbName)

	return res.Int()
}

func ListRange(dbName string, key string, start int, end int) ([]string, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return nil, err
	}

	res := conn.Cmd("lrange", key, start, end)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return nil, res.Err
	}

	respList, err := res.Array()
	if err != nil {
		return nil, err
	}

	resList := make([]string, len(respList))
	for i, resp := range respList {
		resList[i], _ = resp.Str()
	}

	return resList, nil
}

func ListElementDelete(dbName string, key string, value string) (err error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return err
	}

	res := conn.Cmd("lrem", key, 0, value)
	defer putConnToPool(conn, res, dbName)

	return res.Err
}

func HashGet(dbName string, hKey string, key string) (string, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return "", err
	}

	res := conn.Cmd("hget", hKey, key)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return "", res.Err
	}
	if res.IsType(redis.Nil) {
		return "", nil
	}

	return res.Str()
}

func SortedSetAdd(dbName string, key string, val string, score float64) error {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return err
	}

	res := conn.Cmd("zadd", key, score, val)
	defer putConnToPool(conn, res, dbName)

	return res.Err
}

func SortedSetElementDelete(dbName string, key string, val string) error {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return err
	}

	res := conn.Cmd("zrem", key, val)
	defer putConnToPool(conn, res, dbName)

	return res.Err
}

func SortedSetCards(dbName string, key string) (int, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return 0, err
	}

	res := conn.Cmd("zcard", key)
	defer putConnToPool(conn, res, dbName)

	return res.Int()
}

func SortedSetCardsBigInt(dbName string, key string) (int64, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return 0, err
	}

	res := conn.Cmd("zcard", key)
	defer putConnToPool(conn, res, dbName)

	return res.Int64()
}

func SortedSetElemIndex(dbName string, key string, elem string) (int, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return 0, err
	}

	res := conn.Cmd("zrevrank", key, elem)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return 0, res.Err
	}
	if res.IsType(redis.Nil) {
		return -1, nil
	}

	return res.Int()
}

func SortedSetScore(dbName string, key string, val string) (float64, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return -1, err
	}

	res := conn.Cmd("zscore", key, val)
	defer putConnToPool(conn, res, dbName)

	if res.Err != nil {
		return -1, res.Err
	}

	if res.IsType(redis.Nil) {
		return -1, ValueNilError
	}

	return res.Float64()
}

func SortedSetIncrScore(dbName string, key string, val string, score float64) (float64, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return -1, err
	}

	res := conn.Cmd("zincrby", key, score, val)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return -1, res.Err
	}

	return res.Float64()
}

func SortedSetExist(dbName string, key string, val string) (bool, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return false, err
	}

	res := conn.Cmd("zscore", key, val)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return false, res.Err
	}

	if res.IsType(redis.Nil) {
		return false, nil
	}

	return true, nil
}

func SortedSetRevRange(dbName string, key string, start int, end int) ([]string, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return nil, err
	}

	res := conn.Cmd("zrevrange", key, start, end)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return nil, res.Err
	}

	respList, err := res.Array()
	if err != nil {
		return nil, err
	}

	resList := make([]string, len(respList))
	for i, resp := range respList {
		resList[i], _ = resp.Str()
	}

	return resList, nil
}

func SortedSetRevRangeByScore(dbName string, cacheKey string, maxScore string, minScore string, count int) ([]string, []float64, error) {
	conn, err := getConnFromPool(dbName)
	if err != nil {
		return nil, nil, err
	}

	res := conn.Cmd("zrevrangebyscore", cacheKey, maxScore, minScore, "withscores", "limit", 0, count)
	defer putConnToPool(conn, res, dbName)
	if res.Err != nil {
		return nil, nil, fmt.Errorf("redis error: ", res.Err.Error())
	}

	redList, err := res.List()
	if err != nil {
		return nil, nil, fmt.Errorf("redis error: ", res.Err.Error())
	}

	// redis sorted set结构 ["id1", "score1", "id2", "score2"....]
	listLen := len(redList)
	idList := make([]string, listLen/2)
	scoreList := make([]float64, listLen/2)
	for ix := 0; ix < listLen; {
		idList[ix/2] = redList[ix]
		score, _ := strconv.ParseFloat(redList[ix+1], 10)
		scoreList[ix/2] = score
		ix += 2
	}

	return idList, scoreList, nil
}
