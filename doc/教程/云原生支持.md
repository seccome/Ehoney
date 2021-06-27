## 概念

## 使用说明

### Ehoney对云原生的支持主要是通过K3s以及Falco实现的
* 蜜罐作为Pod容器部署在K3s上、而Falco就是为云原生容器安全而生，可以通过灵活的规则引擎来描述任何类型的主机或者容器的行为和活动。

* Ehoney 通过自定义Falco检测规则实现对K3s上的任意容器进行监测，发现可疑的文件操作都将进行上报, 并配合服务端对容器内部的文件拷贝操作能够获取到攻击者的攻击脚本样例和操作过程等信息。

### Falco 是什么？

* Falco 是云原生的容器运行时（Runtime）安全项目，主要有 Sysdig 为主开发，为 CNCF 项目

* Falco 可以轻松使用内核事件，并使用 Kubernetes 和其他云本机中的信息补充和丰富事件。Falco 具有一组专门为 Kubernetes，Linux 和云原生构建的安全规则。如果系统中违反了规则，Falco 将发送警报，通知到用户


![](https://www.showdoc.com.cn/server/api/attachment/visitfile/sign/44963f789b39c221893f03bd21c5d807)

## 实现设计
