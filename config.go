package request

import (
	"io"
	"net/http"
	"net/url"
)

// Config 请求配置
type Config struct {
	// header 头
	Headers http.Header

	// 请求地址
	URL string

	// 请求方法
	Method string

	// 请求 query 或者 search 参数
	Params url.Values

	// 请求主体
	Data io.Reader
}

// NewConfig 新建 config 实例
func NewConfig(rawurl string, method ...string) *Config {
	mth := http.MethodGet
	if method != nil && len(method) > 0 {
		mth = method[0]
	}

	return &Config{
		URL:     rawurl,
		Method:  mth,
		Headers: make(http.Header),
		Params:  make(url.Values),
	}
}

// SetContentType 添加 content-type
func (c *Config) SetContentType(val string) *Config {
	c.Headers.Add("Content-Type", val)
	return c
}

// AddHeader 添加 header 头
func (c *Config) AddHeader(key, val string) *Config {
	c.Headers.Add(key, val)
	return c
}

// AddParam 添加 query 或者 search 参数
func (c *Config) AddParam(key, val string) *Config {
	c.Params.Add(key, val)
	return c
}

// SetData 设置请求主体
func (c *Config) SetData(dt io.Reader) *Config {
	c.Data = dt
	return c
}
