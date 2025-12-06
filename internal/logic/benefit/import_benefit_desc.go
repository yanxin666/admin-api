package benefit

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportBenefitDesc struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportBenefitDesc(ctx context.Context, svcCtx *svc.ServiceContext) *ImportBenefitDesc {
	return &ImportBenefitDesc{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportBenefitDesc) ImportBenefitDesc(req *types.ImportBenefitDescReq) (resp *types.ImportBenefitDescResp, err error) {
	resp = &types.ImportBenefitDescResp{
		List: make([]types.UserBenefitInfo, 0),
	}

	// 获取子用户权益列表
	subList, err := l.svcCtx.ChannelSubOrderModel.FindByChannelOrderId(l.ctx, req.ChannelOrderId)
	if err != nil {
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	for _, v := range subList {
		re, err := l.svcCtx.BenefitResourceModel.FindOne(l.ctx, v.BenefitsResourceId)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		// data, err := l.svcCtx.UserBenefitModel.FindOneByUserIdBenefitsResourceId(l.ctx, v.UserId, v.BenefitsResourceId)
		// if err != nil {
		//	return nil, errs.WithCode(err, errs.ServerErrorCode)
		// }
		baseInfo, err := l.svcCtx.BaseUserModel.FindOne(l.ctx, v.BaseUserId)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		// 解密手机
		unMaskPhone, err := util.AesDecrypt(l.ctx, baseInfo.MaskPhone, l.svcCtx.Config.EncryptKey)
		if err != nil {
			return nil, err
		}
		userInfo, err := l.svcCtx.UserModel.FindOne(l.ctx, v.UserId)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		resp.List = append(resp.List, types.UserBenefitInfo{
			Id:          v.Id,
			BaseUserId:  v.BaseUserId,
			Phone:       util.GenerateMaskPhone(unMaskPhone),
			UserId:      v.UserId,
			RealName:    userInfo.RealName,
			Status:      v.Status,
			BenefitName: re.Name,
			StartTime:   v.CreatedAt.Unix(),
			EndTime:     v.UpdatedAt.Unix(),
		})
	}

	return
}
