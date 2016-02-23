package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	math_rand "math/rand"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

func init() {
	math_rand.Seed(time.Now().Unix())
}

func GenerateUUID() string {
	u := feeds.NewUUID()
	return string(u.String())
}

func GetNowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond/time.Nanosecond)
}

func GetNowSecond() int {
	return int(time.Now().UnixNano() / int64(time.Second/time.Nanosecond))
}

func GetNowStringYMD() string {
	return time.Now().Format("2006-01-02")
}

func HttpGet(uri string, cookies []*http.Cookie) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to call [%s], status code: %d", uri, res.StatusCode)
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func HttpPost(uri string, bodyData []byte) ([]byte, error) {
	b := bytes.NewReader(bodyData)
	res, err := http.Post(uri, "plain/text", b)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to call [%s], status code: %d", uri, res.StatusCode)
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func RunesTruncate(source string, count int) string {
	rs := bytes.Runes([]byte(source))
	if count > len(rs) {
		count = len(rs)
	}

	return string(rs[:count])
}

func FormPost(uri, form string, cookies []*http.Cookie) ([]byte, error) {
	buf := bytes.NewBufferString(form)
	req, err := http.NewRequest("POST", uri, buf)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookies {
		req.AddCookie(c)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func MergeSortedListResult(idList [][]string, scoreList [][]float64, count int, cmp func(float64, float64) bool) []string {
	res := make([]string, 0)
	insertedCount := 0
	indexMap := make(map[int]int)
	for ix := range idList {
		indexMap[ix] = 0
	}

	for {
		if len(indexMap) == 0 || insertedCount == count {
			break
		}

		var selectedScore float64
		var selectedScoreIx int
		var resStr string
		if insertedCount > 0 {
			for ix, index := range indexMap {
				if cmp(scoreList[ix][index], selectedScore) {
					selectedScore = scoreList[ix][index]
					selectedScoreIx = ix
					resStr = idList[ix][index]
				}
			}
		} else {
			for ix, index := range indexMap {
				if ix == 0 {
					selectedScore = scoreList[ix][index]
					selectedScoreIx = ix
					resStr = idList[ix][index]
				} else {
					if cmp(scoreList[ix][index], selectedScore) {
						selectedScore = scoreList[ix][index]
						selectedScoreIx = ix
						resStr = idList[ix][index]
					}
				}
			}
		}

		indexMap[selectedScoreIx]++
		if indexMap[selectedScoreIx] >= len(idList[selectedScoreIx]) {
			delete(indexMap, selectedScoreIx)
		}
		res = append(res, resStr)
		insertedCount++
	}

	return res
}
