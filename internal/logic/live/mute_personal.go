package live

import (
	"context"
	"muse-admin/internal/tools"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/proto/core"
	"github.com/zeromicro/go-zero/core/logx"
)

type MutePersonal struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMutePersonal(ctx context.Context, svcCtx *svc.ServiceContext) *MutePersonal {
	return &MutePersonal{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MutePersonal) MutePersonal(req *types.MutePersonalReq) error {
	if err := l.checkParams(req); err != nil {
		return err
	}
	userId := tools.GetUserIdByCtx(l.ctx)
	_, err := l.svcCtx.CoreRPC.LiveClient.MutePersonal(l.ctx, &core.MutePersonalReq{
		RoomId:         req.RoomId,
		MuteStatus:     req.IsMute,
		TargetImUserId: req.ImUserId,
		OperateUserId:  userId,
		Duration:       req.Duration,
	})
	if err != nil {
		logz.Errorf(l.ctx, "全体禁言失败 err ：%s", err.Error())
		return err
	}
	return nil
}

// 校验参数
func (l *MutePersonal) checkParams(req *types.MutePersonalReq) error {
	// todo: add your logic check params and delete this line
	// 若校验的参数过少，可删除此方法，直接在上方 MutePersonal 中编写

	return nil
}
