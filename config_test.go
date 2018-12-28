package request_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/zjxpcyc/request"
)

func TestConfig(t *testing.T) {
	conf := request.NewConfig("http://baidu.com")
	conf.SetContentType("text/plain")
	conf.AddHeader("token", "123-456")
	conf.AddParam("q", "keyword")
	conf.SetData(bytes.NewBuffer([]byte(`yansenisahero`)))

	if conf.URL != "http://baidu.com" {
		t.Fatalf("Test request.Config Fail: init url error")
	}

	if ctt := conf.Headers.Get("Content-Type"); ctt != "text/plain" {
		t.Fatalf("Test request.Config Fail: set content-type error")
	}

	if tk := conf.Headers.Get("token"); tk != "123-456" {
		t.Fatalf("Test request.Config Fail: set header error")
	}

	if q := conf.Params.Get("q"); q != "keyword" {
		t.Fatalf("Test request.Config Fail: set query params error")
	}

	if dt, err := ioutil.ReadAll(conf.Data); err != nil {
		t.Fatalf("Test request.Config Fail: %v", err)
	} else if string(dt) != "yansenisahero" {
		t.Fatalf("Test request.Config Fail: set body data error")
	}
}
