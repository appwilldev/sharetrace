package caches

import "fmt"

/**************************************************************
NOTE !!!!!:
 1. 不要修改这些函数和const值
 2. 新增cache key要确保和原有的cache key不冲突

**************************************************************/

const (
	TEST_QUEUE = "test_queue"
)

// sharetrace
func getShareURLInfoCacheKey(id int64) string {
	return fmt.Sprintf("sui_%d", id)
}

func getShareURLInfoCacheKeyStr(str string) string {
	return "sui_" + str
}

func getShareURLIdCacheKeyByURL(str string) string {
	return fmt.Sprintf("su_id_by_url_%s", str)
}

// clicksession
func getClickSessionInfoCacheKey(id int64) string {
	return fmt.Sprintf("csi_%d", id)
}

func getClickSessionInfoCacheKeyStr(str string) string {
	return "csi_" + str
}

func getClickSessionIdCacheKeyByCookieid(str string) string {
	return fmt.Sprintf("cs_id_by_cookieid_%s", str)
}

func getClickSessionIdCacheKeyByIP(str string) string {
	return fmt.Sprintf("cs_id_by_agentip_%s", str)
}

func getClickSessionIdCacheKeyByAgentId(str string) string {
	return fmt.Sprintf("cs_id_by_agentid_%s", str)
}
