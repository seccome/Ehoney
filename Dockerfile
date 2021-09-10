FROM centos:7
RUN  sed -i -e '/mirrors.cloud.aliyuncs.com/d' -e '/mirrors.aliyuncs.com/d' /etc/yum.repos.d/CentOS-Base.repo

RUN yum -y install gcc gcc-c++ kernel-devel
RUN mkdir -p /go/src
COPY . /go/src
WORKDIR /go/src
RUN cp -r tool/go /usr/local/go
RUN chmod -R 755 /usr/local/go
RUN export PATH=$PATH:/usr/local/go/bin
ENV PATH /usr/local/go/bin:$PATH
RUN go version

ENV GO111MODULE on
ENV GOPROXY "https://gocenter.io"
ENV APP_NAME "decept-defense"
RUN go build .

WORKDIR /ehoney
VOLUME /ehoney
RUN rpm --rebuilddb && yum -y install kde-l10n-Chinese telnet net-tools && localedef -c -f UTF-8 -i zh_CN zh_CN.utf8

# 设置系统环境变量
ENV LANG zh_CN.UTF-8
ENV LANGUAGE zh_CN:zh
ENV LC_ALL zh_CN.UTF-8

# 创建应用安装目录和ehoney账户
RUN mkdir ehoney/nginx -p && mkdir ehoney/html -p

# 安装nginx
ADD ./front/nginx_rpm/*.rpm /ehoney/nginx/

RUN chmod 755 /ehoney/nginx/*.rpm
RUN yum -y install /ehoney/nginx/*.rpm

ADD ./front/ehoney.conf  /etc/nginx/conf.d/ehoney.conf
ADD ./front/nginx.conf  /etc/nginx/nginx.conf


# 安装前端应用
ADD ./front/decept-defense.tar.gz /ehoney/html/decept-defense


#删除压缩包
RUN rm -rf /ehoney/nginx/*.rpm

COPY dockerStart.sh /ehoney/dockerStart.sh
RUN chmod +x /ehoney/dockerStart.sh

WORKDIR /go/src

## 配置 对外端口
EXPOSE 8080 8082
RUN echo "$CONFIGS"
CMD ["/ehoney/dockerStart.sh", "$CONFIGS"]