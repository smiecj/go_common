package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func DoGetRequest(url string, parameters map[string]string) []byte {
	emptyArr := make([]byte, 0)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[DoGetRequest] 当前请求URL: %s, 初始化http对象失败原因: %s", url, err.Error())
		return emptyArr
	}

	query := req.URL.Query()

	for k, v := range parameters {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	startTime := time.Now()
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[DoGetRequest] 当前请求URL: %s, 失败原因: %s", url, err.Error())
		return emptyArr
	}
	defer rsp.Body.Close()

	rspBytes, _ := ioutil.ReadAll(rsp.Body)
	endTime := time.Now()
	log.Printf("[DoGetRequest] 请求URL: %s 成功，耗时: %d秒", url, endTime.Unix()-startTime.Unix())

	return rspBytes
}

func DoPostRequest(url string, parameters map[string]string) []byte {
	emptyArr := make([]byte, 0)
	jsonBytes, _ := json.Marshal(parameters)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Printf("[DoPostRequest] 当前请求URL: %s, 初始化http对象失败原因: %s", url, err.Error())
		return emptyArr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", "10056")

	client := &http.Client{
		Timeout:   time.Second * 60000,
		Transport: http.DefaultTransport,
	}

	startTime := time.Now()
	rsp, err := client.Do(req)
	if err != nil {
		log.Printf("[DoPostRequest] 当前请求URL: %s, 失败原因: %s", url, err.Error())
		return emptyArr
	}

	rspBytes, _ := ioutil.ReadAll(rsp.Body)
	endTime := time.Now()
	log.Printf("[DoPostRequest] 请求URL: %s 成功，耗时: %d秒", url, endTime.Unix()-startTime.Unix())

	return rspBytes
}
