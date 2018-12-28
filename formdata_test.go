package request_test

import (
	"net/textproto"
	"strings"
	"testing"

	"github.com/zjxpcyc/request"
)

func TestFormData(t *testing.T) {
	fdt := request.NewFormData()

	err := fdt.SetBoundary("yansenisahero")
	if err != nil {
		t.Fatalf("Test request.FormData fail: %v", err)
	}

	fdt.Add("name", "yansen", textproto.MIMEHeader{"Content-Type": []string{"text/plain"}})
	fdt.Add("sex", 1)

	// fdt.Close()

	dt := fdt.Bytes()
	expected := []string{
		"--yansenisahero",
		`Content-Disposition: form-data; name="name"`,
		"Content-Type: text/plain",
		"",
		"yansen",
		"--yansenisahero",
		`Content-Disposition: form-data; name="sex"`,
		"",
		"1",
		"--yansenisahero--",
		"",
	}

	if string(dt) != strings.Join(expected, "\r\n") {
		t.Fatalf("Test request.FormData fail")
	}
}
