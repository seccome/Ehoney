# requests

> golang http请求库，类似python requests。

## 特点

* `GET`、`POST`、`PUT`、`DELETE`（Common HTTP methods）
* `application/json`、`application/x-www-form-urlencoded`、`multipart/form-data`
* 全局请求拦截器

## 例子

### get请求

```go
func getText() {
	text, err := requests.Get("http://127.0.0.1:8080/ping").
		Params(url.Values{
			"param1": {"value1"},
			"param2": {"123"},
		}).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
```

http报文：

```
GET http://127.0.0.1:8080/ping?param1=value1&param2=123 HTTP/1.1
```

### 表单请求

```go
func postForm() {
	text, err := requests.Post("http://127.0.0.1:8080/ping").
		Params(url.Values{
			"param1": {"value1"},
			"param2": {"123"},
		}).
		Form(url.Values{
			"form1": {"value1"},
			"form2": {"123"},
		}).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
```

http报文：

```
POST http://127.0.0.1:8080/ping?param1=value1&param2=123 HTTP/1.1
Content-Type: application/x-www-form-urlencoded

form1=value1&form2=123
```

### json请求

```go
func postJson() {
	text, err := requests.Post("http://127.0.0.1:8080/ping").
		Params(url.Values{
			"param1": {"value1"},
			"param2": {"123"},
		}).
		Json(map[string]interface{}{
			"json1": "value1",
			"json2": 2,
		}).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
```

http报文：

```
POST http://127.0.0.1:8080/ping?param1=value1&param2=123 HTTP/1.1
Content-Type: application/json

{"json1": "value1", "json2": 2}
```

### 文件上传

```go
func postMultipart() {
	text, err := requests.Post("http://127.0.0.1:8080/ping").
		Params(url.Values{
			"param1": {"value1"},
			"param2": {"123"},
		}).
		Multipart(requests.FileForm{
			Value: url.Values{
				"form1": {"value1"},
				"form2": {"value2"},
			},
			File: map[string]string{
				"file1": "./examples/main.go",
				"file2": "./examples/main.go",
			},
		}).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
```

http报文：

```
POST http://127.0.0.1:8080/ping?param1=value1&param2=123 HTTP/1.1
Content-Type: multipart/form-data; boundary=947f4ca12d44786ccda8f8cd60e083fca2ec1ede6d8f1bad69f4cf03bc8a

--947f4ca12d44786ccda8f8cd60e083fca2ec1ede6d8f1bad69f4cf03bc8a
Content-Disposition: form-data; name="file1"; filename="./examples/main.go"
Content-Type: application/octet-stream

bytes...

--947f4ca12d44786ccda8f8cd60e083fca2ec1ede6d8f1bad69f4cf03bc8a
Content-Disposition: form-data; name="file2"; filename="./examples/main.go"
Content-Type: application/octet-stream

bytes...

--947f4ca12d44786ccda8f8cd60e083fca2ec1ede6d8f1bad69f4cf03bc8a
Content-Disposition: form-data; name="form1"

value1
--947f4ca12d44786ccda8f8cd60e083fca2ec1ede6d8f1bad69f4cf03bc8a
Content-Disposition: form-data; name="form2"

value2
--947f4ca12d44786ccda8f8cd60e083fca2ec1ede6d8f1bad69f4cf03bc8a--
```

### 保存文件

```go
func save() {
	err := requests.Get("https://github.com/xuanbo/requests").
		Send().
		Save("./requests.html")
	if err != nil {
		panic(err)
	}
}
```

### 检查响应状态码

```go
func save() {
	err := requests.Get("https://github.com/xuanbo/requests").
		Send().
		// resp status code must be 200.
		StatusOk().
		Save("./requests.html")
	if err != nil {
		panic(err)
	}
}
```

* `StatusOk()`
* `Status2xx()`

### 自定义http

```go
func customHttp() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	text, err := requests.Request("https://github.com/xuanbo", "OPTIONS", client).
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
```

### 全局请求拦截器

```go
func requestInterceptor() {
	logRequestInterceptor := func(request *http.Request) error {
		fmt.Println(request.URL)
		return nil
	}
	requests.AddRequestInterceptors(logRequestInterceptor)

	text, err := requests.Get("https://github.com/xuanbo").
		Send().
		Text()
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
```