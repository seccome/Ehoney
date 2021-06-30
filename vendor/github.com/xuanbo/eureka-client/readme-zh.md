# Eureka-Client

> golang版eureka客户端

![ui](./doc/eureka-server.jpg)

## 特点

* 心跳
* 刷新服务列表（仅仅支持全量拉取）

## 未完成

* ~~当eureka服务端重启后，心跳404时支持重新注册~~

* 以delta方式刷新服务列表（增量拉取）

如果delta被禁用或者首次刷新，则使用全量拉取

## 例子

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	eureka "github.com/xuanbo/eureka-client"
)

func main() {
	// 创建客户端
	client := eureka.NewClient(&eureka.Config{
		DefaultZone:           "http://127.0.0.1:8080/eureka/",
		App:                   "go-example",
		Port:                  10000,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
	})
	// 启动客户端，同步组册服务，异步拉取服务列表、心跳，监听退出信号删除注册信息
	client.Start()

	// http server
	http.HandleFunc("/services", func(writer http.ResponseWriter, request *http.Request) {
		// 获取所有的服务列表
		apps := client.Applications

		b, _ := json.Marshal(apps)
		_, _ = writer.Write(b)
	})

	// 启动http服务
	if err := http.ListenAndServe(":10000", nil); err != nil {
		fmt.Println(err)
	}
}
```

[例子](./examples/main.go)

## 测试

我使用的是Java`spring-cloud-starter-netflix-eureka-server`.

```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-starter-netflix-eureka-server</artifactId>
    <version>2.1.0.RELEASE</version>
</dependency>
```

代码如下:

[spring-cloud-v2](https://github.com/xuanbo/spring-cloud-v2)