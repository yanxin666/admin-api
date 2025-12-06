package public

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/telemetry"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"go.opentelemetry.io/otel/codes"
	taskModel "muse-admin/internal/model/task"
	"muse-admin/internal/svc"
	_import "muse-admin/internal/svc/import"
	"muse-admin/internal/svc/public"
	ctxt "muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
)

type ImportFile struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	alarmService *public.BizAlarm
}

func NewImportFile(ctx context.Context, svcCtx *svc.ServiceContext) *ImportFile {
	return &ImportFile{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		alarmService: public.NewAlarmService(svcCtx),
	}
}

func (l *ImportFile) ImportFile(req *types.ImportFileReq) (resp *types.Result, err error) {
	userId := ctxt.GetUserIdByCtx(l.ctx)

	insertData := &taskModel.SyncTask{
		OperateId:     userId,
		Type:          req.Type,
		FileId:        req.FileId,
		FileName:      req.Filename,
		FileSheet:     req.FileSheet,
		FileSheetName: req.FileSheetName,
	}

	res, err := l.svcCtx.SyncTaskModel.Insert(l.ctx, insertData)
	if err != nil {
		return nil, err
	}
	taskId, err := res.LastInsertId()
	if err != nil {
		return nil, errs.WithMsg(err, errs.ErrCodeAbnormal, fmt.Sprintf("未获取到添加的任务Id"))
	}

	// // 推送任务ID到MQ
	// err = l.svcCtx.MQProducer.SendSync(l.ctx, define.ProducerTopics[l.svcCtx.Config.Mode].ExcelImport, map[string]any{
	//	"id": taskId,
	// })
	// if err != nil {
	//	return nil, errs.WithMsg(err, errs.ErrImportPushMQ, fmt.Sprintf("推送任务失败"))
	// }

	// 创建新的 context，避免继承原始请求的取消信号
	newCtx, span, end := telemetry.StartSpanInGoroutine(l.ctx, "syncExecute", nil)
	// 提取线索
	threading.GoSafeCtx(newCtx, func() {
		defer end()
		err = l.syncExecute(newCtx, taskId)
		if err != nil {
			logz.Errorf(newCtx, "[文件导入] syncExecute 失败 error: %v", err)
			span.SetStatus(codes.Error, err.Error())
			return
		}
	})

	return &types.Result{Result: true}, nil
}

func (l *ImportFile) syncExecute(ctx context.Context, taskId int64) error {
	// 查询任务详细信息
	taskInfo, err := l.svcCtx.SyncTaskModel.FindOne(ctx, taskId)
	if err != nil {
		logz.Errorf(ctx, "获取任务失败,Err:%v,Id:%d", err, taskId)
		_ = l.alarmService.ImportAlarm(ctx, taskInfo.OperateId, "获取任务失败", err)
		return err
	}

	// 文件下载到本地
	filename, err := public.DownloadFiled(ctx, l.svcCtx.Config.Oss, taskId, taskInfo.FileId)
	if err != nil {
		logz.Errorf(ctx, "导入源文件下载失败,Err:%v,Id:%d", err, taskId)
		_ = l.alarmService.ImportAlarm(ctx, taskInfo.OperateId, "导入源文件下载失败", err)
		return err
	}

	// 执行构造器
	_import.NewBuilder(l.svcCtx, taskInfo, filename).BuildStage(ctx)
	return nil
}
