package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	//"strconv"
)

const (
	SERVER_ERROR = iota
	BAD_REQUEST
	BAD_POST_DATA
	LOGIN_NEEDED
	LOGIN_FAILED
	NOT_PERMITTED
	DATA_NOT_FOUND
	REGISTER_FAILED
	DATA_DUPLICATED
)

type RequestLogData struct {
	Status bool
	Error  string
	Msg    string
}

var (
	errorStr = map[int][2]string{
		SERVER_ERROR:    [2]string{"sever_error", "服务器错误"},
		BAD_REQUEST:     [2]string{"bad_request", "客户端请求错误"},
		BAD_POST_DATA:   [2]string{"bad_post_data", "客户端请求体错误"},
		LOGIN_NEEDED:    [2]string{"login_needed", "未登录"},
		LOGIN_FAILED:    [2]string{"login_failed", "登录失败"},
		NOT_PERMITTED:   [2]string{"not_permitted", "无权进行此次操作"},
		DATA_NOT_FOUND:  [2]string{"data_not_fond", "没有找到该数据"},
		REGISTER_FAILED: [2]string{"register_failed", "注册失败, 登录名或邮箱重复"},
		DATA_DUPLICATED: [2]string{"data_duplicated", "数据重复"},
	}
)

func Success(c *gin.Context, data interface{}) {
	res := gin.H{"status": true}
	if data != nil {
		res["data"] = data
	}

	c.Set("request_log", &RequestLogData{Status: true})
	c.JSON(200, res)
}

func Error(c *gin.Context, errorCode int, data ...interface{}) {
	var (
		errCodeStr = errorStr[errorCode][0]
		errMsg     = errorStr[errorCode][1]
		errMsgLog  = errMsg
	)

	if len(data) >= 1 {
		if data[0] != nil {
			errMsg = data[0].(string)
		}
		if len(data) >= 2 {
			if data[1] != nil {
				errMsgLog = data[1].(string)
			} else {
				errMsgLog = errMsg
			}
		}
	}

	log.Println("api_request code:", errCodeStr, "url:", c.Request.URL.Path, "err_msg:", errMsgLog)

	res := gin.H{"status": false, "code": errCodeStr, "msg": errMsg}
	c.Set("request_log", &RequestLogData{Status: false, Error: errCodeStr, Msg: errMsgLog})
	c.JSON(200, res)
}

func getUserIdFromContext(c *gin.Context) int64 {
	s, _ := c.Get("userid")
	//ck, err := c.Request.Cookie("userid")
	//s, err := strconv.ParseInt(ck.Value, 10, 64)
	//if err != nil {
	//	Error(c, LOGIN_NEEDED, err.Error())
	//	c.Abort()
	//	return -1
	//}
	//return s
	log.Println("getUserIdFromContext userid:", s)
	return s.(int64)
}
