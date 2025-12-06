package lesson

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonExample struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLessonExample(ctx context.Context, svcCtx *svc.ServiceContext) *LessonExample {
	return &LessonExample{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LessonExample) LessonExample(req *types.LessonExampleReq) (resp *types.LessonExampleResp, err error) {
	data, err := l.svcCtx.ExampleModel.FindOne(l.ctx, req.ExampleId)
	if err != nil {
		logc.Errorf(l.ctx, "查询例题json失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	// 暂无数据
	if data == nil {
		return
	}

	return &types.LessonExampleResp{
		Content: data.Explain,
	}, nil
}
