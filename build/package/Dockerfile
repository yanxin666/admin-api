# 使用golang:1.21.1-alpine镜像作为基础镜像，并命名为builder
FROM registry-ze.tencentcloudcr.com/basic/golang:1.24.2-alpine3.21-git AS builder

# 设置标签为gobuilder
LABEL stage=gobuilder

# 禁用CGO，以提高构建速度
ENV CGO_ENABLED=0
ENV GOAUTH=netrc

# 设置Go代理为https://goproxy.cn和direct
ENV GOPROXY https://goproxy.cn,direct

## 设置私有仓库链接
ENV GOPRIVATE="e.coding.net/zmexing/*"
#
## 检查是否已经安装 Git，如果没有则安装
#RUN command -v git || apk add --no-cache git
#
## 在容器中运行 git config 命令，将https方式替换为ssh方式访问
#RUN git config --add --global url."git@e.coding.net:".insteadOf https://e.coding.net

# 将Alpine Linux的仓库地址替换为阿里云的镜像源，加速镜像资源拉取速度
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 更新APK包列表并安装tzdata
RUN apk update --no-cache && apk add --no-cache tzdata

# 设置工作目录为/build
WORKDIR /build

# 将go.mod文件复制到容器中
ADD go.mod .

# 将go.sum文件复制到工作容器中
ADD go.sum .

# 下载依赖
RUN go mod download

# 将当前目录下的所有文件复制到 /build 中
COPY . .

# 将etc目录下的文件复制到容器中的app/etc目录下
COPY ./etc /app/etc

# 使用go build命令编译应用程序，并将生成的可执行文件命名为application
RUN go build -ldflags="-s -w" -o /app/application application.go

# 使用alpine镜像作为基础镜像
FROM registry-ze.tencentcloudcr.com/basic/alpine

# 将builder镜像中的ca-certificates.crt文件复制到容器中的/etc/ssl/certs/目录下
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# 将builder镜像中的Asia/Shanghai时区文件复制到容器中的/usr/share/zoneinfo/目录下
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai

# 设置时区环境变量TZ为Asia/Shanghai
ENV TZ Asia/Shanghai

# 设置工作目录为/app
WORKDIR /app

# 文件挂载，如日志文件路径
VOLUME /tmp

# 将builder镜像中的application文件复制到容器中的/app/目录下
COPY --from=builder /app/application /app/application

# 将builder镜像中的etc目录复制到容器中的/app/目录下
COPY --from=builder /app/etc /app/etc

# 设置容器启动时执行的命令，运行application程序，并传入配置文件application.yaml
CMD ["./application", "-f", "etc/application.yaml"]
