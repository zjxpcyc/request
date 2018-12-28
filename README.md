# Request
简单版本的 http request . 主要用途是将日常的 http 接口请求抽离出来

## Config
请求内容统一放到一个模块中
```golang
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
```

同时也封装了部分方法
```golang
// NewConfig 新建 config 实例
func NewConfig(rawurl string, method ...string) *Config

// SetContentType 添加 content-type
func (c *Config) SetContentType(val string) *Config

// AddHeader 添加 header 头
func (c *Config) AddHeader(key, val string) *Config

// AddParam 添加 query 或者 search 参数
func (c *Config) AddParam(key, val string) *Config

// SetData 设置请求主体
func (c *Config) SetData(dt io.Reader) *Config
```
上述方法均返回 `*Config` 的原因是，可以实现链式调用


## FormData
请求体的部分, 通用的主要有三种
1. application/x-www-form-urlencoded
2. multipart/form-data
3. raw

其中 1 可以通过 `url.Values` 实现, 3 `raw` 的方式不确定, 不容易封装。 `FormData` 是针对方式2的实现

主要方法有:
```golang
// NewFormData 初始化 FormData
// 必须通过此方法进行实例化
func NewFormData() *FormData

// Add 添加表单字段
// key 是表单字段名称, val 是对应值, 支持的类型有
// int, int64, float64, string, []byte, time.Time, *os.File, io.Reader
// header 为自定义的头信息, 可以不传
func (f *FormData) Add(key string, val interface{}, header ...textproto.MIMEHeader) error

// AddFile 添加本地文件 fpath 为本地文件路径
func (f *FormData) AddFile(key, fpath string, header ...textproto.MIMEHeader) error

// SetBoundary 设置 boundary, 不设置的话, 则使用默认设置
// 如果不清楚 boundary 是什么, 那么这个方法以及下面的 Boundary 方法, 你可以当做不存在, 不会影响你的使用
func (f *FormData) SetBoundary(b string) error

// Boundary 返回 boundary
// 注意, 这个仅仅是字符串的 boundary 前面的 - 需要自动
func (f *FormData) Boundary() string

// Close 最好在字段添加之后调用一次
func (f *FormData) Close() error

// Bytes 获取 bytes 数据, 默认会自动调用 Close 方法
func (f *FormData) Bytes() []byte

// Data 返回 io.Reader 可以直接传入 http 请求体, 默认会自动调用 Close 方法
func (f *FormData) Data() io.Reader

// ContentType 获取 content-type
// 辅助功能, 可以自己组装
func (f *FormData) ContentType() string
```

## Request 方法
封装了几个通用的调用方法 `Do`, `Get`, `Post`, `Put`, `Delete`
```golang
// Do 通用方法, 可以使用任意 http method
// 第一个返回值 io.ReadCloser 其实就是 *http.Response 里面的 body
// 因此在调用的时候，需要进行 close 操作
func Do(config *Config) (io.ReadCloser, *http.Response, error)

// Get 请求, 其实就是调用的 Do, 只是将 method = http.MethodGet
func Get(config *Config) (io.ReadCloser, *http.Response, error)

// Post
func Post(config *Config) (io.ReadCloser, *http.Response, error)

// Put
func Put(config *Config) (io.ReadCloser, *http.Response, error)

// Delete
func Delete(config *Config) (io.ReadCloser, *http.Response, error)
```

## Demo
如果是一般的请求, 比如请求的参数或者表单都只是简单的字段，没有文件发送的情况
```golang
// 1. 先确认 config
config := request.NewConfig("http://baidu.com")
// 2. 添加一些简单的请求字段放在 query 里面
conf.AddParam("q", "小说将夜")
// 3. 确认请求方式, 如果没有, 默认是 Get
conf.Method = http.MethodGet
// 4. 进行远程请求
data, response, err := request.Do(conf)
// TODO
```


如果是带有文件传输的, 可以使用 `FormData` 工具
```golang
// 1. 初始化一个实例
formData := request.NewFormData()

// 2. 加入文件, AddFile 的第二个参数是文件的路径
// 如果你已经有一个 file handle, 可以直接使用 Add 方法
err := formData.AddFile("file", "./README.md")
if err != nil {
	// TODO
}

// 3. 请求内容组合
config := request.NewConfig("http://baidu.com")
// formData.Data() 就是一个 io.Reader 因此可以直接传值
config.SetData(formData.Data())
// 设置 content-type
config.SetContentType(formData.ContentType())

data, response, err := request.Post(config)
// TODO
```
