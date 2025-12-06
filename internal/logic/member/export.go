package member

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"muse-admin/internal/define"
	"muse-admin/internal/model/member"
	"muse-admin/internal/model/tools"
	"muse-admin/internal/model/user"
	ctxt "muse-admin/internal/tools"
	"muse-admin/pkg/errs"
	"muse-admin/pkg/excel"
	"time"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Export struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExport(ctx context.Context, svcCtx *svc.ServiceContext) *Export {
	return &Export{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Export) Export(req *types.ExportMemberReq) (resp *types.ExportMemberResp, err error) {
	var (
		keyword   string
		userInfo  *user.User
		userInfos map[int64][]user.User

		userInfoArr   []user.User
		userIds       []any
		douShenVipMap map[int64]member.DoushenVip
	)

	// 构建查询条件
	condition := tools.FilterConditions(req)
	delete(condition, "phone")
	delete(condition, "user_id")
	delete(condition, "export_type")
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
			return nil, errors.New("暂无数据，无法导出")
		}
		condition["id"] = userInfo.BaseUserId
	}

	// 获取主用户列表
	list, _, err := l.svcCtx.BaseUserModel.FindPageByCondition(l.ctx, req.Page, req.Limit, keyword, condition)
	if err != nil {
		logc.Errorf(l.ctx, "批量查询用户失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 空数据
	if len(list) == 0 {
		return nil, errors.New("暂无数据，无法导出")
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

	var (
		userArr []types.ExportMemberInfo // 导出用户列表
		remark  string                   // 备注
	)
	for _, v := range list {
		maxTime := time.Time{} // 初始时间为零时间

		for _, u := range userInfos[v.Id] {
			if u.UpdatedAt.After(maxTime) {
				maxTime = u.UpdatedAt
			}

			var (
				vipName, isVipValid, vipValidTime = "-", "-", "-" // 会员相关
				isInfoFinished                    = "未完善"         // 是否完善信息
			)
			if vip, ok := douShenVipMap[u.Id]; ok {
				vipName = util.GetValueOrDefaultByMap(define.VipTypeMap, douShenVipMap[u.Id].VipType)
				vipValidTime = util.UnixToDateTime(vip.EndTime.Unix())
				if vip.EndTime.Unix() > time.Now().Unix() {
					isVipValid = "会员有效"
				} else {
					isVipValid = "会员过期"
				}
			}
			if u.IsInfoFinished == define.IsInfoFinishedStatus.Finish {
				isInfoFinished = "已完善"
			}

			// 优先取备注里的值
			if u.Remark != "" {
				remark = u.Remark
			} else if maxTime.Unix() == v.UpdatedAt.Unix() {
				remark = "最近登录"
			}

			// 解密手机
			unMaskPhone, err := util.AesDecrypt(l.ctx, v.MaskPhone, l.svcCtx.Config.EncryptKey)
			if err != nil {
				return nil, err
			}

			// 用户列表
			userArr = append(userArr, types.ExportMemberInfo{
				BaseUserId:     v.Id,                                                            // 主用户ID
				UserNo:         v.UserNo,                                                        // 用户编号
				Phone:          unMaskPhone,                                                     // 明文手机号
				Product:        util.GetValueOrDefaultByMap(define.UserProductMap, v.Product),   // 产品线
				Source:         util.GetValueOrDefaultByMap(define.UserSourceMap, v.Source),     // 来源
				Status:         util.GetValueOrDefaultByMap(define.BaseUserStatusMap, v.Status), // 状态
				UserId:         u.Id,                                                            // 用户ID
				RoleType:       define.GetUserRole(u.RoleType),                                  // 角色类型 1.孩子[默认] 2.家长
				Nickname:       u.Nickname,                                                      // 昵称
				RealName:       u.RealName,                                                      // 姓名
				Gender:         define.GetUserGenderInfo(u.Gender),                              // 性别 0.未知 1.男 2.女
				Grade:          define.GetUserGradeInfo(u.Grade),                                // 年级 10.学龄前 11.一年级 12.二年级 13.三年级 14.四年级 15.五年级 16.六年级 20.中小衔接 21.初一 22.初二 23.初三 31.高一 32.高二 33.高三 41.大学语文 42.中文专业 43.文学专家
				Birthday:       util.UnixToDate(u.Birthday.Time.Unix()),                         // 出生日期
				Region:         util.ConcatString(u.ProvinceName, u.CityName, u.AreaName),       // 地区
				IsInfoFinished: isInfoFinished,                                                  // 是否完善信息 1.未完善 2.已完善
				VipName:        vipName,                                                         // 会员名称
				IsVipValid:     isVipValid,                                                      // 会员是否过期
				VipEndTime:     vipValidTime,                                                    // 会员过期时间
				SubUserStatus:  util.GetValueOrDefaultByMap(define.UserStatusMap, u.Status),     // 子用户状态
				Remark:         remark,                                                          // 备注
				CreatedAt:      util.UnixToDateTime(u.CreatedAt.Unix()),                         // 创建时间
				UpdatedAt:      util.UnixToDateTime(u.UpdatedAt.Unix()),                         // 最近登录时间
			})
		}
	}

	// 设置表头的映射
	headerMap := map[string]string{
		"BaseUserId":     "主用户ID",
		"UserNo":         "用户编号",
		"Phone":          "手机号",
		"Product":        "产品线",
		"Source":         "首次来源",
		"Status":         "账号状态",
		"UserId":         "用户ID",
		"RoleType":       "角色",
		"Nickname":       "昵称",
		"RealName":       "真实姓名",
		"Gender":         "性别",
		"Grade":          "年级",
		"Birthday":       "生日",
		"Region":         "地区",
		"IsInfoFinished": "是否完善信息",
		"VipName":        "会员类型",
		"VipEndTime":     "会员过期时间",
		"IsVipValid":     "会员是否过期",
		"SubUserStatus":  "角色状态",
		"Remark":         "备注",
		"CreatedAt":      "注册时间",
		"UpdatedAt":      "最近登录时间",
	}

	o := excel.NewExCelExport(
		l.svcCtx.Config.Oss,
		excel.SetFilepath("member"),
		excel.SetFilename(util.GetNowTimeNoFormat()+"-"+ctxt.GetUserInfoByCtx(l.ctx).Username),
	)

	url, err := o.GenExcelFile(l.ctx, headerMap, userArr)
	if err != nil {
		logz.Errorf(l.ctx, "[用户列表导出]upload to cos failed, err:%v", err)
		return nil, err
	}

	return &types.ExportMemberResp{
		Url: url,
	}, nil
}
