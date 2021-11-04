package client

import (
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

const (
	MethodGet  HTTPRequestMethod = "GET"
	MethodPost HTTPRequestMethod = "POST"

	headerContentType     = "Content-Type"
	contentTypeJson       = "application/json"
	contentTypeUrlEncoded = "application/x-www-form-urlencoded"
)

type HTTPRequestMethod string

// http 客户端定义
type Client interface {
	Do(...ConfigRequestFunc) (*Response, error)
}

// 请求体
type Request struct {
	url    string
	header map[string]string
	param  map[string]string
	body   string
	method string
}

// 获取一个新request
// 注意对map 数据结构的初始化
func buildRequest() *Request {
	req := new(Request)
	req.header = make(map[string]string)
	req.param = make(map[string]string)
	return req
}

// 构造请求结构体
// 需要借助 method + header 共同构建
func (req *Request) makeRequestBody() {
	switch req.method {
	case string(MethodGet):
		// GET 方法统一不需要设置body
		break
	case string(MethodPost):
		// POST 方法，要根据header 中的参数格式 Content-Type 来设置格式
		// 如果body 已经被设置，则直接跳过
		if req.body != "" {
			break
		}
		if req.header[headerContentType] == contentTypeJson {
			jsonBytes, _ := json.Marshal(req.param)
			req.body = string(jsonBytes)
		} else if req.header[headerContentType] == contentTypeUrlEncoded {
			data := url.Values{}
			for key, value := range req.param {
				data.Set(key, value)
			}
			req.body = data.Encode()
		}
	}
}

type Response struct {
	RspBody   string
	RspHeader map[string]string
}

// 配置请求函数
type ConfigRequestFunc func(*Request)

// 配置url
func ConfigRequestUrl(url string) func(*Request) {
	return func(request *Request) {
		request.url = url
	}
}

// 配置 自定义 header，每次只设置一个key-value对
func ConfigRequestAddHeader(key, value string) func(*Request) {
	return func(request *Request) {
		if nil == request.header {
			request.header = make(map[string]string)
		}
		request.header[key] = value
	}
}

// 配置 param，每次只设置一个 key-value 对
func ConfigRequestAddParam(key, value string) func(*Request) {
	return func(request *Request) {
		request.param[key] = value
	}
}

// 配置 param，直接配置整个body
func ConfigRequestSetParam(body string) func(*Request) {
	return func(request *Request) {
		request.body = body
	}
}

// 配置 method
func ConfigRequestMethod(method HTTPRequestMethod) func(*Request) {
	return func(request *Request) {
		request.method = string(method)
		// POST 请求默认设置为 json body 格式
		request.header[headerContentType] = contentTypeJson
	}
}

// 配置 post 请求方式为 urlencode
// 注意这个方法需要在 ConfigRequestMethod 后执行
func ConfigRequestSetPostUrlEncode() func(*Request) {
	return func(request *Request) {
		request.header[headerContentType] = contentTypeUrlEncoded
	}
}

// http 客户端实现
type httpClient struct {
	Client *http.Client
}

// 获取 http 客户端
func GetHTTPClient() Client {
	httpClient := new(httpClient)
	httpClient.Client = defaultClient

	return httpClient
}

func (client *httpClient) Do(configFuncArr ...ConfigRequestFunc) (rsp *Response, err error) {
	request := buildRequest()
	rsp = new(Response)
	for _, currentConfigFunc := range configFuncArr {
		currentConfigFunc(request)
	}

	err = client.commonSendRequest(request, rsp)
	return
}

func (client *httpClient) commonSendRequest(request *Request, response *Response) (err error) {
	req, err := http.NewRequest(request.method, request.url, strings.NewReader(request.body))
	if err != nil {
		log.Error("[http.client.Do] 当前请求URL: %s, 请求 method: %s, 初始化http对象失败原因: %s",
			request.url, request.method, err.Error())
		return
	}

	if request.method == string(MethodGet) {
		query := req.URL.Query()
		for k, v := range request.param {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}
	startTime := time.Now()
	rsp, err := defaultClient.Do(req)
	if err != nil {
		log.Error("[http.client.Do] 请求URL: %s, 方法: %s, 失败原因: %s", request.url,
			request.method, err.Error())
		return err
	}

	rspBytes, _ := ioutil.ReadAll(rsp.Body)
	endTime := time.Now()
	log.Info("[http.client.Do] 请求URL: %s 成功，耗时: %d 秒", request.url,
		endTime.Unix()-startTime.Unix())
	response.RspBody = string(rspBytes)
	response.RspHeader = make(map[string]string)

	for key, valueArr := range rsp.Header {
		response.RspHeader[key] = valueArr[0]
	}

	return nil
}
