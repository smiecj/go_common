package http

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
)

var (
	defaultClient = &http.Client{
		Timeout:   time.Second * 60000,
		Transport: http.DefaultTransport,
	}
	defaultClientOnce sync.Once
	clientSingleton   Client
)

const (
	methodGet  HTTPRequestMethod = "GET"
	methodPost HTTPRequestMethod = "POST"

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
	ctx    context.Context
	cancel context.CancelFunc
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
	req.ctx, req.cancel = context.Background(), func() {}
	req.header = make(map[string]string)
	req.param = make(map[string]string)
	req.method = string(methodGet)
	return req
}

// 构造请求结构体
// 需要借助 method + header 共同构建
func (req *Request) buildRequestBody() {
	switch req.method {
	case string(methodGet):
		// GET 方法统一不需要设置body
		break
	case string(methodPost):
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

// 检查 request 参数是否合法
func (req *Request) checkIsValid() error {
	if req.url == "" {
		return errorcode.BuildErrorWithMsg(errorcode.NetHandleFailed, "url is not valid: "+req.url)
	}
	if req.method != string(methodGet) && req.method != string(methodPost) {
		return errorcode.BuildErrorWithMsg(errorcode.NetHandleFailed, "request method is not valid: "+req.method)
	}
	return nil
}

type Response struct {
	Body   string
	Header map[string]string
}

// 配置请求函数
type ConfigRequestFunc func(*Request)

// 配置 timeout context
func Timeout(timeout time.Duration) func(*Request) {
	return func(request *Request) {
		request.ctx, request.cancel = context.WithTimeout(request.ctx, timeout)
	}
}

// 配置url
func Url(url string) func(*Request) {
	return func(request *Request) {
		request.url = url
	}
}

// 配置 自定义 header，每次只设置一个key-value对
func AddHeader(key, value string) func(*Request) {
	return func(request *Request) {
		if nil == request.header {
			request.header = make(map[string]string)
		}
		request.header[key] = value
	}
}

// 配置 param，每次只设置一个 key-value 对
func AddParam(key, value string) func(*Request) {
	return func(request *Request) {
		request.param[key] = value
	}
}

// 配置 param
func SetParam(body string) func(*Request) {
	return func(request *Request) {
		request.body = body
	}
}

// 配置 param
func SetBody(body interface{}) func(*Request) {
	return func(request *Request) {
		bodyStr, _ := json.Marshal(body)
		request.body = string(bodyStr)
		request.header[headerContentType] = contentTypeJson
	}
}

// 配置 method get
func Get() func(*Request) {
	return func(request *Request) {
		request.method = string(methodGet)
		// POST 请求默认设置为 json body 格式
		request.header[headerContentType] = contentTypeJson
	}
}

// 配置 method post
func Post() func(*Request) {
	return func(request *Request) {
		request.method = string(methodPost)
		// POST 请求默认设置为 json body 格式
		request.header[headerContentType] = contentTypeJson
	}
}

// 配置 post 请求方式为 urlencode
// 注意这个方法需要在 ConfigRequestMethod 后执行
func PostWithUrlEncode() func(*Request) {
	return func(request *Request) {
		request.method = string(methodPost)
		request.header[headerContentType] = contentTypeUrlEncoded
	}
}

// http 客户端实现
type httpClient struct {
	Client *http.Client
}

// 获取 http 客户端
func DefaultHTTPClient() Client {
	// 单例模式
	defaultClientOnce.Do(func() {
		httpClient := new(httpClient)
		httpClient.Client = defaultClient
		clientSingleton = httpClient
	})

	return clientSingleton
}

// 发起http请求
func (client *httpClient) Do(configFuncArr ...ConfigRequestFunc) (rsp *Response, err error) {
	request := buildRequest()
	defer request.cancel()
	rsp = new(Response)
	for _, currentConfigFunc := range configFuncArr {
		currentConfigFunc(request)
	}
	request.buildRequestBody()
	err = request.checkIsValid()
	if nil != err {
		return nil, err
	}

	err = client.commonSendRequest(request, rsp)
	return
}

func (client *httpClient) commonSendRequest(request *Request, response *Response) (err error) {
	req, err := http.NewRequestWithContext(request.ctx, request.method, request.url, strings.NewReader(request.body))
	for key, value := range request.header {
		req.Header.Set(key, value)
	}
	if err != nil {
		log.Error("[http.client.Do] request url: %s, method: %s, make request failed: %s",
			request.url, request.method, err.Error())
		return
	}

	if request.method == string(methodGet) {
		query := req.URL.Query()
		for k, v := range request.param {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
	}
	startTime := time.Now()
	rsp, err := defaultClient.Do(req)
	if err != nil {
		log.Error("[http.client.Do] request url: %s, method: %s, do request failed: %s", request.url,
			request.method, err.Error())
		return err
	}

	// check return code
	if rsp.StatusCode != 200 && rsp.StatusCode != 204 {
		log.Error("[http.client.Do] request url: %s, method: %s, return code: %d", request.url,
			request.method, rsp.StatusCode)
		return errorcode.BuildErrorWithMsg(errorcode.NetReturnCode, strconv.Itoa(rsp.StatusCode))
	}

	rspBytes, _ := ioutil.ReadAll(rsp.Body)
	endTime := time.Now()
	log.Info("[http.client.Do] request url: %s success, cost: %d seconds", request.url,
		endTime.Unix()-startTime.Unix())
	response.Body = string(rspBytes)
	response.Header = make(map[string]string)

	for key, valueArr := range rsp.Header {
		response.Header[key] = valueArr[0]
	}

	return nil
}
