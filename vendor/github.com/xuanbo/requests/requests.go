package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

const (
	ContentType               = "Content-Type"
	ApplicationJSON           = "application/json"
	ApplicationFormUrlencoded = "application/x-www-form-urlencoded"
)

// RequestInterceptor 请求拦截器
// 返回不为nil，即有错误会终止后续执行
type RequestInterceptor func(request *http.Request) error

// requestInterceptorChain 请求拦截链
type requestInterceptorChain struct {
	mutex        *sync.RWMutex
	interceptors []RequestInterceptor
}

// defaultRequestInterceptorChain 默认的请求拦截链实例
var defaultRequestInterceptorChain = &requestInterceptorChain{
	mutex:        new(sync.RWMutex),
	interceptors: make([]RequestInterceptor, 0),
}

// Client 封装了http的参数等信息
type Client struct {
	// 自定义Client
	client *http.Client

	url    string
	method string
	header http.Header
	params url.Values

	form      url.Values
	json      interface{}
	multipart FileForm
}

// FileForm form参数和文件参数
type FileForm struct {
	Value url.Values
	File  map[string]string
}

// Result http响应结果
type Result struct {
	Resp *http.Response
	Err  error
}

// Get http `GET` 请求
func Get(url string) *Client {
	return newClient(url, http.MethodGet, nil)
}

// Post http `POST` 请求
func Post(url string) *Client {
	return newClient(url, http.MethodPost, nil)
}

// Put http `PUT` 请求
func Put(url string) *Client {
	return newClient(url, http.MethodPut, nil)
}

// Delete http `DELETE` 请求
func Delete(url string) *Client {
	return newClient(url, http.MethodDelete, nil)
}

// Request 用于自定义请求方式，比如`HEAD`、`PATCH`、`OPTIONS`、`TRACE`
// client参数用于替换DefaultClient，如果为nil则会使用默认的
func Request(url, method string, client *http.Client) *Client {
	return newClient(url, method, client)
}

// Params http请求中url参数
func (c *Client) Params(params url.Values) *Client {
	for k, v := range params {
		c.params[k] = v
	}
	return c
}

// Header http请求头
func (c *Client) Header(k, v string) *Client {
	c.header.Set(k, v)
	return c
}

// Headers http请求头
func (c *Client) Headers(header http.Header) *Client {
	for k, v := range header {
		c.header[k] = v
	}
	return c
}

// Form 表单提交参数
func (c *Client) Form(form url.Values) *Client {
	c.header.Set(ContentType, ApplicationFormUrlencoded)
	c.form = form
	return c
}

// Json json提交参数
// 如果是string，则默认当作是json字符串；否则会序列化为json字节数组，再发送
func (c *Client) Json(json interface{}) *Client {
	c.header.Set(ContentType, ApplicationJSON)
	c.json = json
	return c
}

// Multipart form-data提交参数
func (c *Client) Multipart(multipart FileForm) *Client {
	c.multipart = multipart
	return c
}

// Send 发送http请求
func (c *Client) Send() *Result {
	var result *Result

	// 处理query string
	if c.params != nil && len(c.params) != 0 {
		// 如果url中已经有query string参数，则只需要&拼接剩下的即可
		encoded := c.params.Encode()
		if strings.Index(c.url, "?") == -1 {
			c.url += "?" + encoded
		} else {
			c.url += "&" + encoded
		}
	}

	// 根据不同的Content-Type设置不同的http body
	contentType := c.header.Get(ContentType)
	if c.multipart.Value != nil || c.multipart.File != nil {
		result = c.createMultipartForm()
	} else if strings.HasPrefix(contentType, ApplicationJSON) {
		result = c.createJson()
	} else if strings.HasPrefix(contentType, ApplicationFormUrlencoded) {
		result = c.createForm()
	} else {
		// 不是以上类型，就不设置http body
		result = c.createEmptyBody()
	}

	return result
}

// createMultipartForm 创建form-data的请求
func (c *Client) createMultipartForm() *Result {
	var result = new(Result)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 设置文件字节
	for name, filename := range c.multipart.File {
		file, err := os.Open(filename)
		if err != nil {
			result.Err = err
			return result
		}

		part, err := writer.CreateFormFile(name, filename)
		if err != nil {
			result.Err = err
			return result
		}

		// todo 这里的io.Copy实现，会把file文件都读取到内存里面，然后当做一个buffer传给NewRequest。对于大文件来说会占用很多内存
		_, err = io.Copy(part, file)
		if err != nil {
			result.Err = err
			return result
		}

		err = file.Close()
		if err != nil {
			result.Err = err
			return result
		}
	}

	// 设置field
	for name, values := range c.multipart.Value {
		for _, value := range values {
			_ = writer.WriteField(name, value)
		}
	}

	err := writer.Close()
	if err != nil {
		result.Err = err
		return result
	}

	req, err := http.NewRequest(c.method, c.url, body)
	req.Header = c.header
	req.Header.Set(ContentType, writer.FormDataContentType())
	c.doSend(req, result)
	return result
}

// createForm 创建application/json请求
func (c *Client) createJson() *Result {
	var result = new(Result)

	b, err := json.Marshal(c.json)
	if err != nil {
		result.Err = err
		return result
	}

	req, err := http.NewRequest(c.method, c.url, bytes.NewReader(b))
	if err != nil {
		result.Err = err
		return result
	}

	req.Header = c.header
	c.doSend(req, result)
	return result
}

// createForm 创建application/x-www-form-urlencoded请求
func (c *Client) createForm() *Result {
	var result = new(Result)

	form := c.form.Encode()

	req, err := http.NewRequest(c.method, c.url, strings.NewReader(form))
	if err != nil {
		result.Err = err
		return result
	}

	req.Header = c.header
	c.doSend(req, result)
	return result
}

// createEmptyBody 没有内容的body
func (c *Client) createEmptyBody() *Result {
	var result = new(Result)

	req, err := http.NewRequest(c.method, c.url, nil)
	if err != nil {
		result.Err = err
		return result
	}

	req.Header = c.header
	c.doSend(req, result)
	return result
}

// doSend 发送请求
func (c *Client) doSend(req *http.Request, result *Result) {
	// 调用拦截器，遇到错误就退出
	if err := c.beforeSend(req); err != nil {
		result.Err = err
		return
	}

	// 发送请求
	result.Resp, result.Err = c.client.Do(req)
}

// beforeSend 发送请求前，调用拦截器
func (c *Client) beforeSend(req *http.Request) error {
	mutex := defaultRequestInterceptorChain.mutex
	mutex.RLock()
	defer mutex.RUnlock()

	// 遍历调用拦截器
	for _, interceptor := range defaultRequestInterceptorChain.interceptors {
		err := interceptor(req)
		if err != nil {
			return err
		}
	}
	return nil
}

// StatusOk 判断http响应码是否为200
func (r *Result) StatusOk() *Result {
	if r.Err != nil {
		return r
	}
	if r.Resp.StatusCode != http.StatusOK {
		r.Err = errors.New("status code is not 200")
		return r
	}

	return r
}

// Status2xx 判断http响应码是否为2xx
func (r *Result) Status2xx() *Result {
	if r.Err != nil {
		return r
	}
	if r.Resp.StatusCode < http.StatusOK || r.Resp.StatusCode >= http.StatusMultipleChoices {
		r.Err = errors.New("status code is not match [200, 300)")
		return r
	}

	return r
}

// Raw 获取http响应内容，返回字节数组
func (r *Result) Raw() ([]byte, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	b, err := ioutil.ReadAll(r.Resp.Body)
	if err != nil {
		r.Err = err
		return nil, r.Err
	}
	defer r.Resp.Body.Close()

	return b, r.Err
}

// Text 获取http响应内容，返回字符串
func (r *Result) Text() (string, error) {
	b, err := r.Raw()
	if err != nil {
		r.Err = err
		return "", r.Err
	}

	return string(b), nil
}

// Json 获取http响应内容，返回json
func (r *Result) Json(v interface{}) error {
	b, err := r.Raw()
	if err != nil {
		r.Err = err
		return r.Err
	}

	return json.Unmarshal(b, v)
}

// Save 获取http响应内容，保存为文件
func (r *Result) Save(name string) error {
	if r.Err != nil {
		return r.Err
	}

	f, err := os.Create(name)
	if err != nil {
		r.Err = err
		return r.Err
	}
	defer f.Close()

	_, err = io.Copy(f, r.Resp.Body)
	if err != nil {
		r.Err = err
		return r.Err
	}
	defer r.Resp.Body.Close()

	return nil
}

// newClient 创建Client
func newClient(u string, method string, client *http.Client) *Client {
	// client为nil则使用默认的DefaultClient
	if client == nil {
		client = http.DefaultClient
	}
	return &Client{
		client: client,
		url:    u,
		method: method,
		header: make(http.Header),
		params: make(url.Values),
		form:   make(url.Values),
	}
}

// AddRequestInterceptors 添加请求拦截器
func AddRequestInterceptors(interceptors ...RequestInterceptor) {
	mutex := defaultRequestInterceptorChain.mutex
	mutex.Lock()
	defer mutex.Unlock()

	// 添加到拦截器链
	for _, interceptor := range interceptors {
		defaultRequestInterceptorChain.interceptors = append(defaultRequestInterceptorChain.interceptors, interceptor)
	}
}
