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

func getShareURLInfoCacheKeyStr(idstr string) string {
	return "sui_" + idstr
}

func getShareURLIdCacheKeyByURL(url string) string {
	return fmt.Sprintf("su_id_by_url_%s", url)
}

// clicksession
func getClickSessionInfoCacheKey(id int64) string {
	return fmt.Sprintf("csi_%d", id)
}

func getClickSessionInfoCacheKeyStr(idStr string) string {
	return "csi_" + idStr
}

func getClickSessionIdCacheKeyByCookieid(idStr string) string {
	return fmt.Sprintf("cs_id_by_cookieid_%s", idStr)
}

func getClickSessionIdCacheKeyByIP(IPStr string) string {
	return fmt.Sprintf("cs_id_by_agentip_%s", IPStr)
}

func getClickSessionIdCacheKeyByAgentId(idStr string) string {
	return fmt.Sprintf("cs_id_by_agentid_%s", idStr)
}
