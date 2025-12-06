package member

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/define"
	"muse-admin/internal/model/user"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	ctxt "muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
)

type LogOff struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const (
	RedisKeyJwt         = "muse-passport:user:token:%d"        // jwt授权访问令牌缓存key
	RedisKeyJwtRefresh  = "muse-passport:user:tokenRefresh:%d" // jwt授权刷新令牌缓存key
	RedisKeyUserInfoKey = "muse-passport:user:userInfo:%d"     // 用户缓存
)

var desc = map[int64]string{
	1: "恢复用户",
	2: "注销删除",
}

func NewLogOff(ctx context.Context, svcCtx *svc.ServiceContext) *LogOff {
	return &LogOff{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogOff) LogOff(req *types.MemberLogOffReq) (resp *types.MemberLogOffResp, err error) {
	data, err := l.svcCtx.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		return nil, errs.NewCode(errs.UserIdErrorCode)
	}

	var status int64
	switch data.Status {
	case define.UserStatus.Normal:
		status = define.UserStatus.LogOff // 注销
	case define.UserStatus.LogOff:
		status = define.UserStatus.Normal // 恢复
	default:
		return nil, errs.NewCode(errs.ServerErrorCode)
	}

	_, err = l.svcCtx.UserModel.UpdateUserById(l.ctx, req.UserId, &user.User{
		Status: status,
	})
	if err != nil {
		return nil, errs.NewCode(errs.ServerErrorCode)
	}

	// 删除Token
	_ = l.DeleteToken(l.ctx, data.Id)
	_ = l.DeleteRefreshToken(l.ctx, data.Id)

	// 删除用户信息缓存
	_ = l.DeleteUserInfo(l.ctx, data.Id)

	// 注销成功日志
	l.InsertLoginLog(l.ctx, data.Id, ctxt.GetUserIdByCtx(l.ctx), status)

	return &types.MemberLogOffResp{
		Result: true,
	}, nil
}

// DeleteToken 删除token缓存
func (l *LogOff) DeleteToken(ctx context.Context, userId int64) error {
	key := fmt.Sprintf(RedisKeyJwt, userId)
	count, err := l.svcCtx.RedisClient.Del(ctx, key)
	if err != nil {
		logz.Warnf(ctx, "删除token缓存失败，key:%s, err:%v", key, err)
		return errs.WithMsg(err, errs.ErrCodeAbnormal, "删除token缓存失败")
	}

	if count <= 0 {
		return errs.NewMsg(errs.ErrCodeAbnormal, "删除token缓存失败")
	}

	return nil
}

// DeleteRefreshToken 删除refresh_token缓存
func (l *LogOff) DeleteRefreshToken(ctx context.Context, userId int64) error {
	key := fmt.Sprintf(RedisKeyJwtRefresh, userId)
	count, err := l.svcCtx.RedisClient.Del(ctx, key)
	if err != nil {
		logz.Warnf(ctx, "删除token缓存失败，key:%s, err:%v", key, err)
		return errs.WithMsg(err, errs.ErrCodeAbnormal, "删除refresh_token缓存失败")
	}

	if count <= 0 {
		return errs.NewMsg(errs.ErrCodeAbnormal, "删除refresh_token缓存失败")
	}

	return nil
}

// DeleteUserInfo 删除用户信息缓存
func (l *LogOff) DeleteUserInfo(ctx context.Context, userId int64) error {
	var keyArr []string

	key := fmt.Sprintf(RedisKeyUserInfoKey, userId)
	keys, err := l.svcCtx.RedisClient.HGetAll(ctx, key)
	if err != nil {
		return err
	}

	// 遍历所有 key，并删除哈希表
	for k, _ := range keys {
		keyArr = append(keyArr, k)
	}

	_, err = l.svcCtx.RedisClient.HDel(ctx, key, keyArr...)
	if err != nil {
		return err
	}

	return nil
}

// InsertLoginLog 记录注销相关日志
func (l *LogOff) InsertLoginLog(ctx context.Context, userId, operateId, status int64) {
	loginLog := workbench.Log{
		UserId:  operateId,
		Ip:      ctxt.GetRequestLogsByCtx(ctx).ClientIP,
		Uri:     ctxt.GetRequestLogsByCtx(ctx).URL,
		Type:    2,
		Status:  1,
		Request: fmt.Sprintf("用户ID：%d, %s", userId, desc[status]),
	}
	_, _ = l.svcCtx.SysLogModel.Insert(l.ctx, &loginLog)
}
