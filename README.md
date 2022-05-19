<p align="center">
  <img width="200px" src="./doc/img/logo_ehoney_black.svg" alt="Ehoney" />
  <br/>
  <br/>
  <br> 中文 | <a href="README-EN.md">English</a>
  <h1 align="center">欢迎来到Ehoney 👋</h1>
  <br/>
  <p align="center">
  <img align="center" src="https://img.shields.io/badge/release-v1.0.0-green" />
  <img align="center" src="https://img.shields.io/badge/language-goland-orange" />
  <img align="center" src="https://img.shields.io/badge/documentation-yes-ff69b4" />
  <img align="center" src="https://img.shields.io/badge/license-Apache%202-blue" />
  </p>
</p>


> ⭐️ Seccome Teamer积累十几年的安全经验，都将对外逐步开放，首开的Ehoney欺骗防御系统，该系统是基于云原生的欺骗防御系统，也是业界唯一开源的对标商业系统的产品，欺骗防御系统通过部署高交互高仿真蜜罐及流量代理转发，再结合自研密签及诱饵，将攻击者攻击引导到蜜罐中达到扰乱引导以及延迟攻击的效果，可以很大程度上保护业务的安全。`护网必备良药`，该平台只提供安全技术防护能力，任何人不得用于任何不法行为⭐️   


![介绍视频](./doc/img/介绍.gif)


🏠 [使用文档](https://seccome.github.io/Ehoney/) &nbsp;&nbsp; :triangular_flag_on_post: [演示环境](http://47.98.206.178:8080/decept-defense)   

## 📝 特点

- **支持丰富的蜜罐类型**

1. **通用蜜罐**： SSH 蜜罐、Http蜜罐、Redis蜜罐、Telnet蜜罐、Mysql蜜罐、RDP 蜜罐
2. **IOT蜜罐**：  RTSP 蜜罐
3. **工控蜜罐**： ModBus 蜜罐

- **基于云原生技术**<br>
基于k3s打造saas平台欺骗防御，无限生成蜜罐，真实仿真业务环境

- **业内独一无二密签技术**<br>
独创的密签技术，支持20多种密签，如文件、图片，邮件等

- **强大诱饵**<br>
支持数十种诱饵，通过探针管理，进行欺骗引流

- **可视化拓扑**<br>
可以可视化展示攻击视图，让所有攻击可视化，形成完整的攻击链路

- **动态对抗技术**<br>
基于LSTM的预测算法，可以预测黑客下一步攻击手段，动态欺骗，延缓黑客攻击时间，保护真实业务

- **强大的定制化**<br>
支持自定义密签、诱饵、蜜罐等，插件化安装部署，满足一切特性需求

## ⛴ 环境准备

- **系统要求**: CentOS 7 及以上
- **最低配置：**: 内存4G、磁盘空间10G以上
- **建议配置：**: 内存8G、磁盘空间30G以上

<br>


## 🔧 快速部署

```shell
git clone https://github.com/seccome/Ehoney.git
cd Ehoney && chmod +x quick-start.sh && ./quick-start.sh

# 此安装过程会比较耗时、耐心等待

**all the services are ready and happy to use!!!**
# 代表安装成功。
```

访问 `http://IP:8080/decept-defense` 进入系统登录页

默认账户
       <font color=Blue>用户名: `admin`</font>
       <font color=Blue>密码: `123456`</font>

<br>

## 🖥️ 使用演示

![操作视频](./doc/img/操作视频.gif)

<br>

## 🚀 效果展示

- **攻击大屏**

![攻击事件大屏](./doc/img/攻击事件大屏.png)

- **蜜罐拓扑**

![蜜罐拓扑图](./doc/img/蜜罐拓扑图.png)

- **告警列表**

![告警列表](./doc/img/告警列表.png)


## 🙏 讨论区 

如有问题可以在 GitHub 提 issue, 也可在下方的讨论组里，问题我们都会及时处理

1. GitHub issue: [创建issue](https://github.com/seccome/Ehoney/issues/new )
2. Ehoney 技术交流群: 679424748
3. 邮箱: sihan@tsign.cn
