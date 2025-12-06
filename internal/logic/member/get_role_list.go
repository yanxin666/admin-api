package member

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/define"
	"muse-admin/internal/model/member"
	"muse-admin/pkg/errs"
	"time"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleList(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleList {
	return &GetRoleList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleList) GetRoleList(req *types.GetRoleListReq) (resp *types.GetRoleListResp, err error) {
	resp = &types.GetRoleListResp{
		List: make([]types.MemberSubInfo, 0),
	}

	var (
		list          []types.MemberSubInfo
		userIds       []any
		douShenVipMap map[int64]member.DoushenVip
	)

	// 手机号加密，判断手机号用户是否存在
	encryptPhone, err := util.AesEncrypt(l.ctx, req.Phone, l.svcCtx.Config.EncryptKey)
	if err != nil {
		return nil, err
	}
	// 只获取产品线对应的用户
	baseUserInfo, err := l.svcCtx.BaseUserModel.FindOneByMaskPhoneProduct(l.ctx, encryptPhone, req.Product)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}
	if baseUserInfo == nil {
		return nil, errs.NewMsg(errs.UserIdErrorCode, "输入的手机号暂无用户").ShowMsg()
	}

	// 获取用户角色列表
	userList, err := l.svcCtx.UserModel.FindOneByBaseUserId(l.ctx, baseUserInfo.Id)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}
	if userList == nil {
		return nil, errs.NewMsg(errs.ServerErrorCode, "角色不存在").ShowMsg()
	}

	userIds, _ = util.ArrayColumnByPtr(userList, "Id")
	douShenVipMap, err = l.svcCtx.DouShenVipModel.BatchMapByUserIds(l.ctx, userIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	for _, u := range userList {
		maxTime := time.Time{} // 初始时间为零时间
		if u.UpdatedAt.After(maxTime) {
			maxTime = u.UpdatedAt
		}
		var (
			isVipValid bool
			endTime    int64
		)
		if vip, ok := douShenVipMap[u.Id]; ok {
			endTime = vip.EndTime.Unix()
			isVipValid = vip.EndTime.Unix() > time.Now().Unix()
		}

		user := types.MemberSubInfo{
			BaseUserId:     u.BaseUserId,                                              // 主用户ID
			UserId:         u.Id,                                                      // 用户ID
			RoleType:       u.RoleType,                                                // 角色类型 1.孩子[默认] 2.家长
			Nickname:       u.Nickname,                                                // 昵称
			RealName:       u.RealName,                                                // 姓名
			Avatar:         u.Avatar,                                                  // 头像
			Gender:         u.Gender,                                                  // 性别 0.未知 1.男 2.女
			Grade:          u.Grade,                                                   // 年级 10.学龄前 11.一年级 12.二年级 13.三年级 14.四年级 15.五年级 16.六年级 20.中小衔接 21.初一 22.初二 23.初三 31.高一 32.高二 33.高三 41.大学语文 42.中文专业 43.文学专家
			Birthday:       util.UnixToDate(u.Birthday.Time.Unix()),                   // 出生日期
			Region:         util.ConcatString(u.ProvinceName, u.CityName, u.AreaName), // 地区
			IsInfoFinished: u.IsInfoFinished,                                          // 是否完善信息 1.未完善 2.已完善
			VipName:        define.VipTypeMap[douShenVipMap[u.Id].VipType],            // 会员名称
			VipType:        douShenVipMap[u.Id].VipType,                               // 会员类型
			IsVipValid:     isVipValid,                                                // 会员是否过期
			VipEndTime:     endTime,                                                   // 会员过期时间
			Status:         u.Status,                                                  // 状态
			Remark:         u.Remark,                                                  // 备注
			CreatedAt:      u.CreatedAt.Unix(),                                        // 创建时间
			UpdatedAt:      u.UpdatedAt.Unix(),                                        // 最近登录时间
		}

		list = append(list, user)
	}

	return &types.GetRoleListResp{
		List: list,
	}, nil
}
