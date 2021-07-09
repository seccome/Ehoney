## 系统设置

> 系统设置可进行镜像源配置、查看镜像列表、进行协议配以及密签跟踪服务器设置



**镜像源采用harbor为镜像源、版本为2.0、指定harbor的URL、用户名、密码以及项目名**

![系统配置-镜像源设置](../img/系统配置-镜像源设置.png)
注意：如果设置自定义harbor 需要修改 /usr/lib/systemd/system/docker.service 文件 
设置 ExecStart=/usr/bin/dockerd --insecure-registry=47.96.71.197:90" 中的 --insecure-registry的值为harbor地址。
并执行 1、sudo systemctl daemon-reload  2、sudo systemctl restart docker

**镜像列表展示当前可用来创建蜜罐的镜像列表**

!!! 注意： 这里可以对镜像的端口进行设置、当前端口是确定的、除非你明确修改的意义、否则请不要随意修改

![系统设置-镜像列表](../img/系统设置-镜像列表.png)



**进行协议转发的服务**

![系统设置-协议配置](../img/系统设置-协议配置.png)



**密签跟踪服务器设置**

![系统设置-跟踪服务器](../img/系统设置-跟踪服务器.png)
