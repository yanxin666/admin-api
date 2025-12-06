package mqdef

import (
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/grouper"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/types"
	"time"
)

// 新增消费者组时需要在此定义
var (
	// GroupTest 测试专用 todo 后续删除
	GroupTest = grouper.NewPushGroup(types.PushGroupConfig{
		BaseGroupConfig: types.BaseGroupConfig{
			GroupName:        "muse_ability_simple", // 消费者组名称
			ConsumerType:     grouper.PushConsumer,  // 消费者类型
			IsFixedGroupName: true,                  // 组名固定，不随环境变量变化
			Timeout:          30 * time.Second,      // 超时时间
			AwaitDuration:    10 * time.Second,
		},
		MaxThreads:                 20,              // 消费最大线程数，默认：20
		MaxCacheMessageCount:       500,             // 本地缓存消息数量上限，单位：条数，
		MaxCacheMessageSizeInBytes: 1024 * 1024 * 5, // 本地缓存消息大小上限，单位：字节，
	})

	// GroupAdmin 后台服务消费者组 muse_admin
	GroupAdmin = grouper.NewPushGroup(types.PushGroupConfig{
		BaseGroupConfig: types.BaseGroupConfig{
			GroupName:     "muse_admin",         // 消费者组名称
			ConsumerType:  grouper.PushConsumer, // 消费者类型
			Timeout:       30 * time.Second,     // 超时时间
			AwaitDuration: 10 * time.Second,
		},
		MaxThreads:                 20,              // 消费最大线程数，默认：20
		MaxCacheMessageCount:       500,             // 本地缓存消息数量上限，单位：条数，
		MaxCacheMessageSizeInBytes: 1024 * 1024 * 5, // 本地缓存消息大小上限，单位：字节，
	})

	// GroupAdminProToTest 由【线上环境】向【测试环境】发送数据的消费者组
	GroupAdminProToTest = grouper.NewPushGroup(types.PushGroupConfig{
		BaseGroupConfig: types.BaseGroupConfig{
			GroupName:        "group_hub_consumer", // 接收由【线上环境】向【测试环境】发送数据进行消费的消费者组
			ConsumerType:     grouper.PushConsumer, // 消费者类型
			IsFixedGroupName: true,                 // 组名固定，不随环境变量变化
			Timeout:          30 * time.Second,     // 超时时间
			AwaitDuration:    10 * time.Second,
		},
		MaxThreads:                 20,              // 消费最大线程数，默认：20
		MaxCacheMessageCount:       500,             // 本地缓存消息数量上限，单位：条数，
		MaxCacheMessageSizeInBytes: 1024 * 1024 * 5, // 本地缓存消息大小上限，单位：字节，
	})

	// GroupAdminProToPre 由【线上环境】向【预发环境】发送数据的消费者组
	GroupAdminProToPre = grouper.NewPushGroup(types.PushGroupConfig{
		BaseGroupConfig: types.BaseGroupConfig{
			GroupName:        "group_hub_consumer_online2pre", // 接收由【线上环境】向【测试环境】发送数据进行消费的消费者组
			ConsumerType:     grouper.PushConsumer,            // 消费者类型
			IsFixedGroupName: true,                            // 组名固定，不随环境变量变化
			Timeout:          30 * time.Second,                // 超时时间
			AwaitDuration:    10 * time.Second,
		},
		MaxThreads:                 20,              // 消费最大线程数，默认：20
		MaxCacheMessageCount:       500,             // 本地缓存消息数量上限，单位：条数，
		MaxCacheMessageSizeInBytes: 1024 * 1024 * 5, // 本地缓存消息大小上限，单位：字节，
	})
)
