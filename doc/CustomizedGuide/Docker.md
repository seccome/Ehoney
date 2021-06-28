## 功能描述
docker镜像是用于k3s创建蜜罐时，需要的容器镜像。

## 开发说明
* 可以通过自定义Dockerfile，build完之后，导出docker镜像。命令如下：<br>
1、docker build ：`docker build -t 镜像名 .`<br>
2、导出docker镜像 ：`docker save -o /xxx/镜像名.tar 镜像名`<br>
3、将镜像包同步给 ask@seccome.com <br>
4、欺骗防御开发团队对镜像包进行安全扫描，如果无安全风险，团队会将该镜像包上传到公用Harbor上。

## 使用说明
- 在系统设置-镜像列表中可以对线上拉取的镜像（除了系统默认提供的镜像）进行编辑，包括端口号，服务类型等
- 在新建蜜罐的时候，拉取配置好的镜像列表，创建蜜罐

## 注意事项
1、镜像必须支持/bin/sh
2、镜像必须安装wget命令
