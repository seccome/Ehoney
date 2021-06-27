## 概念
蜜网是由单个或者多个部署了K3S的服务器组成的，蜜网中的服务器支持通过管理平台部署蜜罐，可以仿真企业内网的真实环境。
## 使用说明
#### 创建蜜罐
![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/bb2930f42f0abf97d93a6dcd20d8e41c)
>输入蜜罐名称和镜像地址，新建蜜罐

![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/cf38647efe9b75b83b493209dd0bf224)


### 创建协议代理
![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/bae7dc0400132636ddae652e24add9b5)
>选择服务类型，选择对应服务的蜜罐，输入需要转发的端口（建议1025-65535范围内的端口号）

![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/ecc6dfa08a2ba6ceda56a0393c33a078)


#### 部署探针
![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/fe8039241bf9e5fa6effe9118000918d)
>通过下载支持，下载agent，在业务服务器上部署探针后，探针会自动注册上线

### 透明转发
![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/4f2f44a33bde0ae732d87dc98f489908)
>通过蜜罐、协议代理转发端口选择已经创建，且在线的协议代理信息，选择已经上线的探针，输入需要转发的端口创建透明代理。

### 模拟黑客链路

## 实现设计
* 通过K3S的云原生技术，部署蜜网服务器。
* 支持只部署Master节点的单机模式，也支持部署多个Worker节点的集群模式。可以根据实际场景进行部署。
* 在蜜网内的服务器上部署agent，可以通过协议代理技术，将透明代理转发过来的流量，根据协议类型代理到蜜罐中去，使黑客最终攻击的是蜜罐环境。
