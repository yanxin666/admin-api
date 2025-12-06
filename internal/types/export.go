package types

type ExportMemberInfo struct {
	BaseUserId int64  `json:"baseUserId"` // 主用户ID
	UserNo     string `json:"userNo"`     // 用户编号
	Phone      string `json:"phone"`      // 手机号
	Product    string `json:"product"`    // 产品线 0.辞源 1.AI搜索 2.白帝辞 3.智普乌托邦(1001~1002) 4.豆子机器人
	Source     string `json:"source"`     // 首次来源 100.提审标记 0.APP 1.豆伴匠 2.风的颜色 1001.听力熊 1002.微软
	Status     string `json:"status"`     // 账号状态 0.正常 1.冻结 2.注销

	UserId         int64  `json:"userId"`         // 用户ID
	RoleType       string `json:"role"`           // 角色类型 1.孩子[默认] 2.家长
	Nickname       string `json:"nickname"`       // 昵称
	RealName       string `json:"realName"`       // 姓名
	Gender         string `json:"gender"`         // 性别 0.未知 1.男 2.女
	Grade          string `json:"grade"`          // 年级 10.学龄前 11.一年级 12.二年级 13.三年级 14.四年级 15.五年级 16.六年级 20.中小衔接 21.初一 22.初二 23.初三 31.高一 32.高二 33.高三 41.大学语文 42.中文专业 43.文学专家
	Birthday       string `json:"birthday"`       // 生日[年月日]
	Region         string `json:"region"`         // 地区
	IsInfoFinished string `json:"isInfoFinished"` // 是否完善信息 1.未完善 2.已完善
	VipName        string `json:"vipName"`        // 会员名称 1.豆神会员（基础版） 2.豆神会员（典藏版） 3.豆神超级会员 4.豆神会员（特别版）
	IsVipValid     string `json:"isVipValid"`     // 会员是否到期
	VipEndTime     string `json:"vipEndTime"`     // 会员过期时间
	SubUserStatus  string `json:"subUserStatus"`  // 状态 0.正常 1.冻结 2.注销
	Remark         string `json:"remark"`         // 备注
	CreatedAt      string `json:"createdAt"`      // 创建时间
	UpdatedAt      string `json:"updatedAt"`      // 更新时间
}
