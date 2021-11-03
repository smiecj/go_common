package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/smiecj/go_common/util/log"
)

var (
	defaultClient = &http.Client{
		Timeout:   time.Second * 60000,
		Transport: http.DefaultTransport,
	}
)

func DoGetRequest(url string, parameters map[string]string) []byte {
	emptyArr := make([]byte, 0)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("[DoGetRequest] 当前请求URL: %s, 初始化http对象失败原因: %s", url, err.Error())
		return emptyArr
	}

	query := req.URL.Query()

	for k, v := range parameters {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	return commonSendRequest(req)
}

func DoPostRequest(requestUrl string, parameters map[string]string) []byte {
	emptyArr := make([]byte, 0)
	jsonBytes, _ := json.Marshal(parameters)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Error("[DoPostRequest] 当前请求URL: %s, 初始化http对象失败原因: %s", requestUrl, err.Error())
		return emptyArr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "10056")

	return commonSendRequest(req)
}

// 发起 post、form 格式请求
func DoPostFormRequest(requestUrl string, parameters map[string]string) []byte {
	emptyArr := make([]byte, 0)
	data := url.Values{}
	for key, value := range parameters {
		data.Set(key, value)
	}
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("[DoPostRequest] 当前请求URL: %s, 初始化http对象失败原因: %s", requestUrl, err.Error())
		return emptyArr
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return commonSendRequest(req)
}

func commonSendRequest(request *http.Request) []byte {
	emptyBytes := make([]byte, 0)
	startTime := time.Now()
	rsp, err := defaultClient.Do(request)
	if err != nil {
		log.Error("[http request] 请求URL: %s, 方法: %s, 失败原因: %s", request.URL.RawPath,
			request.Method, err.Error())
		return emptyBytes
	}

	rspBytes, _ := ioutil.ReadAll(rsp.Body)
	endTime := time.Now()
	log.Info("[http request] 请求URL: %s 成功，耗时: %d秒", request.URL.RawPath,
		endTime.Unix()-startTime.Unix())
	return rspBytes
}
