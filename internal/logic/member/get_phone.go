package member

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPhone struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPhone(ctx context.Context, svcCtx *svc.ServiceContext) *GetPhone {
	return &GetPhone{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPhone) GetPhone(req *types.MemberGetPhoneReq) (resp *types.MemberGetPhoneResp, err error) {
	resp = &types.MemberGetPhoneResp{
		Phone: 0,
	}

	baseInfo, err := l.svcCtx.BaseUserModel.FindOne(l.ctx, cast.ToInt64(req.BaseUserId))
	if err != nil && !errors.Is(err, sqlc.ErrNotFound) {
		return nil, errs.NewCode(errs.ServerErrorCode)
	}
	if baseInfo == nil {
		return resp, nil
	}

	resp.Phone, err = util.AesDecryptIdToInt(l.ctx, baseInfo.MaskPhone, l.svcCtx.Config.EncryptKey)
	if err != nil {
		return nil, errs.NewCode(errs.ServerErrorCode)
	}

	return resp, nil
}
