package order

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/tencent"
	"e.coding.net/zmexing/nenglitanzhen/proto/ability"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"io"
	"muse-admin/internal/define"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportWaitSend struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportWaitSend(ctx context.Context, svcCtx *svc.ServiceContext) *ExportWaitSend {
	return &ExportWaitSend{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportWaitSend) ExportWaitSend(req *types.OrderWaitSendReq) (resp *types.ExportWaitSendResp, err error) {
	param := &ability.OrderWaitSendReq{
		OrderNo:     req.OrderNo,
		SubOrderNo:  req.SubOrderNo,
		GoodsName:   req.GoodsName,
		UserName:    req.UserName,
		UserPhone:   req.UserPhone,
		GrantStatus: req.GrantStatus,
	}
	if req.ExportAll == 2 {
		if req.Page > 0 && req.Limit > 0 {
			param.Page = &ability.Page{
				Page: req.Page,
				Size: req.Limit,
			}
		}
	}
	timeFilter := &ability.TimeFilter{}
	if req.StartTime != "" {
		timeFilter.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		timeFilter.EndTime = req.EndTime
	}
	param.TimeFilter = timeFilter

	data, err := l.svcCtx.AbilityRPC.OrderClient.OrderWaitSend(l.ctx, param)
	if err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	reader, err := l.buildIO(data.List)
	if err != nil {
		return nil, err
	}

	url, err := l.uploadToCos(reader)
	if err != nil {
		logz.Errorf(l.ctx, "[待发货导出]upload to cos failed, err:%v", err)
		return nil, err
	}

	return &types.ExportWaitSendResp{
		Url: url,
	}, nil
}

func (l *ExportWaitSend) buildIO(data []*ability.OrderWaitSend) (io.Reader, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 创建Sheet
	sheetName := "sheet1"
	defaultSheet, err := f.NewSheet(sheetName)
	if err != nil {
		logz.Errorf(l.ctx, "[待发货导出]create new sheet failed, err:%v", err)
		return nil, err
	}
	// 设置默认Sheet
	f.SetActiveSheet(defaultSheet)
	// 设置表头
	titles := []string{
		"订单号",
		"子订单号",
		"商品名称",
		"数量",
		"收货人",
		"支付渠道",
		"收货人电话",
		"省",
		"市",
		"区",
		"详细地址",
		"发货状态",
		"创建时间",
	}
	err = f.SetSheetRow(sheetName, "A1", &titles)
	if err != nil {
		logz.Errorf(l.ctx, "[待发货导出]set sheet row failed, err:%v", err)
		return nil, err
	}
	// 设置所有数据
	for i := 0; i < len(data); i++ {
		row := []interface{}{
			data[i].OrderNo,
			data[i].SubOrderNo,
			data[i].GoodsName,
			data[i].Pcs,
			data[i].UserName,
			define.PayChannelStr[data[i].PayChannel],
			data[i].UserPhone,
			data[i].Province,
			data[i].City,
			data[i].Area,
			data[i].Address,
			define.GrantStatusStr[data[i].GrantStatus],
			data[i].CreatedAt,
		}
		err = f.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+2), &row)
		if err != nil {
			logz.Errorf(l.ctx, "[待发货导出]set sheet row failed, err:%v", err)
			return nil, err
		}
	}
	// 写入Buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		logz.Errorf(l.ctx, "[待发货导出]write to buffer failed, err:%v", err)
	}
	return strings.NewReader(buf.String()), nil
}

func (l *ExportWaitSend) uploadToCos(r io.Reader) (string, error) {
	cos := tencent.NewCos(l.ctx, tencent.CocConf{
		SecretId:  l.svcCtx.Config.Oss.SecretId,
		SecretKey: l.svcCtx.Config.Oss.SecretKey,
		Appid:     l.svcCtx.Config.Oss.Appid,
		Bucket:    l.svcCtx.Config.Oss.Bucket,
		Region:    l.svcCtx.Config.Oss.Region,
	})
	fileName := fmt.Sprintf("order/发货数据_%s.xlsx", time.Now().Format("20060102150405"))
	return cos.Put(l.ctx, fileName, r, 300)
}
