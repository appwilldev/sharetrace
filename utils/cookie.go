package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getCookieVersion(cookie string) int {
	var _signed_value_version_re = "^([1-9][0-9]*)\\|(.*)$"
	m, _ := regexp.MatchString(_signed_value_version_re, cookie)
	if m {
		return 2
	}

	return 1
}

var (
	CookieKey    string = "userid"
	cookieSecret string = "howareyoudoing"
)

func createSignature(h func() hash.Hash, parts ...string) string {
	c := hmac.New(h, []byte(cookieSecret))
	for _, v := range parts {
		c.Write([]byte(v))
	}

	return hex.EncodeToString(c.Sum(nil))
}

func decodeCookieV1(cookie string) int64 {
	// Tornado secure cookie v1
	if len(cookie) <= 0 {
		return 0
	}

	parts := strings.Split(cookie, "|")
	if len(parts) != 3 {
		return 0
	}

	sign := createSignature(sha1.New, CookieKey, parts[0], parts[1])
	if sign != parts[2] {
		return 0
	}

	d, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return 0
	}
	id, _ := strconv.ParseInt(string(d), 10, 64)
	return id
}

func decodeCookieV2(cookie string) int64 {
	// Tornado secure cookie v2
	if len(cookie) <= 0 {
		return 0
	}

	parts := strings.Split(cookie, "|")
	if len(parts) != 6 {
		return 0
	}

	signStr := cookie[0 : len(cookie)-len(parts[5])]
	sign := createSignature(sha256.New, signStr)
	if sign != parts[5] {
		return 0
	}

	d, err := base64.StdEncoding.DecodeString(parts[4][2:])
	if err != nil {
		return 0
	}
	id, _ := strconv.ParseInt(string(d), 10, 64)
	return id
}

func DecodeCookie(cookie string) int64 {
	version := getCookieVersion(cookie)
	switch version {
	case 1:
		return decodeCookieV1(cookie)
	case 2:
		return decodeCookieV2(cookie)
	}

	return 0
}

func NewCookie(userid int64) *http.Cookie {
	userIdStr := fmt.Sprintf("%d", userid)

	// Tornado cookie v1: three parts: base64(x)|timestamp|signature
	// 1: base64(x)
	e := base64.StdEncoding.EncodeToString([]byte(userIdStr))

	// 2: timestamp
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// 3: signature
	sign := createSignature(sha1.New, CookieKey, e, timestamp)
	value := e + "|" + timestamp + "|" + sign

	cookie := new(http.Cookie)
	cookie.Name = CookieKey
	cookie.Expires = time.Now().Add(time.Duration(365*86400) * time.Second)
	cookie.Value = value
	cookie.Path = "/"

	return cookie
}

func EmptyCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = CookieKey
	cookie.Path = "/"
	cookie.MaxAge = 0

	return cookie
}
