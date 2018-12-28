package request_test

import (
	"testing"

	"github.com/zjxpcyc/request"
)

func TestDo(t *testing.T) {
	config := request.NewConfig("http://baidu.com")

	_, res, err := request.Do(config)
	if err != nil {
		t.Fatalf("Test request.Do error: %v", err)
	}

	if ctt := res.Header.Get("Content-Type"); ctt != "text/html" {
		t.Fatalf("Test request.Do error")
	}
}

func TestPost(t *testing.T) {
	formData := request.NewFormData()
	err := formData.AddFile("file", "./README.md")
	if err != nil {
		t.Fatalf("Test request.Post error: %v", err)
	}

	config := request.NewConfig("http://baidu.com").
		SetData(formData.Data()).
		SetContentType(formData.ContentType())

	_, res, err := request.Post(config)
	if err != nil {
		t.Fatalf("Test request.Do error: %v", err)
	}

	if ctt := res.Header.Get("Content-Type"); ctt != "text/html" {
		t.Fatalf("Test request.Do error")
	}
}
