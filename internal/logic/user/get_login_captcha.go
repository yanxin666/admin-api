package user

import (
	"context"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/other"

	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/logx"
	utils2 "github.com/zeromicro/go-zero/core/utils"
)

type GetLoginCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoginCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginCaptchaLogic {
	return &GetLoginCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoginCaptchaLogic) GetLoginCaptcha() (resp *types.LoginCaptchaResp, err error) {
	var store = base64Captcha.DefaultMemStore
	captcha := other.NewCaptcha(45, 80, 4, 40, 30, 89, 0)
	driver := captcha.DriverString()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := c.Generate()
	val := store.Get(id, true)
	captchaId := utils2.NewUuid()
	err = l.svcCtx.Redis.Setex(define.SysLoginCaptchaCachePrefix+captchaId, val, 300)

	return &types.LoginCaptchaResp{
		CaptchaId:  captchaId,
		VerifyCode: b64s,
	}, nil
}
