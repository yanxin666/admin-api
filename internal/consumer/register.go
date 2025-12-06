package consumer

import (
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/grouper"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/producer"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/types"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"muse-admin/internal/config"
	"muse-admin/internal/consumer/hub"
	"muse-admin/internal/consumer/hub/builder"
	"muse-admin/internal/define/mqdef"
	"muse-admin/internal/svc"
	"os"
)

func Register(c config.Config, svcCtx *svc.ServiceContext) {
	// 初始化自定义MQ日志
	_ = os.Setenv(rmq.CLIENT_LOG_LEVEL, c.RocketMQ.LogLevel)
	_ = os.Setenv(rmq.CLIENT_LOG_ROOT, c.RocketMQ.LogRoot)
	_ = os.Setenv(rmq.ENABLE_CONSOLE_APPENDER, c.RocketMQ.LogConsole)
	rmq.ResetLogger()

	// 注册MQ生产者
	mqBase := buildMQBaseConfig(c.RocketMQ)
	svcCtx.MQProducer = producer.NewProducerV5(mqBase)

	// 慎用！MQ测试环境生产者，用作【线上环境】向【测试环境】发送数据
	mqBaseToTest := buildMQConfProToTest(c.RocketMQ)
	svcCtx.MQProducerToTest = producer.NewProducerV5(mqBaseToTest, producer.WithAlias("admin_pro_to_test"))

	// 慎用！MQ预发环境生产者，用作【线上环境】向【预发环境】发送数据
	mqBaseToPre := buildMQConfProToPre(c.RocketMQ)
	svcCtx.MQProducerToPre = producer.NewProducerV5(mqBaseToPre, producer.WithAlias("admin_pro_to_pre"))

	// 注册MQ消费者
	registerQueue(mqBase, mqBaseToTest, mqBaseToPre, svcCtx)

	// // 注册MQ消费
	// registerQueue(buildConfigMQ(c, svcCtx), svcCtx)

	// 注册其他服务
	// ...
}

func buildMQBaseConfig(svcCtxMQ config.MQConfig) types.BaseConfig {
	return types.BaseConfig{
		EndPoint:          svcCtxMQ.EndPoint,
		AccessKey:         svcCtxMQ.AccessKey,
		SecretKey:         svcCtxMQ.SecretKey,
		ProducerRetryTime: svcCtxMQ.ProducerRetryTime,
		ProducerTimeout:   svcCtxMQ.ProducerTimeout,
		LogLevel:          svcCtxMQ.LogLevel,
		LogRoot:           svcCtxMQ.LogRoot,
		LogConsole:        svcCtxMQ.LogConsole,
		Env:               svcCtxMQ.Env,
	}
}

// buildMQConfProToTest 构建生产环境向测试环境的发送的MQ配置
func buildMQConfProToTest(svcCtxMQ config.MQConfig) types.BaseConfig {
	point := "rmq-jaz2q432.rocketmq.bj.public.tencenttdmq.com:8080"
	accessKey := "akjaz2q4323fe4c263e2bd"
	secretKey := "skf535d5b64980290a"

	return types.BaseConfig{
		EndPoint:          point,
		AccessKey:         accessKey,
		SecretKey:         secretKey,
		ProducerRetryTime: svcCtxMQ.ProducerRetryTime,
		ProducerTimeout:   svcCtxMQ.ProducerTimeout,
		LogLevel:          svcCtxMQ.LogLevel,
		LogRoot:           svcCtxMQ.LogRoot,
		LogConsole:        svcCtxMQ.LogConsole,
		Env:               svcCtxMQ.Env,
	}
}

// buildMQConfProToPre 构建生产环境向预发环境的发送的MQ配置
func buildMQConfProToPre(svcCtxMQ config.MQConfig) types.BaseConfig {
	point := "rmq-8v9gnzp9o.rocketmq.bj.qcloud.tencenttdmq.com:8080"
	accessKey := "ak8v9gnzp9ocb22a29a9c61"
	secretKey := "sk9b80b588d8c7ffd5"

	return types.BaseConfig{
		EndPoint:          point,
		AccessKey:         accessKey,
		SecretKey:         secretKey,
		ProducerRetryTime: svcCtxMQ.ProducerRetryTime,
		ProducerTimeout:   svcCtxMQ.ProducerTimeout,
		LogLevel:          svcCtxMQ.LogLevel,
		LogRoot:           svcCtxMQ.LogRoot,
		LogConsole:        svcCtxMQ.LogConsole,
		Env:               svcCtxMQ.Env,
	}
}

func registerQueue(c, c2t, c2p types.BaseConfig, svcCtx *svc.ServiceContext) {
	// 设置MQ日志等级
	rlog.SetLogLevel(c.LogLevel)

	/*------------------- 注册所需要消费的Topic -------------------*/

	// mq.Register(define.ConsumerTopic[c.Mode].ExcelImport, NewImport(svcCtx))         // Excel导入任务[2024.10.11删除，换成go func处理]

	// todo 测试使用，后续删除
	mqdef.GroupTest.Register(mqdef.TopicTest, func() types.Consumer { return NewTest(svcCtx) })

	// 消费【三方平台】向【线上环境】发送的数据
	mqdef.GroupAdmin.Register(mqdef.TopicScheduleKnowledge, func() types.Consumer {
		return hub.NewDataHub(svcCtx, builder.NewGenWithOpt(builder.WithReqFrom("ThirdToOnline")))
	})
	// 消费由【线上环境】向【测试环境】发送的数据
	mqdef.GroupAdminProToTest.Register(mqdef.TopicHubConsumerToTest, func() types.Consumer {
		return hub.NewDataHub(svcCtx, builder.NewGenWithOpt(builder.WithReqFrom("OnlineToTest")))
	})
	// 消费由【线上环境】向【预发环境】发送的数据
	mqdef.GroupAdminProToPre.Register(mqdef.TopicHubConsumerToPre, func() types.Consumer {
		return hub.NewDataHub(svcCtx, builder.NewGenWithOpt(builder.WithReqFrom("OnlineToPre")))
	})
	// 消费由【管理后台审核通过后】向【线上环境】发送的数据
	mqdef.GroupAdmin.Register(mqdef.TopicHubConsumerToPro, func() types.Consumer {
		return hub.NewDataHub(svcCtx, builder.NewGenWithOpt(builder.WithReqFrom("ManageToOnline")))
	})

	// 注册由【线上环境】发送的消息在【测试环境】进行消费的消费组
	if svcCtx.Config.Mode == "test" {
		gp := mqdef.GroupAdminProToTest.(*grouper.PushGroup)
		go mq.NewPushConsume().Start(gp.GetThreads(), gp, c2t)
	}

	// 注册由【线上环境】发送的消息在【预发环境】进行消费的消费组
	if svcCtx.Config.Mode == "pre" {
		gp := mqdef.GroupAdminProToPre.(*grouper.PushGroup)
		go mq.NewPushConsume().Start(gp.GetThreads(), gp, c2p)
	}

	// 注册消费者组
	grouper.AddGroups(
		mqdef.GroupTest,
		mqdef.GroupAdmin,
	)

	// 初始化全局消费者
	mq.StartConsumer(c)
}

// func buildConfigMQ(c config.Config, svcCtx *svc.ServiceContext) mq.ConfigMQ {
// 	configMQ := mq.ConfigMQ{
// 		AccessKey:         c.RocketMQ.AccessKey,
// 		SecretKey:         c.RocketMQ.SecretKey,
// 		EndPoint:          c.RocketMQ.EndPoint,
// 		ProducerRetryTime: c.RocketMQ.ProducerRetryTime,
// 		ProducerTimeout:   c.RocketMQ.ProducerTimeout,
// 		LogLevel:          c.RocketMQ.LogLevel,
// 		Mode:              c.Mode,
// 		ClsConfig:         buildClsConfig(c),
// 	}
// 	var groupMQ []mq.GroupMQ
// 	for _, group := range c.RocketMQ.ConsumerGroup {
// 		groupMQ = append(groupMQ, mq.GroupMQ{
// 			GroupName:           group.GroupName,
// 			Threads:             group.Threads,
// 			Retry:               group.Retry,
// 			Timeout:             group.Timeout,
// 			GoroutineNums:       group.GoroutineNums,
// 			MessageBatchMaxSize: group.MessageBatchMaxSize,
// 		})
// 	}
//
// 	// 注册由【线上环境】发送的消息在【测试环境】进行消费的消费组
// 	if c.Mode == "test" {
// 		groupMQ = append(groupMQ, mq.GroupMQ{
// 			GroupName:           define.GroupHubConsumer,
// 			Threads:             1, // 目前此group先开一个协程去消费
// 			Retry:               groupMQ[0].Retry,
// 			Timeout:             groupMQ[0].Timeout,
// 			GoroutineNums:       groupMQ[0].GoroutineNums,
// 			MessageBatchMaxSize: groupMQ[0].MessageBatchMaxSize,
// 		})
// 	}
//
// 	configMQ.ConsumerGroup = groupMQ
//
// 	// 注册MQ消费所需要的数据库
// 	configMQ.ConfigDB = buildDBConfig(svcCtx)
//
// 	return configMQ
// }
//
// func buildClsConfig(c config.Config) config.AutoCls {
// 	return config.AutoCls{
// 		SecretId:  c.TencentCloud.SecretId,
// 		SecretKey: c.TencentCloud.SecretKey,
// 		TopicID:   c.Cls.TopicID,
// 		Endpoint:  c.Cls.Endpoint,
// 		Mode:      c.Mode,
// 	}
// }
//
// func buildDBConfig(svcCtx *svc.ServiceContext) mq.ConfigDB {
// 	return mq.ConfigDB{
// 		SysUserModel: svcCtx.SysUserModel,
// 	}
// }
