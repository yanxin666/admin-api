package svc

import (
	"context"
	"fmt"
	"math/rand"
	"muse-admin/internal/config"
	"muse-admin/internal/middleware"
	"muse-admin/internal/model/behavior"
	benefitModel "muse-admin/internal/model/benefit"
	"muse-admin/internal/model/code"
	"muse-admin/internal/model/hub"
	"muse-admin/internal/model/kingclub/supertrain"
	knowledgeModel "muse-admin/internal/model/knowledge"
	lessonModel "muse-admin/internal/model/lesson"
	"muse-admin/internal/model/live"
	"muse-admin/internal/model/live_course"
	"muse-admin/internal/model/member"
	"muse-admin/internal/model/nightstar"
	student "muse-admin/internal/model/student"
	"muse-admin/internal/model/supervisor"
	systemModel "muse-admin/internal/model/system"
	taskModel "muse-admin/internal/model/task"
	"muse-admin/internal/model/teacher"
	userModel "muse-admin/internal/model/user"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/model/write_ppt"
	"os"
	"strings"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mq/producer"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/mysql"
	archRedis "e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/redis"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/proto/ability"
	"e.coding.net/zmexing/nenglitanzhen/proto/core"
	"e.coding.net/zmexing/nenglitanzhen/proto/passport"
	"github.com/go-pay/util/snowflake"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config       config.Config
	Redis        *redis.Redis
	RedisClient  *archRedis.ClientRedis
	PermMenuAuth rest.Middleware

	MQProducer       *producer.ProducerV5
	MQProducerToTest *producer.ProducerV5
	MQProducerToPre  *producer.ProducerV5

	Snowflake *snowflake.Node

	MysqlConnAdmin    sqlx.SqlConn
	MysqlConnAbility  sqlx.SqlConn
	MysqlConnCenter   sqlx.SqlConn
	MysqlConnHub      sqlx.SqlConn
	MysqlConnKingClub sqlx.SqlConn

	// muse-admin 数据库
	SysUserModel        workbench.UserModel
	SysPermMenuModel    workbench.PermMenuModel
	SysRoleModel        workbench.RoleModel
	SysDeptModel        workbench.DeptModel
	SysJobModel         workbench.JobModel
	SysProfessionModel  workbench.ProfessionModel
	SysDictionaryModel  workbench.DictionaryModel
	SysLogModel         workbench.LogModel
	StudentModel        student.UserInfoModel
	StudentOpinionModel student.UserOpinionModel
	SyncTaskModel       taskModel.SyncTaskModel
	SyncTaskLogModel    taskModel.SyncTaskLogModel
	ScheduleLogModel    taskModel.ScheduleLogModel
	SessionModel        behavior.SessionModel
	SessionUserModel    behavior.SessionUserModel
	SessionRecordModel  behavior.SessionRecordModel

	// muse-center 数据库
	BaseUserModel  userModel.BaseUserModel
	UserModel      userModel.UserModel
	UserLogModel   userModel.UserLoginLogModel
	UserWhiteModel userModel.WhiteListModel
	DsTeacherModel teacher.TeacherModel

	// muse-ability 数据库
	MaterialModel                 knowledgeModel.MaterialModel
	ExampleModel                  knowledgeModel.ExampleModel
	LessonQuestionModel           knowledgeModel.LessonQuestionModel
	LessonQuestionOptionModel     knowledgeModel.LessonQuestionOptionModel
	LessonQuestionTtsModel        knowledgeModel.LessonQuestionTtsModel
	CourseOutlineModel            lessonModel.CourseOutlineModel
	LessonModel                   knowledgeModel.LessonModel
	LessonResourceModel           knowledgeModel.LessonResourceModel
	LessonPointModel              knowledgeModel.LessonPointModel
	ChannelOrderModel             benefitModel.ChannelOrderModel
	ChannelSubOrderModel          benefitModel.ChannelSubOrderModel
	SysSmsRecordModel             systemModel.SmsRecordModel
	BenefitResourceModel          benefitModel.BenefitsResourceModel
	BenefitGroupModel             benefitModel.BenefitsGroupModel
	UserBenefitModel              benefitModel.UserBenefitsDetailModel
	UserBenefitRecordModel        benefitModel.BenefitsUserRecordModel
	PPTCourseSeriesModel          write_ppt.GradeCourseSeriesModel
	PPTCourseOutlineModel         write_ppt.CourseOutlineModel
	PPTLessonModel                write_ppt.LessonModel
	DouShenVipModel               member.DoushenVipModel
	LiveTeacherModel              nightstar.LiveTeacherModel
	LiveCourseModel               nightstar.LiveCourseModel
	LiveCourseDetailModel         nightstar.LiveCourseDetailModel
	LiveCourseTimelineModel       nightstar.LiveCourseTimelineModel
	LiveCourseEventModel          nightstar.LiveCourseEventModel
	LiveCourseQuestionModel       nightstar.LiveCourseQuestionModel
	LiveCourseQuestionOptionModel nightstar.LiveCourseQuestionOptionModel
	LiveBetaRecordModel           live.LiveBetaRecordsModel
	RedeemCodeModel               code.RedeemCodeModel
	DictWhiteModel                systemModel.DictWhiteModel

	// muse-hub 数据库
	HubLessonModel   hub.LessonSnapshotModel
	HubWritePPTModel hub.WritePptSnapshotModel
	// HubLessonQuestionTtsModel knowledgeModel.LessonQuestionTtsModel
	HubLiveModel       hub.LiveSnapshotModel
	HubSuperTrainModel hub.SuperTrainSnapshotModel

	// king-club 数据库
	TrainCourseModel                  supertrain.CourseModel
	TrainCourseChapterModel           supertrain.CourseChapterModel
	TrainChapterTaskModel             supertrain.ChapterTaskModel
	TrainSubTaskModel                 supertrain.SubTaskModel
	TrainMedalModel                   supertrain.MedalModel
	TrainSubTaskMedalModel            supertrain.SubTaskMedalModel
	SupeScheduleModel                 supervisor.ScheduleModel
	SupeStreamModel                   supervisor.StreamModel
	SupeScheduleInteractRelationModel supervisor.ScheduleInteractionLinkModel
	AcTeacherModel                    supervisor.TeacherModel
	InteractionModel                  supervisor.InteractionModel
	LiveUserModel                     live_course.UserModel
	ChatRoomModel                     live_course.ChatRoomModel
	// muse-hub 数据库
	// HubLessonQuestionModel    knowledgeModel.LessonQuestionModel
	// HubLessonQuestionTtsModel knowledgeModel.LessonQuestionTtsModel

	// rpc
	AbilityRPC  *ability.Client
	PassportRPC *passport.Client
	CoreRPC     *core.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConnAdmin := mysql.GetMysqlConn(c.Mysql.DataSourceAdmin)
	mysqlConnAbility := mysql.GetMysqlConn(c.Mysql.DataSourceAbility)
	mysqlConnCenter := mysql.GetMysqlConn(c.Mysql.DataSourceCenter)
	mysqlConnHub := mysql.GetMysqlConn(c.Mysql.DataSourceHub)
	mysqlConnKingClub := mysql.GetMysqlConn(c.Mysql.DataSourceKingClub)
	redisClient := redis.New(c.Redis.Host, func(r *redis.Redis) {
		r.Type = c.Redis.Type
		r.Pass = c.Redis.Pass
	})

	// 设置日志模式
	logz.SetLogzMode(c.Log.Mode)
	logz.SetEncoding(c.Log.Encoding)
	if c.Log.Mode == logz.File {
		// 初始化日志对象
		logz.SelfSetLogzWriter(c.Log.Path)
	}
	// 关闭统计日志
	stat.DisableLog()

	return &ServiceContext{
		Config:       c,
		Redis:        redisClient,
		RedisClient:  archRedis.NewClientRedisBuilder().SetConnect(c.Redis).Build(), // Redis连接,
		PermMenuAuth: middleware.NewPermMenuAuthMiddleware(c, redisClient).Handle,
		Snowflake:    initSnowflake(), // 雪花算法ID发号器
		// 数据库配置
		MysqlConnAdmin:    mysqlConnAdmin,
		MysqlConnAbility:  mysqlConnAbility,
		MysqlConnCenter:   mysqlConnCenter,
		MysqlConnHub:      mysqlConnHub,
		MysqlConnKingClub: mysqlConnKingClub,
		// muse-admin 数据库
		SysUserModel:        workbench.NewUserModel(mysqlConnAdmin),
		SysPermMenuModel:    workbench.NewPermMenuModel(mysqlConnAdmin),
		SysRoleModel:        workbench.NewRoleModel(mysqlConnAdmin),
		SysDeptModel:        workbench.NewDeptModel(mysqlConnAdmin),
		SysJobModel:         workbench.NewJobModel(mysqlConnAdmin),
		SysProfessionModel:  workbench.NewProfessionModel(mysqlConnAdmin),
		SysDictionaryModel:  workbench.NewDictionaryModel(mysqlConnAdmin),
		SysLogModel:         workbench.NewLogModel(mysqlConnAdmin),
		StudentModel:        student.NewUserInfoModel(mysqlConnAdmin),
		StudentOpinionModel: student.NewUserOpinionModel(mysqlConnAdmin),
		SyncTaskModel:       taskModel.NewSyncTaskModel(mysqlConnAdmin),
		SyncTaskLogModel:    taskModel.NewSyncTaskLogModel(mysqlConnAdmin),
		ScheduleLogModel:    taskModel.NewScheduleLogModel(mysqlConnAdmin),
		SessionModel:        behavior.NewSessionModel(mysqlConnAdmin),
		SessionUserModel:    behavior.NewSessionUserModel(mysqlConnAdmin),
		SessionRecordModel:  behavior.NewSessionRecordModel(mysqlConnAdmin),
		// muse-center 数据库
		BaseUserModel:  userModel.NewBaseUserModel(mysqlConnCenter),
		UserModel:      userModel.NewUserModel(mysqlConnCenter),
		UserLogModel:   userModel.NewUserLoginLogModel(mysqlConnCenter),
		UserWhiteModel: userModel.NewWhiteListModel(mysqlConnCenter),
		DsTeacherModel: teacher.NewTeacherModel(mysqlConnCenter),

		// muse-ability 数据库
		MaterialModel:                 knowledgeModel.NewMaterialModel(mysqlConnAbility),
		ExampleModel:                  knowledgeModel.NewExampleModel(mysqlConnAbility),
		LessonQuestionModel:           knowledgeModel.NewLessonQuestionModel(mysqlConnAbility),
		LessonQuestionOptionModel:     knowledgeModel.NewLessonQuestionOptionModel(mysqlConnAbility),
		LessonQuestionTtsModel:        knowledgeModel.NewLessonQuestionTtsModel(mysqlConnAbility),
		CourseOutlineModel:            lessonModel.NewCourseOutlineModel(mysqlConnAbility),          // 大纲表
		LessonModel:                   knowledgeModel.NewLessonModel(mysqlConnAbility),              // 课节表
		LessonResourceModel:           knowledgeModel.NewLessonResourceModel(mysqlConnAbility),      // 课节资源表
		LessonPointModel:              knowledgeModel.NewLessonPointModel(mysqlConnAbility),         // 课节知识点表
		ChannelOrderModel:             benefitModel.NewChannelOrderModel(mysqlConnAbility),          // 渠道订单表
		ChannelSubOrderModel:          benefitModel.NewChannelSubOrderModel(mysqlConnAbility),       // 渠道子订单表
		SysSmsRecordModel:             systemModel.NewSmsRecordModel(mysqlConnAbility),              // 短信记录表
		BenefitResourceModel:          benefitModel.NewBenefitsResourceModel(mysqlConnAbility),      // 权益资源表
		BenefitGroupModel:             benefitModel.NewBenefitsGroupModel(mysqlConnAbility),         // 权益组表
		UserBenefitModel:              benefitModel.NewUserBenefitsDetailModel(mysqlConnAbility),    // 用户权益表
		UserBenefitRecordModel:        benefitModel.NewBenefitsUserRecordModel(mysqlConnAbility),    // 用户权益记录表
		PPTCourseSeriesModel:          write_ppt.NewGradeCourseSeriesModel(mysqlConnAbility),        // 写作私教年级可学系列表
		PPTCourseOutlineModel:         write_ppt.NewCourseOutlineModel(mysqlConnAbility),            // 写作私教大纲表
		PPTLessonModel:                write_ppt.NewLessonModel(mysqlConnAbility),                   // 写作私教课节表
		DouShenVipModel:               member.NewDoushenVipModel(mysqlConnAbility),                  // 豆神Vip表
		LiveTeacherModel:              nightstar.NewLiveTeacherModel(mysqlConnAbility),              // AI直播老师表
		LiveCourseModel:               nightstar.NewLiveCourseModel(mysqlConnAbility),               // AI直播课程列表
		LiveCourseDetailModel:         nightstar.NewLiveCourseDetailModel(mysqlConnAbility),         // AI直播课程详情表
		LiveCourseTimelineModel:       nightstar.NewLiveCourseTimelineModel(mysqlConnAbility),       // AI直播课程时间线表
		LiveCourseEventModel:          nightstar.NewLiveCourseEventModel(mysqlConnAbility),          // AI直播课程事件表
		LiveCourseQuestionModel:       nightstar.NewLiveCourseQuestionModel(mysqlConnAbility),       // AI直播课程题目表
		LiveCourseQuestionOptionModel: nightstar.NewLiveCourseQuestionOptionModel(mysqlConnAbility), // AI直播课程题目选项表
		LiveBetaRecordModel:           live.NewLiveBetaRecordsModel(mysqlConnAbility),               // AI直播Beta用户记录表
		RedeemCodeModel:               code.NewRedeemCodeModel(mysqlConnAbility),                    // 兑换码表
		DictWhiteModel:                systemModel.NewDictWhiteModel(mysqlConnAbility),              // 字典白名单表

		// muse-hub 数据库
		HubLessonModel:   hub.NewLessonSnapshotModel(mysqlConnHub),
		HubWritePPTModel: hub.NewWritePptSnapshotModel(mysqlConnHub),
		HubLiveModel:     hub.NewLiveSnapshotModel(mysqlConnHub),
		// HubLessonQuestionTtsModel: knowledgeModel.NewLessonQuestionTtsModel(mysqlConnHub),
		HubSuperTrainModel: hub.NewSuperTrainSnapshotModel(mysqlConnHub),

		// king-club 数据库
		TrainCourseModel:                  supertrain.NewCourseModel(mysqlConnKingClub),
		TrainCourseChapterModel:           supertrain.NewCourseChapterModel(mysqlConnKingClub),
		TrainChapterTaskModel:             supertrain.NewChapterTaskModel(mysqlConnKingClub),
		TrainSubTaskModel:                 supertrain.NewSubTaskModel(mysqlConnKingClub),
		TrainMedalModel:                   supertrain.NewMedalModel(mysqlConnKingClub),
		TrainSubTaskMedalModel:            supertrain.NewSubTaskMedalModel(mysqlConnKingClub),
		SupeScheduleModel:                 supervisor.NewScheduleModel(mysqlConnKingClub),
		SupeStreamModel:                   supervisor.NewStreamModel(mysqlConnKingClub),
		SupeScheduleInteractRelationModel: supervisor.NewScheduleInteractionLinkModel(mysqlConnKingClub),
		AcTeacherModel:                    supervisor.NewTeacherModel(mysqlConnKingClub),
		InteractionModel:                  supervisor.NewInteractionModel(mysqlConnKingClub),

		LiveUserModel: live_course.NewUserModel(mysqlConnKingClub),
		ChatRoomModel: live_course.NewChatRoomModel(mysqlConnKingClub),
		// RPC服务
		AbilityRPC:  ability.NewClient(c.RpcTarget.Ability),
		PassportRPC: passport.NewClient(c.RpcTarget.Passport),
		CoreRPC:     core.NewClient(c.RpcTarget.Core),
	}
}

// 初始化雪花算法发号器
func initSnowflake() *snowflake.Node {
	// 获取hostname
	hostname, err := os.Hostname()
	if err != nil {
		logc.Errorf(context.Background(), "获取hostname失败：%s", err)
		panic(fmt.Sprintf("获取hostname失败：%s", err))
	}

	// 机器名称字符转ASCII
	numStr := strings.Builder{}
	for _, c := range hostname {
		numStr.WriteString(cast.ToString(int(c)))
	}

	// 机器ASCII编号 + 随机数
	randomNumber := (cast.ToInt64(numStr.String()) + rand.Int63n(1024)) % 1024

	// 初始化雪花算法
	node, err := snowflake.NewNode(randomNumber)
	if err != nil {
		logc.Errorf(context.Background(), "初始化雪花算法失败：%s", err)
		panic(fmt.Sprintf("初始化雪花算法失败：%s", err))
	}
	return node
}

// InitTest 单元测试初始化依赖注册
func InitTest() (context.Context, *ServiceContext, logx.Logger) {
	var c config.Config
	var configFile = "../../../../../etc/application.yaml"

	conf.MustLoad(configFile, &c)
	ctx := context.Background()
	svcCtx := NewServiceContext(c)
	log := logx.WithContext(ctx)
	return ctx, svcCtx, log
}
