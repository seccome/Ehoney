## FAQ

### 一 安装问题

1. 服务器环境问题

```
1. 仅支持centos7
2. 磁盘空间必须大于10GB
```

2. 安装问题

```
1. 如果出去端口占用优先确定以及关闭清除之前部署的相关docker容器、然后使用lsof命令查看是否由其他进程占用
2. 镜像拉取失败原因可能由于网络和网速的问题、请确保网络配置正确
3. 启动quick-start.sh 如果出现编码问题，请安装dos2unix进行编码后再执行
```



3. 探针安装

```
1. 如果探针启动后后端探针列表查看不到，需要确定探针目录下的conf/agent.json的Redis配置是否正确
```

### 二 使用问题

1. 创建蜜罐问题

```
1. 首次安装会较慢，原因是需要从harbor服务器上拉取容器镜像。
2. 镜像拉取失败，可能原因为修改了harbor地址，但是没有在docker中修改，需要修改 /usr/lib/systemd/system/docker.service，配置ExecStart=/usr/bin/dockerd --insecure-registry=（harbor的ip）
```

2. 透明代理问题

```
1. 透明代理下发后创建失败 请查看探针目录下的proxy/log/proxy.log 查看日志确定是否是用完端口被占用导致
```

3. 伪装代理问题

```
1. 伪装代理无法下发成功
   1.1 登录蜜网服务器，进入/home/ehoney_proxy/目录，查看代理模块是否具有执行权限。
   1.2 请查看项目部署机器的/home/relay/agent/proxy/log/proxy.log 确实是否因为端口占用或端口过大导致。
   1.3  请查看项目部署机器的/home/relay/agent/conf/agent.json的Redis配置是否正确
```

4. 探针问题

```
1.暂无
```

5. 密签问题

```
1. 密签无法下发到蜜罐，原因可能是该蜜罐镜像不支持/bin/sh或者无法执行wget命令
2. 密签文件打开后无法跟踪，原因可能是密签追踪url配置有问题，可以访问系统设置-密签配置，进行确认。
3. 确认密签追踪url配置方法：访问该url，加上/api/health，如果返回是200，说明部署成功。
```

6. 诱饵问题

```
1. 诱饵无法下发到蜜罐，原因可能是该蜜罐镜像不支持/bin/sh或者无法执行wget命令
```

7. 攻击溯源问题


```
1.暂无
```


### 三 其他问题

1. 一键安装成功，但是浏览器无法访问？

* 检查docker容器状态是否正常

* 检查web服务器容器日志
  ```docker logs -f $(docker ps  | grep decept-defense:latest | awk '{print $1}')```

2. 为什么模拟攻击没有攻击日志？

* 通过透明代理列表，网络探测检查攻击IP、端口可以访问。

* 确保密网服务器上/home/ehoney_proxy目录下对应的代理文件存在。

3. 为什么我部署了多台探针Agent，但是探针列表只显示一个？

* 删除agent/conf目录下agent文件，重新启动agent进程，确保每台服务器上agentID不一致。

4. 为什么蜜罐列表中新建蜜罐失败？

* 检查密网服务器资源状态，有可能是磁盘、cpu利用率满了。一般一台2核4G服务器上最多支持部署20个蜜罐。

* 如果访问harbor不是https协议，需要在docker启动的时候指定insecure-registry。可以通过```kubectl describe pod```查看具体报错信息。

5. 为什么蜜罐拓扑图不显示蜜罐和相关连线？

* 蜜罐拓扑图中是以协议代理、透明代理为维度，如果没有代理则不会显示相关蜜罐。

6. 为什么我部署了探针但是探针列表却不显示？

* 检查agent中配置是否正确。主要检查conf/agent.json中strategyAddr(redis地址)、strategyPass(redis密码)、sshKeyUploadUrl(服务端更新sshkey地址)。

* 确认探针服务器和web/蜜网服务器网络可以通信。

7. 为什么协议代理创建失败？

* 查看协议代理日志/home/relay/agent/proxy/log/proxy.log

* 检查转发端口是否被占用/ssh_proxy文件是否正确
