# FROM centos:7
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
# 构建物打包阶段 final stage
#FROM alpine:latest
## 配置 apk包加速镜像为阿里
#RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
## 安装 一些基础包
#RUN apk update \
#    && apk upgrade \
#    && apk add bash \
#    && apk add ca-certificates \
#    && apk add wget \
#    && apk add curl \
#    && apk add libc6-compat \
#    && apk add -U tzdata \
#    && rm -rf /var/cache/apk/*
## 设置 操作系统时区
#RUN rm -rf /etc/localtime \
#    && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime


#ENV APP_NAME="decept-defense"
#ENV APP_ROOT="/apps/"
#RUN mkdir -p $APP_ROOT
#WORKDIR $APP_ROOT
## 从构建阶段复制构建物
#COPY /go/src/${APP_NAME}  $APP_ROOT/
#COPY /go/src/conf  $APP_ROOT/conf
#COPY /go/src/policy  $APP_ROOT/policy
#COPY /go/src/upload  $APP_ROOT/upload
#COPY --from=build-env /go/src/${APP_NAME}  $APP_ROOT/
#COPY --from=build-env /go/src/conf  $APP_ROOT/conf
#COPY --from=build-env /go/src/policy  $APP_ROOT/policy
#COPY --from=build-env /go/src/upload  $APP_ROOT/upload

# 设置启动时预期的命令参数, 可以被 docker run 的参数覆盖掉.
#CMD $APP_ROOT/$APP_NAME


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

COPY dockerStart.sh /ehoney/dockerStart.sh
RUN chmod +x /ehoney/dockerStart.sh

WORKDIR /go/src

## 配置 对外端口
EXPOSE 8080 8082

# ENTRYPOINT ["nginx"]
CMD ["/ehoney/dockerStart.sh", "$CONFIGS"]