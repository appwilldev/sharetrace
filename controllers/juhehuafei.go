package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//----------------------------------
// 手机话费充值调用示例代码 － 聚合数据
// 在线接口文档：http://www.juhe.cn/docs/85
//----------------------------------

const OPENID_JUHE = "JHe639b4ea7c06d3513125eaea4aea95ce"
const APPKEY_HUAFEI = "44636e2c2810178b075b3344fd8d2d4d" //您申请的APPKEY

func Remaining(c *gin.Context) {
	JHHFYue()
}

func TelCheck(c *gin.Context) {
	q := c.Request.URL.Query()
	phoneno := q["phoneno"][0]
	cardnum := q["cardnum"][0]
	JHHFTelCheck(phoneno, cardnum)
	JHHFTelQuery(phoneno, cardnum)
}

//1.账户余额查询
func JHHFYue() {
	//请求地址
	juheURL := "http://op.juhe.cn/ofpay/mobile/yue"

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	time_now := time.Now()
	timestamp := fmt.Sprintf("%d", time_now.Unix())
	log.Println("timestamp:", timestamp)
	param.Set("timestamp", timestamp) //当前时间戳，如：1432788379
	param.Set("key", APPKEY_HUAFEI)   //应用APPKEY_HUAFEI(应用详细页查询)

	md5Ctx := md5.New()
	md5_org_str := fmt.Sprintf("%s%s%s", OPENID_JUHE, APPKEY_HUAFEI, timestamp)
	log.Println("md5_org_str:", md5_org_str)
	md5Ctx.Write([]byte(md5_org_str))
	sign := hex.EncodeToString(md5Ctx.Sum(nil))

	param.Set("sign", sign) //校验值

	//发送请求
	data, err := Get(juheURL, param)
	if err != nil {
		fmt.Errorf("请求失败,错误信息:\r\n%v", err)
	} else {
		var netReturn map[string]interface{}
		json.Unmarshal(data, &netReturn)
		if netReturn["error_code"].(float64) == 0 {
			fmt.Printf("接口返回result字段是:\r\n%v", netReturn["result"])
		}
	}
}

//2.订单状态查询
//{
//    "reason": "查询成功",
//    "result": {
//        "uordercash": "5.000", /*订单扣除金额*/
//        "sporder_id": "20150511163237508",/*聚合订单号*/
//        "game_state": "1" /*状态 1:成功 9:失败 0：充值中*/
//    },
//    "error_code": 0
//}
func JHHFOrderSta(orderid string) {
	//请求地址
	juheURL := "http://op.juhe.cn/ofpay/mobile/ordersta"

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("orderid", orderid)   //商家订单号，8-32位字母数字组合
	param.Set("key", APPKEY_HUAFEI) //应用APPKEY_HUAFEI(应用详细页查询)

	//发送请求
	data, err := Get(juheURL, param)
	if err != nil {
		fmt.Errorf("请求失败,错误信息:\r\n%v", err)
	} else {
		var netReturn map[string]interface{}
		json.Unmarshal(data, &netReturn)
		if netReturn["error_code"].(float64) == 0 {
			fmt.Printf("接口返回result字段是:\r\n%v", netReturn["result"])
		}
	}
}

//3.历史订单列表检索
func JHHFOrderList() {
	//请求地址
	juheURL := "http://op.juhe.cn/ofpay/mobile/orderlist"

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("page", "")           //当前页数，默认1
	param.Set("pagesize", "")       //每页显示条数，默认50，最大100
	param.Set("mobilephone", "")    //检索指定手机号码
	param.Set("orderid", "")        //需要检索的商户订单号
	param.Set("key", APPKEY_HUAFEI) //应用APPKEY_HUAFEI(应用详细页查询)

	//发送请求
	data, err := Get(juheURL, param)
	if err != nil {
		fmt.Errorf("请求失败,错误信息:\r\n%v", err)
	} else {
		var netReturn map[string]interface{}
		json.Unmarshal(data, &netReturn)
		if netReturn["error_code"].(float64) == 0 {
			fmt.Printf("接口返回result字段是:\r\n%v", netReturn["result"])
		}
	}
}

//4.状态回调配置
//func JHHF() {
//	//请求地址
//	juheURL := "充值接口测试完毕，联系在线客服进行回调配置"
//
//	//初始化参数
//	param := url.Values{}
//
//	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
//
//	//发送请求
//	data, err := Post(juheURL, param)
//	if err != nil {
//		fmt.Errorf("请求失败,错误信息:\r\n%v", err)
//	} else {
//		var netReturn map[string]interface{}
//		json.Unmarshal(data, &netReturn)
//		if netReturn["error_code"].(float64) == 0 {
//			fmt.Printf("接口返回result字段是:\r\n%v", netReturn["result"])
//		}
//	}
//}

//5.检测手机号码是否能充值
//{
//    "reason": "允许充值的手机号码及金额",
//    "result": null,
//    "error_code": 0
//}
func JHHFTelCheck(phoneno string, cardnum string) {
	//请求地址
	juheURL := "http://op.juhe.cn/ofpay/mobile/telcheck"

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("phoneno", phoneno)   //手机号码
	param.Set("cardnum", cardnum)   //充值金额,目前可选：5、10、20、30、50、100、300
	param.Set("key", APPKEY_HUAFEI) //应用APPKEY_HUAFEI(应用详细页查询)

	//发送请求
	data, err := Get(juheURL, param)
	if err != nil {
		fmt.Errorf("请求失败,错误信息:\r\n%v", err)
	} else {
		var netReturn map[string]interface{}
		json.Unmarshal(data, &netReturn)
		if netReturn["error_code"].(float64) == 0 {
			fmt.Printf("接口返回:\r\n%v", netReturn)
			//fmt.Printf("接口返回result字段是:\r\n%v", netReturn["result"])
		}
	}
}

//6.根据手机号和面值查询商品信息
func JHHFTelQuery(phoneno string, cardnum string) {
	//请求地址
	juheURL := "http://op.juhe.cn/ofpay/mobile/telquery"

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("phoneno", phoneno)   //手机号码
	param.Set("cardnum", cardnum)   //充值金额,目前可选：5、10、20、30、50、100、300
	param.Set("key", APPKEY_HUAFEI) //应用APPKEY_HUAFEI(应用详细页查询)

	//发送请求
	data, err := Get(juheURL, param)
	if err != nil {
		fmt.Errorf("请求失败,错误信息:\r\n%v", err)
	} else {
		var netReturn map[string]interface{}
		json.Unmarshal(data, &netReturn)
		if netReturn["error_code"].(float64) == 0 {
			fmt.Printf("接口返回result字段是:\r\n%v", netReturn["result"])
		}
	}
}

//7.手机直充接口
//{
//    "reason": "订单提交成功，等待充值",
//    "result": {
//        "cardid": "1900212", /*充值的卡类ID*/
//        "cardnum": "1", /*数量*/
//        "ordercash": 49.25, /*进货价格*/
//        "cardname": "江苏电信话费50元直充", /*充值名称*/
//        "sporder_id": "20141120174602882", /*聚合订单号*/
//        "uorderid":"2014123115121",/*商户自定的订单号*/
//        "game_userid": "18913515122", /*充值的手机号码*/
//        "game_state": "0" /*充值状态:0充值中 1成功 9撤销，刚提交都返回0*/
//    },
//    "error_code": 0
//}
func JHHFOnlineOrder(phoneno string, cardnum string, orderid string) {
	//请求地址
	juheURL := "http://op.juhe.cn/ofpay/mobile/onlineorder"

	//初始化参数
	param := url.Values{}

	//配置请求参数,方法内部已处理urlencode问题,中文参数可以直接传参
	param.Set("phoneno", "")        //手机号码
	param.Set("cardnum", "")        //充值金额,目前可选：5、10、20、30、50、100、300
	param.Set("orderid", "")        //商家订单号，8-32位字母数字组合
	param.Set("key", APPKEY_HUAFEI) //应用APPKEY_HUAFEI(应用详细页查询)

	md5Ctx := md5.New()
	md5_org_str := fmt.Sprintf("%s%s%s%s%s", OPENID_JUHE, APPKEY_HUAFEI, phoneno, cardnum, orderid)
	log.Println("md5_org_str:", md5_org_str)
	md5Ctx.Write([]byte(md5_org_str))
	sign := hex.EncodeToString(md5Ctx.Sum(nil))

	param.Set("sign", sign) //校验值，md5(OpenID+key+phoneno+cardnum+orderid)

	//发送请求
	data, err := Get(juheURL, param)
	if err != nil {
		fmt.Errorf("请求失败,错误信息:\r\n%v", err)
	} else {
		var netReturn map[string]interface{}
		json.Unmarshal(data, &netReturn)
		if netReturn["error_code"].(float64) == 0 {
			fmt.Printf("接口返回result字段是:\r\n%v", netReturn["result"])
		}
	}
}

// get 网络请求
func Get(apiURL string, params url.Values) (rs []byte, err error) {
	var Url *url.URL
	Url, err = url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	resp, err := http.Get(Url.String())
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// post 网络请求 ,params 是url.Values类型
func Post(apiURL string, params url.Values) (rs []byte, err error) {
	resp, err := http.PostForm(apiURL, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
