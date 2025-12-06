package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Salt         string
	JwtAuth      Auth
	Mysql        Mysql
	Cache        cache.CacheConf
	Redis        redis.RedisConf
	Cls          Cls
	Oss          Oss
	TencentCloud TencentCloud
	// Tts          Tts
	EncryptKey string
	Im         Im
	XxlJob     XxlJob
	RocketMQ   MQConfig
	RpcTarget  RpcTarget
	MiniMax    MiniMax
}

type Mysql struct {
	DataSourceAdmin, DataSourceAbility, DataSourceCenter, DataSourceHub, DataSourceKingClub string
	MaxOpenConn                                                                             int
	MaxIdleConn                                                                             int
	ConnMaxLifetime                                                                         int
	ConnMaxIdleTime                                                                         int
}

type Cls struct {
	Endpoint string
	TopicID  string
}

type Auth struct {
	AccessSecret string
	AccessExpire int64
}

type Oss struct {
	SecretId       string
	SecretKey      string
	Appid          int
	Bucket         string
	BehaviorBucket string
	Region         string
	ReplaceDomain  string
}

type TencentCloud struct {
	SecretId     string
	SecretKey    string
	SmsAppId     string
	SmsSign      string
	SmsSize      int // 验证码长度
	SmsSendLimit int // 当天最大发送量
}

type Tts struct {
	SecretId  string
	SecretKey string
	Region    string
	Appid     int
}

type Im struct {
	AppId int
	Key   string
}

type XxlJob struct {
	Host, Port, AccessToken, RegistryKey string
}

type MQConfig struct {
	EndPoint          string // 消息队列地址
	AccessKey         string // 消息队列AK
	SecretKey         string // 消息队列SK
	ProducerRetryTime int    // 生产者重试次数
	ProducerTimeout   int    // 生产者超时时间(秒)
	LogLevel          string // 日志级别 // debug:调试，info:信息，warn:警告，error:错误，fatal:致命错误
	LogRoot           string // 日志根目录
	LogConsole        string // 是否开启控制台日志输出，true：开启，false: 关闭
	Env               string // 环境变量，pro: 生产环境，dev: 开发环境，test: 测试环境
}

type RpcTarget struct {
	Ability  string
	Passport string
	Core     string
}

type MiniMax struct {
	ApiKey   string
	GroupIds struct {
		CiYuan string
	}
}
