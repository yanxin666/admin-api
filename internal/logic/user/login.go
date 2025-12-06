package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/spf13/cast"
	"muse-admin/internal/define"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq, r *http.Request) (resp *types.LoginResp, err error) {
	verifyCode, _ := l.svcCtx.Redis.Get(define.SysLoginCaptchaCachePrefix + req.CaptchaId)
	if verifyCode != req.VerifyCode {
		return nil, errs.NewCode(errs.CaptchaErrorCode)
	}

	sysUser, err := l.svcCtx.SysUserModel.FindOneByAccount(l.ctx, req.Account)
	if err != nil {
		return nil, errs.NewCode(errs.AccountErrorCode)
	}

	if sysUser.Password != util.GenerateMD5Str(req.Password+l.svcCtx.Config.Salt) {
		return nil, errs.NewCode(errs.PasswordErrorCode)
	}

	if sysUser.Status != define.SysEnable {
		return nil, errs.NewCode(errs.AccountDisableErrorCode)
	}

	if sysUser.Id != define.SysSuperUserId {
		dept, _ := l.svcCtx.SysDeptModel.FindOne(l.ctx, sysUser.DeptId)
		if dept.Status == define.SysDisable {
			return nil, errs.NewCode(errs.AccountDisableErrorCode)
		}
	}

	token, _ := l.getJwtToken(sysUser.Id)
	_, err = l.svcCtx.Redis.Del(req.CaptchaId)

	loginLog := workbench.Log{
		UserId: sysUser.Id,
		Ip:     r.Header.Get("X-Forwarded-For"),
		Uri:    r.RequestURI,
		Type:   1,
		Status: 1,
	}
	_, err = l.svcCtx.SysLogModel.Insert(l.ctx, &loginLog)

	err = l.svcCtx.Redis.Setex(define.SysOnlineUserCachePrefix+cast.ToString(sysUser.Id), token, int(l.svcCtx.Config.JwtAuth.AccessExpire))
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	userInfoMap := map[string]string{
		"account":  sysUser.Account,
		"username": sysUser.Username,
		"mobile":   sysUser.Mobile,
	}
	// 将用户信息缓存
	_ = l.svcCtx.RedisClient.HMSet(l.ctx, define.SysOnlineUserInfoPrefix+cast.ToString(sysUser.Id), userInfoMap, int(l.svcCtx.Config.JwtAuth.AccessExpire))

	return &types.LoginResp{
		Token: token,
	}, nil
}

func (l *LoginLogic) getJwtToken(userId int64) (string, error) {
	iat := time.Now().Unix()
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + l.svcCtx.Config.JwtAuth.AccessExpire
	claims["iat"] = iat
	claims[define.SysJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(l.svcCtx.Config.JwtAuth.AccessSecret))
}
