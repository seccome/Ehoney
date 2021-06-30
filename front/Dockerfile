# centos:7
FROM 47.96.71.197:90/ehoney-images/gopy:v1
ENV TZ=Asia/Shanghai

#RUN apt-get update --fix-missing && apt-get install -y python-pip --fix-missing
RUN pip install requests
RUN pip install pypdf2
RUN pip install openpyxl==2.6.4
RUN pip install click

ADD . /go/src
WORKDIR /go/src
COPY . /go/src


ENV GO111MODULE=on
ENV APP_NAME="decept-defense"
RUN go build -mod=vendor -o ${APP_NAME}
RUN ls -ltr /go/src
RUN mkdir /go/src/upload/honeysign





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

#CMD ["nginx"]

#删除压缩包
RUN rm -rf /ehoney/nginx/*.rpm
## 配置 对外端口
EXPOSE 8080

ENTRYPOINT ["nginx"]
## CMD  /go/src/$APP_NAME  && nginx &
#
#CMD [ "sh", "-c", "nginx &  && /go/src/$APP_NAME"]

