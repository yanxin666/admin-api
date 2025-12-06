# 开发规范

## 基本规范

**PS：该部分规范，强制遵守，确保整体代码风格统一。**

- 文件命名规范：文件名全部小写，以_做分词，不要使用拼音和无关单词。
- 文件夹命名规范：文件夹名全部小写，不要包含下划线和数字，不要使用拼音和无关单词。
- 变量命名规范：驼峰命名，首字母小写，不要包含下划线和数字，不要使用拼音和无关单词。
- 常量命名规范：驼峰命名，首字母大写，不要包含下划线和数字，不要使用拼音和无关单词。
- 函数/结构体命名规范：驼峰命名，可导出首字母大写，不可导出小写，不要包含下划线和数字，不要使用拼音和无关单词。
- 代码注释规范：注释是必要的，在复杂逻辑或者BUG修改处注释尽量详细，写清楚代码思路和问题点；代码有修改时注释要及时更新；且每行注释后要跟一个前导空格。
- 异常处理规范：遵循有错必处理的规则，禁止使用_忽略遇到的错误。返回错误时要提供清晰的错误信息。
- 日志记录规范：当遇到错误或者recover panic的时候应使用日志包记录对应等级的错误。

## 建议规范

**PS：该部分规范推荐大家使用、遵守。**

- 目录创建规范：建议不要在```/internal```以外的目录创建自定义文件和目录。
- 单元测试规范：复杂逻辑或者头脑风暴产出的方法和函数建议编写单元测试函数并进行测试。
- 全局变量规范：开发中建议避免使用全局变量。
- 接口命名规范：接口的命名建议以"er"结尾，例如：```type Reader interface{}```
- 包命名规范：包命名使用全小写[尽量使用全小写]，若全小写表明不清其含义，则用下划线来区分[尽量不要]。

**接收器命名：当为方法定义接收器时不建议使用this、self来命名接收器。建议使用接收类型名称首字母小写。**
例如：```func (s *Struct) methodName()```

错误变量检查：在处理错误时，建议使用具体的错误变量检查，而不是直接使用字符串比较。例如：```if err == ErrRecordNotFound {}```

- 错误处理返回值：在函数有多个返回值时，建议将错误作为最后一个返回值。例如：```func Func() (Type,error)```
- 空指针检查：在访问指针类型的字段或调用方法之前，建议检查指针是否为```nil```，以避免空指针异常。
- Panic使用：建议尽量不要使用```panic```，除非你知道你在做什么。

## 细则

此处详细描述一些具体规范的细则。

**后缀为.api的文件里需要注意：**

- 接口报文格式。统一一下：Get请求传参使用form格式，Post请求传参统一使用JSON。
- 结构体及对应的报文，以驼峰标识来命名:

```
type UserData struct {
    UserLevelID   int64  `json:"userLevelID"`
    UserAvatarURL string `json:"userAvatarURL"`
}
```

- 针对一些特殊含义的缩写，建议以全大写，比如ID、HTTP、URL等。

### 代码分层

可分成以下四层：

- Handler [handler]
    - 请求处理层，主要职责是对用户的请求做出响应，不做具体业务逻辑。
- Logic [logic]
    - 业务逻辑块的处理
       ``` golang
      // 校验参数
      
      // 执行service层的 test1()
      
      // 执行service层的 test2()
      
      // ...
      ```
- Service [svc]
    - Logic逻辑块的具体实现
    - 与Model交互
      ``` golang
      func test1(){
        // 数据处理
      }
      
      func test2(){
        // model查表
      }
      
      // ...
      ```
- Model 层 [model]
    - 数据访问层，与MySQL进行数据交互。一个Model应该只处理和自己相关的操作，不做任何组装，更不涉及任何业务逻辑

### 数据库

关于数据库的一些基本共识：

1、数据库表名定义格式，单个单词，多个单词用下划线（_）进行分割，如：u_user，u_order，u_log等，
表名前缀采用业务字母简写(一个或两个字母)，在生成代码Model时采用goctl（公司内部二次开发工具，工具安装参考readme中go-zero操作手册）将前缀去掉。

2、一般情况，建议在数据库建表时，统一保留记录的创建时间（created_at）和更新时间（updated_at），类型统一为datetime。

3、富文本字段，尽量采用utf8mb4编码格式。类型、状态等字段，尽量采用tinyint类型。

#### 后缀为.sql的文件里需要注意：

- 字段：
    - 每个表里关联的外键格式为：表名_id。eg：uid(bad) -> user_id(good)
- int类型：
    - id(primary key)  统一用 bigint
    - 枚举值(enumerate) 统一用 tinyint
    - 且使用 unsigned 标识
- 索引：
    - 唯一索引使用 uni 开头：uni_xxx
    - 普通索引使用 idx 开头：idx_xx
- 字符集：
    - 使用 uft8mb4_general_ci 且为 not null
    - string 类型 默认为空串
    - int 类型 默认为0
- 时间：
    - 使用 datetime 格式
    - 统一使用 created_at 和 updated_at
    - created_at datetime xxx
    - updated_at datetime xxx on update xxx
- 以下为参考实例:

```sql
CREATE TABLE `u_info`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
    `nickname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
    `avatar` varchar(256) NOT NULL DEFAULT '' COMMENT '用户头像',
    `country_code` tinyint unsigned NOT NULL DEFAULT '86' COMMENT '国家代码',
    `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
    `source` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '来源 1.APP',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_phone` (`phone`),
    KEY `idx_phone_source` (`phone`, `source`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci COMMENT ='用户信息表';

```

### 指针 or 结构体
对于结构体这种类型，指针 or 结构体，无论使用哪一个，都尽量保持统一。保证数据的安全，也防止乱用造成的代码阅读感混乱。所以目前推荐的选择方式如下。

约定内存=N, 单位为字节 (字节数大家约定)
1. 内存 <= N 时，优先使用结构体
2. 内存 > N 时，必须使用指针（极特殊情况除外）

使用指针时一定注意 ：
- 使用前使用 nil 判断，防止 panic
- 需要明确对数据是否有影响， 防止在函数内部对参数修改，导致外部受影响
- 是否是并发安全的，当需要的时候可以使用深度拷贝等方式解决一些共享等读写问题

### 单元测试

关于单元测试一些的基本共识：

1、鉴于本项目采用go-zero框架，所有的API都被抽象到xxx_logic.go当中，所有的API业务逻辑都可比较容易实现单元测试，建议后端开发人员，针对对外接口在同目录下提供单元对应的测试。

2、较复杂的核心业务逻辑，均需提供单元测试；

3、一般业务逻辑实现可根据具体情况，确定是否编写单元测试。

## git开发规范

### 注释格式

提交代码注释需按照指定格式，对应正则表达式如下：
```text
^T\d{4}(:|：).{1,100}
// （1）第一部分为任务编号格式为：T0033
// （2）第二部分为冒号，支持中英文
// （3）第三部分为具体说明，限定1到100个字符

示例：
T0033：修改代码规范

PS：其中T0033为通用修改代码规范的Task，T0034为修改基础架构功能的Task
```

```text
开发前，先切换到主分支master，并拉取最新代码，再创建对应任务的开发分支，开发自测结束之后，在当前开发分支上合并远端master分支，如有冲突需解决冲突，合并成功后，
再切换到master分支，先更新本地master分支，再在master上合并开发分支，如不需审核则直接合并成功，否则需联系管理员审核合并
```

- 切换主干分支

```text
git checkout master
```

- 拉取最新master代码

```text
git pull
```

- 创建并切换到开发分支

```text
git checkout -b devBranch
```

- 开发结束合并远端master分支

```text
git merge origin master
```

- 切换到master分支

```text
git checkout master
```

- 更新本地master分支

```text
git pull
```

- 合并开发分支

```text
git merge devBranch
```

- 提交分支或发起合并申请

```text
git push
```

## 报文规范

### 接口路径

1、请求路径

- 前缀：`/v1`，后续有需要区分版本再依次递增编号

### 请求响应

```text
// 定义接口的请求体
DemoReq {
    UserName string `json:"userName"`
}

// 定义接口的响应体
DemoResp {
    Id       int64  `json:"id"`
    Name     string `json:"name"`
    Token    string `json:"token"`
    ExpireAt string `json:"expireAt"`
}
```

1、请求体规则

- 统一后缀`Req`，如：`DemoReq`
- 请求体参数使用`驼峰`命名法，如`UserName`，字段类型按照实际场景定义
- `POST`请求时统一使用`json`关键字接收参数，`GET`请求时统一使用`form`关键字接收参数，参数名统一使用`驼峰`命名法

2、响应体规则

- 统一后缀`Resp`，如：`DemoResp`
- 响应体参数使用`驼峰`命名法，如`ExpireAt`，字段类型按照实际场景定义
- 响应体参数统一使用`json`格式返回，参数名统一使用`驼峰`命名法

### Header

1、请求体规则

- `GET`请求无需设置
- `POST`需设置header`"Content-Type": "application/json"`

2、响应体规则

- 默认 `Content-Type` 为 `application/json; charset=utf-8`

3、token认证

- header发送jwt授权token：`Authorization: Bearer <token>`

4、版本号

- header发送版本号：`Version: <version>`，如：Version: 1.0
- header发送平台：`Platform: <platform>`，如：Platform: ios

大家在实践的过程中，可不断完善、丰富该文档中的最佳实践案例。

## 日志规范

### <span style="color:green">强烈推荐写法</span>
<H4 style="color:green">带上下文的日志记录，可完整记录链路追踪相关的trace，方便日志查询</h4>

1、直接使用logc
```go
logc.Info(ctx, "logc 日志记录")
```

带上下文的日志内容
```text
2024-01-06T11:27:18.594+08:00 info logc 测试日志记录 caller=user/opinion.go:33 trace=f8fb31fe8bde4dd2bab979533efab4b8 span=3cddd74920e62b24   path=/v1/user/opinion
```

### <span style="color:yellow">不推荐写法</span>
<H4 style="color:yellow">带上下文的日志记录，可完整记录链路追踪相关的trace，但相较logc多一步调用，容易遗忘，因此不推荐</h4>

1、使用logx.WithContext()
```go
logx.WithContext(ctx).Info("logx WithContext 日志记录")
```

2、使用注入了上下文的Logger对象
```go
Logger := logx.WithContext(ctx)
Logger.Info("Logger 日志记录")
```

带上下文的日志内容
```text
2024-01-06T11:27:18.594+08:00 info logc 测试日志记录 caller=user/opinion.go:33 trace=f8fb31fe8bde4dd2bab979533efab4b8 span=3cddd74920e62b24   path=/v1/user/opinion
```

### <span style="color:red">严格禁用写法</span>
<H4 style="color:red">不带上下文的日志记录，没有链路追踪相关的trace，不利于日志查询</h4>

1、直接使用logx
```go
logx.Info("logx 日志记录")
```

不带上下文的日志内容
```text
2024-01-06T11:27:18.594+08:00    info   logx 测试日志记录       caller=user/opinion.go:32
```

## 最佳实践

更多"编程最佳实践"，参考文档：https://douyuxingchen.feishu.cn/docx/PnsEdzPpgo4NdFxLlVLcvIyTnAw

## Something from official golang documentation

### Package Naming

Good package names are short and clear. They are lower case, with no under_scores or mixedCaps. They are often simple
nouns, such as:

`time` (provides functionality for measuring and displaying time)

`list` (implements a doubly linked list)

`http` (provides HTTP client and server implementations)
