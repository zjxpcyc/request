package request

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

// Do 进行 http 请求
func Do(config *Config) (io.ReadCloser, *http.Response, error) {
	if config == nil || config.URL == "" {
		return nil, nil, errors.New("请求地址或者内容不能为空")
	}

	if config.Method == "" {
		config.Method = http.MethodGet
	}

	// get 请求不允许有 请求体
	if config.Method == http.MethodGet {
		config.Data = nil
	}

	// 请求地址
	remoteAddr, err := getRemoteAddr(config.URL, config.Params)
	if err != nil {
		return nil, nil, err
	}

	// http 请求客户端
	cli := new(http.Client)
	req, err := http.NewRequest(config.Method, remoteAddr, config.Data)
	if err != nil {
		return nil, nil, err
	}

	// 合并 header
	req.Header = mergeHeaders(req.Header, config.Headers)

	// 发送请求
	res, err := cli.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return res.Body, res, nil
}

// Get 请求
func Get(config *Config) (io.ReadCloser, *http.Response, error) {
	config.Method = http.MethodGet
	return Do(config)
}

// Put 请求
func Put(config *Config) (io.ReadCloser, *http.Response, error) {
	config.Method = http.MethodPut
	return Do(config)
}

// Post 请求
func Post(config *Config) (io.ReadCloser, *http.Response, error) {
	config.Method = http.MethodPost
	return Do(config)
}

// Delete 请求
func Delete(config *Config) (io.ReadCloser, *http.Response, error) {
	config.Method = http.MethodDelete
	return Do(config)
}

func getRemoteAddr(rawurl string, params url.Values) (string, error) {
	if params == nil {
		return rawurl, nil
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	searchStr := params.Encode()
	if u.RawQuery == "" {
		u.RawQuery = searchStr
	} else {
		u.RawQuery += "&" + searchStr
	}

	return u.String(), nil
}

// 合并的规则是, h1 的 k-v 替换掉 h2 对应的内容
func mergeHeaders(h1, h2 http.Header) http.Header {
	t := http.Header{}

	for k, v := range h1 {
		t[k] = v
	}

	for k, v := range h2 {
		t[k] = v
	}

	return t
}
