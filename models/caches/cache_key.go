package caches

import "fmt"

/**************************************************************
NOTE !!!!!:
 1. 不要修改这些函数和const值
 2. 新增cache key要确保和原有的cache key不冲突

**************************************************************/

// app_info
func getAppInfoCacheKey(id int64) string {
	return fmt.Sprintf("ai_%d", id)
}
func getAppIdCacheKeyByAppid(str string) string {
	return fmt.Sprintf("ai_id_by_appid_%s", str)
}

// user_info
func getUserInfoCacheKey(id int64) string {
	return fmt.Sprintf("ui_%d", id)
}
func getUserIdCacheKeyByEmail(str string) string {
	return fmt.Sprintf("ui_id_by_email_%s", str)
}

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

func getShareURLIdCacheKeyByTripleID(appid string, fromid string, itemid string) string {
	return fmt.Sprintf("su_id_by_3id_%s_%s_%s", appid, fromid, itemid)
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
