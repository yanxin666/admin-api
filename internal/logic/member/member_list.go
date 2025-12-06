package member

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/base64"
	"errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"muse-admin/internal/define"
	"muse-admin/internal/model/member"
	"muse-admin/internal/model/tools"
	"muse-admin/internal/model/user"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"time"
)

var GatewayLogTopicID = map[string]string{
	"dev":  "e18ac0d0-264b-4bca-96b3-f44b4d9ec891",
	"test": "9e1702a1-3c35-4888-a55e-b03621dba1ac",
	"pro":  "70944513-db91-4d1c-b727-0af9ef9c16ea",
}

type MemberList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMemberList(ctx context.Context, svcCtx *svc.ServiceContext) *MemberList {
	return &MemberList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MemberList) MemberList(req *types.MemberListReq) (resp *types.MemberListResp, err error) {
	var (
		keyword   string
		userInfo  *user.User
		userInfos map[int64][]user.User

		userInfoArr   []user.User
		userIds       []any
		douShenVipMap map[int64]member.DoushenVip
	)
	resp = &types.MemberListResp{
		List:       make([]types.MemberInfo, 0),
		Pagination: types.Pagination{},
	}

	// 构建查询条件
	condition := tools.FilterConditions(req)
	delete(condition, "phone")
	delete(condition, "user_id")
	if len(req.Phone) == 11 {
		// 加密手机号
		condition["mask_phone"], err = util.AesEncrypt(l.ctx, req.Phone, l.svcCtx.Config.EncryptKey)
		if err != nil {
			return nil, err
		}
	} else {
		keyword = req.Phone
	}

	if req.UserId != "" {
		userInfo, err = l.svcCtx.UserModel.FindOne(l.ctx, cast.ToInt64(req.UserId))
		if err != nil && !errors.Is(err, sqlc.ErrNotFound) {
			return nil, errs.NewCode(errs.ServerErrorCode)
		}
		if userInfo == nil {
			return resp, nil
		}
		condition["id"] = userInfo.BaseUserId
	}

	// 获取主用户列表
	list, total, err := l.svcCtx.BaseUserModel.FindPageByCondition(l.ctx, req.Page, req.Limit, keyword, condition)
	if err != nil {
		logc.Errorf(l.ctx, "批量查询用户失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 空数据
	if len(list) == 0 {
		return resp, nil
	}

	// 获取所有子用户
	baseUserIds, _ := util.ArrayColumn(list, "Id")
	userInfos, err = l.svcCtx.UserModel.BatchMapByBaseUserIds(l.ctx, baseUserIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 获取所有子用户ID所对应的会员权益
	userInfoArr, err = l.svcCtx.UserModel.BatchListByBaseUserIds(l.ctx, baseUserIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}
	userIds, _ = util.ArrayColumn(userInfoArr, "Id")
	douShenVipMap, err = l.svcCtx.DouShenVipModel.BatchMapByUserIds(l.ctx, userIds)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	for _, v := range list {
		var userArr []types.MemberSubInfo // 子用户列表
		maxTime := time.Time{}            // 初始时间为零时间

		for _, u := range userInfos[v.Id] {
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

			// 腾讯云访问日志的内容base64编码
			original := "content:\"userLinkRoute\" and userId:\" " + cast.ToString(u.Id) + "\""
			encoded := base64.StdEncoding.EncodeToString([]byte(original))
			logURL := "https://console.cloud.tencent.com/cls/search?region=ap-beijing&topic_id=" + GatewayLogTopicID[l.svcCtx.Config.Mode] + "&queryBase64=" + encoded + "&time=now-7d,now"

			subUser := types.MemberSubInfo{
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
				TencentLogURL:  logURL,                                                    // 腾讯云日志URL
			}

			userArr = append(userArr, subUser)
		}

		// 主用户列表
		resp.List = append(resp.List, types.MemberInfo{
			BaseUserId: v.Id,      // 主用户ID
			UserNo:     v.UserNo,  // 用户编号
			Phone:      v.Phone,   // 手机号
			Product:    v.Product, // 产品线
			Source:     v.Source,  // 来源
			Status:     v.Status,  // 状态
			Remark:     v.Remark,  // 备注
			SubList:    userArr,
			CreatedAt:  v.CreatedAt.Unix(),
			UpdatedAt:  v.UpdatedAt.Unix(),
			RecentlyAt: maxTime.Unix(), // 子用户最近一次登录时间
		})
	}

	// 分页信息
	resp.Pagination = types.Pagination{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	return
}
