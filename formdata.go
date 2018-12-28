package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// FormData 为 multipart/form-data 请求体提供快捷生成工具
type FormData struct {
	buf    *bytes.Buffer
	w      *multipart.Writer
	closed bool
}

// NewFormData 初始化 FormData
func NewFormData() *FormData {
	buf := bytes.NewBuffer([]byte{})
	return &FormData{
		buf:    buf,
		w:      multipart.NewWriter(buf),
		closed: false,
	}
}

// Add 表单字段
func (f *FormData) Add(key string, val interface{}, header ...textproto.MIMEHeader) error {
	// 创建默认 MIMEHeader
	h := make(textproto.MIMEHeader)

	// 设置默认 Content-Disposition
	h.Set(
		"Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(key)),
	)

	if len(header) > 0 {
		h = header[0]

		if disp := h.Get("Content-Disposition"); disp == "" {
			h.Set(
				"Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(key)),
			)
		}
	}

	// 如果是文件
	if file, ok := val.(*os.File); ok {
		disp := h.Get("Content-Disposition")
		if strings.Index(disp, "filename=") == -1 {
			fpath := filepath.ToSlash(file.Name())
			_, fname := filepath.Split(fpath)

			h.Set(
				"Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes(key), escapeQuotes(fname)),
			)
		}

		ctt := h.Get("Content-Type")
		if ctt == "" {
			h.Set("Content-Type", "application/octet-stream")
		}

		w, err := f.w.CreatePart(h)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, file)
		return err
	}

	w, err := f.w.CreatePart(h)
	if err != nil {
		return err
	}

	// 如果是 io.Reader
	r, ok := val.(io.Reader)
	if ok {
		_, err = io.Copy(w, r)
		return err
	}

	// 其他简单类型, 只支持 int, int64, float64, string, []byte, time.Time
	// 不建议直接使用 time.Time 类型, 因为转换的格式被写死为 YYYY-MM-DD HH24:mi:ss
	switch v := val.(type) {
	case int:
		_, err = w.Write([]byte(strconv.Itoa(v)))
		return err
	case int64:
		_, err = w.Write([]byte(strconv.FormatInt(v, 10)))
		return err
	case float64:
		_, err = w.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
		return err
	case string:
		_, err = w.Write([]byte(v))
		return err
	case time.Time:
		ts := v.Format("2016-01-02 15:04:05")
		_, err = w.Write([]byte(ts))
		return err
	case []byte:
		_, err = w.Write(v)
		return err
	default:
		return errors.New("不支持的表单字段类型")
	}
}

// AddFile 添加文件字段
func (f *FormData) AddFile(key, fpath string, header ...textproto.MIMEHeader) error {
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	return f.Add(key, file, header...)
}

// SetBoundary 设置 boundary, 不设置的话, 则使用默认设置
func (f *FormData) SetBoundary(b string) error {
	return f.w.SetBoundary(b)
}

// Boundary 返回 boundary
func (f *FormData) Boundary() string {
	return f.w.Boundary()
}

// Close must be called at the end
func (f *FormData) Close() error {
	if !f.closed {
		f.closed = true
		return f.w.Close()
	}

	return nil
}

// Bytes 获取 bytes 数据
func (f *FormData) Bytes() []byte {
	if !f.closed {
		f.Close()
	}

	bs, _ := ioutil.ReadAll(f.buf)
	return bs
}

// Data 返回 io.Reader 可以直接传入 http 请求体
func (f *FormData) Data() io.Reader {
	if !f.closed {
		f.Close()
	}

	return f.buf
}

// ContentType 获取 content-type
func (f *FormData) ContentType() string {
	return f.w.FormDataContentType()
}
