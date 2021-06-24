<p align="center">
  <img width="200px" src="./doc/img/logo_ehoney_black.svg" alt="Ehoney" />
  <br/>
  <br/>
  <h1 align="center">Ehoney 欺骗防御系统</h1>
  <br/>
</p>


> ⭐️ e签宝安全团队积累十几年的安全经验，都将对外逐步开放，首开的Ehoney欺骗防御系统，该系统是基于云原生的欺骗防御系统，也是业界唯一开源的对标商业系统的产品，欺骗防御系统通过部署高交互高仿真蜜罐及流量代理转发，再结合自研密签及诱饵，将攻击者攻击引导到蜜罐中达到扰乱引导以及延迟攻击的效果，可以很大程度上保护业务的安全。⭐️   

![介绍视频](./doc/img/介绍.gif)



## 📝 特点

- **支持丰富的蜜罐类型**

1. **通用蜜罐**： SSH 蜜罐、Http蜜罐、Redis蜜罐、Telnet蜜罐、Mysql蜜罐、RDP 蜜罐
2. **IOT蜜罐**：  RTSP 蜜罐
3. **工控蜜罐**： ModBus 蜜罐

- **基于云原生技术**
基于k3s打造saas平台欺骗防御，无限生成蜜罐，真实仿真业务环境

- **业内独一无二密签技术**
独创的密签技术，支持20多种密签，如文件、图片，邮件等

- **强大诱饵**
支持数十种诱饵，通过探针管理，进行欺骗引流

- **可视化拓扑**
可以可视化展示攻击视图，让所有攻击可视化，形成完整的攻击链路

- **动态对抗技术**
基于LSTM的预测算法，可以预测黑客下一步攻击手段，动态欺骗，延缓黑客攻击时间，保护真实业务

- **强大的定制化**
支持自定义密签、诱饵、蜜罐等，插件化安装部署，满足一切特性需求

## ⛴ 环境准备

- **系统要求**: CentOS 7
- **配置要求**: 内存4G、磁盘空间10G以上

> Note: 以上的配置要求为系统运行的最优配置

<br>

## 🔧 快速部署

```shell
git clone https://github.com/seccome/Ehoney.git
cd Ehoney && chmod +x quick-start.sh && ./quick-start.sh
```

此安装过程会比较耗时、耐心等待后看到如下提示语

​			**all the services are ready and happy to use!!!**

代表安装成功。

访问http://IP:8080/decept-defense 进入系统登录页
默认账户密码
       <font color=Blue>用户名: admin</font>
       <font color=Blue>密码: 123456</font>

<br>

## 🖥️ 使用演示

![操作视频](./doc/img/操作视频.gif)

系统的详细使用文档请参考[这里](https://www.showdoc.com.cn/1432924569255366?page_id=7002138596779961)

<br>

## 🚀 效果展示

- **攻击大屏**

![攻击事件大屏](./doc/img/攻击事件大屏.png)

- **蜜罐拓扑**

![蜜罐拓扑图](./doc/img/蜜罐拓扑图.png)

- **告警列表**

![告警列表](./doc/img/告警列表.png)



## 🙏 关于 

如果大家对系统有好的建议或想法, 请[创建issue](https://github.com/seccome/Ehoney/issues/new ), 我们会及时回复并处理。

