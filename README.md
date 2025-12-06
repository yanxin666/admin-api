# muse-admin

能力探真项目管理后台

- Doc：[操作手册](https://docs.arklnk.com)
- 开源接口文档：[初始化接口文档](https://apifox.com/apidoc/shared-ddbffed9-fa11-49ed-bc46-8d76ab058d60)
- 开源地址：[https://github.com/arklnk/muse-admin](https://github.com/arklnk/muse-admin)
- 预览地址：[http://arkadmin.si-yee.com](http://arkadmin.si-yee.com)

  | 账号     | 密码     | 备注    |
      |--------|--------|-------|
  | admin | 123456 | 超级管理员 |
  | demo   | 123456 | 演示账号  |

## 基础环境

- Golang版本：1.21+
- 基础框架：go-zero

## 项目介绍

### 目录结构

目录结构遵循 [golang-standards/project-layout](https://github.com/golang-standards/project-layout) 规范

```text
.
├── api API服务配置目录
│    ├── routes 应用api 
│    │    ├── demo
│    │    ├── user
│    │    └── main.api   路由入口文件
│    └── swagger swagger api
├── docs 文档
├── etc 配置项
├── internal 项目私有目录
│    ├── config 配置文件
│    ├── define 常量定义
│    ├── handler api逻辑入口
│    │    ├── demo
│    │    └── user
│    ├── logic api逻辑组装
│    │    ├── demo
│    │    └── user
│    ├── middleware 中间件
│    ├── model 数据模型
│    │    └── user
│    ├── svc 服务层
│    │    ├── oauth
│    │    └── user
│    └── types 结构体定义
└── pkg 外部资源包，该包下的所有功能均需注释说明功能以及单元测试
    ├── arch 资源初始化
    │    ├── http
    │    ├── kafka
    │    ├── mysql
    │    └── redis
    ├── cron 定时任务
    ├── errs 错误处理
    ├── metric 把脉监听
    ├── response 返回处理
    ├── safe 安全处理
    ├── third 第三方服务
    │    ├── gpt 
    │    └── wechat
    ├── wire
    ├── zap
    └── util 工具包
```

### 前提准备

```
1.克隆此二开工具
git clone git@e.coding.net:douyu-devops/zhiyan/goctl.git

2.使用此二开工具生成的二进制文件 goctl 来生成项目文件

3.为开发方便，建议将二进制文件 goctl 放在环境变量中
```

### 常用命令

- 生成API服务示例

```shell
goctl api format --dir ./api/routes; goctl api go --style go_zero --api ./api/routes/main.api --dir .; goctl api plugin -plugin goctl-swagger="swagger -filename ./api/swagger/swagger.json -host 49.232.253.114:8210" -api ./api/routes/main.api -dir .
```

- 生成rpc pb示例

```shell
# 定义生成目录 && 生成proto文件对应的pb.go文件（用于生成RPC的client）
dir="ability" && protoc ./proto/"$dir"/*.proto --go_out=./proto --go-grpc_out=./proto
```

- 生成model示例

```shell
# 默认.sql文件会以 去掉表前缀 生成，指令如下：
goctl model mysql ddl --style go_zero --src ./internal/model/user/sql_user.sql --dir ./internal/model/user

# [前提：文件里只有单个ddl才可用该指令生成自定义文件] 若.sql文件想自定义生成文件名及model中结构体，可用 --typ --name满足，指令如下：
goctl model mysql ddl --style go_zero --src ./internal/model/user/usre.sql --dir ./internal/model/user --typ structName[结构体名] --name fileName[文件名]
```

- 生成swagger

```shell
goctl api plugin -plugin goctl-swagger="swagger -filename ./api/swagger/swagger.json -host 49.232.253.114:8210" -api ./api/routes/main.api -dir .
```

## 开发规范

开发基础规范，参考<a href="./docs/guidelines.md">guidelines.md</a>中的说明。<span style="color:red">
一定要仔细认真阅读，保持代码规范很重要！！！</span>

## 参考文档

- go-zero官方文档：https://go-zero.dev/
- 内部整理手册及使用案例文档：https://douyuxingchen.feishu.cn/docx/NUkDd8nTQoKxqzxH752cAgrjnAe

