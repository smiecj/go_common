package http

import (
	"testing"

	"github.com/smiecj/go_common/util/log"
)

func TestMakeGetRequest(t *testing.T) {
	client := DefaultHTTPClient()
	rsp, err := client.Do(Get(), Url("https://prometheus.io/docs/introduction/overview/"))
	if nil != err {
		log.Error("[TestMakeGetRequest] request failed: %s", err.Error())
		t.FailNow()
	}

	log.Info("[TestMakeGetRequest] request success: %s", rsp.Body)
}

/* func TestMakePostUrlEncodeRequest(t *testing.T) {
	client := DefaultHTTPClient()
	rsp, err := client.Do(PostWithUrlEncode(), Url("http://azkaban_address/userReq/doLogin"),
		AddParam("loginEmail", "admin"), AddParam("password", "admin"))
	if nil != err {
		log.Error("[TestMakePostUrlEncodeRequest] request failed: %s", err.Error())
		t.FailNow()
	}

	log.Info("[TestMakePostUrlEncodeRequest] request success: body: %s", rsp.Body)
	for key, value := range rsp.Header {
		log.Info("[TestMakePostUrlEncodeRequest] rsp header: %s -> %s", key, value)
	}
} */
