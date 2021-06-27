##概念
蜜网是由单个或者多个部署了K3S的服务器组成的，蜜网中的服务器支持通过管理平台部署蜜罐，可以仿真企业内网的真实环境。
##使用说明
####创建蜜罐
![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/bb2930f42f0abf97d93a6dcd20d8e41c)
>描述xxx

![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/cf38647efe9b75b83b493209dd0bf224)

>描述xxxx

###创建协议代理
![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/bae7dc0400132636ddae652e24add9b5)
>描述xxxxx

![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/ecc6dfa08a2ba6ceda56a0393c33a078)
>描述

####部署探针
![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/fe8039241bf9e5fa6effe9118000918d)
>描述

###透明转发
![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/4f2f44a33bde0ae732d87dc98f489908)
>描述

###模拟黑客链路

##实现设计
* 通过K3S的云原生技术，部署蜜网服务器。
* 支持只部署Master节点的单机模式，也支持部署多个Worker节点的集群模式。可以根据实际场景进行部署。
* 在蜜网内的服务器上部署agent，可以通过协议代理技术，将透明代理转发过来的流量，根据协议类型代理到蜜罐中去，使黑客最终攻击的是蜜罐环境。
