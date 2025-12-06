package code

import (
	"context"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"muse-admin/internal/svc"
)

const (
	charset    = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 兑换码格式
	codeLength = 16                                 // 兑换码长度
	maxRetries = 3                                  // 插入失败后的重试次数
	ipLimit    = 5                                  // 每个IP的最大请求次数
	windowTime = 60                                 // 限制时间窗口

	batchSize = 100 // 每次插入100条
)

type GenerateCode struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateCode(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCode {
	return &GenerateCode{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// BatchRedemptionCodes 批量兑换码生成
func (g *GenerateCode) BatchRedemptionCodes(ctx context.Context, codeNumber int) ([]string, error) {
	var codeMap []string
	for i := 0; i < codeNumber; i++ {
		id, err := gonanoid.Generate(charset, codeLength) // 生成兑换码
		if err != nil || id == "" {
			logc.Errorf(ctx, "兑换码生成异常：%s", err)
			continue
		}
		codeMap = append(codeMap, id)
	}
	return codeMap, nil
}
