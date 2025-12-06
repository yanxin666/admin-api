package session

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/logz"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/tencent"
	"encoding/json"
	"fmt"
	"muse-admin/internal/model/behavior"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListBySessionId struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetListBySessionId(ctx context.Context, svcCtx *svc.ServiceContext) *GetListBySessionId {
	return &GetListBySessionId{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetListBySessionId) GetListBySessionId(req *types.GetListBySessionIdReq) (*types.GetListBySessionIdResp, error) {
	// 参数校验
	if err := l.checkParams(req); err != nil {
		return nil, err
	}
	// 优先从数据库获取
	list, err := l.getFromDatabase(req.SessionId)
	if err != nil {
		return nil, fmt.Errorf("从数据库获取会话记录失败: %v", err)
	}
	if list != nil {
		return &types.GetListBySessionIdResp{List: list}, nil
	}

	// 从COS获取
	list, err = l.getFromCOS(req.SessionId)
	if err != nil {
		return nil, fmt.Errorf("从COS获取会话记录失败: %v", err)
	}
	return &types.GetListBySessionIdResp{List: list}, nil
}

// 参数校验
func (l *GetListBySessionId) checkParams(req *types.GetListBySessionIdReq) error {
	if req == nil {
		return errs.NewMsg(errs.ErrCodeProgram, "请求参数为空")
	}
	if req.SessionId <= 0 {
		return errs.NewMsg(errs.ErrCodeProgram, "无效的会话ID")
	}
	return nil
}

// 从数据库获取数据
func (l *GetListBySessionId) getFromDatabase(sessionId int64) ([]*types.EventRecord, error) {
	records, err := l.svcCtx.SessionRecordModel.FindListBySessionId(l.ctx, sessionId)
	if err != nil {
		logz.Errorf(l.ctx, "查询数据库失败, sessionId: %d, err: %v", sessionId, err)
		return nil, err
	}
	if records == nil {
		return nil, nil
	}
	return l.formatSessionRecord(records)
}

// 从COS获取数据
func (l *GetListBySessionId) getFromCOS(sessionId int64) ([]*types.EventRecord, error) {
	// 获取会话信息
	sessionInfo, err := l.svcCtx.SessionModel.FindOne(l.ctx, sessionId)
	if err != nil {
		logz.Errorf(l.ctx, "获取会话信息失败, sessionId: %d, err: %v", sessionId, err)
		return nil, err
	}
	if sessionInfo == nil {
		return make([]*types.EventRecord, 0), nil
	}

	// 从COS获取数据
	bytes, err := l.getObjectFromCos(fmt.Sprintf("%s.log", sessionInfo.Session))
	if err != nil {
		logz.Errorf(l.ctx, "从COS获取数据失败, session: %s, err: %v", sessionInfo.Session, err)
		return nil, err
	}
	logz.Infof(l.ctx, "从COS获取数据成功, session: %s, data: %s", sessionInfo.Session, string(bytes))
	// 解析数据
	var sessionRecord []*behavior.SessionRecord
	if err := json.Unmarshal(bytes, &sessionRecord); err != nil {
		logz.Errorf(l.ctx, "解析COS数据失败, 原始数据: %s, err: %v", string(bytes), err)
		return nil, errs.NewMsg(errs.ErrCodeProgram, "数据格式错误")
	}

	// 格式化数据
	if len(sessionRecord) == 0 {
		return make([]*types.EventRecord, 0), nil
	}
	return l.formatSessionRecord(sessionRecord)
}

// 获取远程cos对象数据
func (l *GetListBySessionId) getObjectFromCos(fileName string) ([]byte, error) {
	if fileName == "" {
		return nil, errs.NewMsg(errs.ErrCodeProgram, "文件名不能为空")
	}

	cosEnv := "test"
	if l.svcCtx.Config.Mode == "pro" {
		cosEnv = "prod"
	}

	cos := tencent.NewCos(l.ctx, tencent.CocConf{
		SecretId:  l.svcCtx.Config.Oss.SecretId,
		SecretKey: l.svcCtx.Config.Oss.SecretKey,
		Appid:     l.svcCtx.Config.Oss.Appid,
		Bucket:    l.svcCtx.Config.Oss.BehaviorBucket,
		Region:    l.svcCtx.Config.Oss.Region,
	})

	cosPath := fmt.Sprintf("app/%s/", cosEnv)
	return cos.OnlyDownLoadToByte(cosPath, fileName)
}

// 格式化数据
func (l *GetListBySessionId) formatSessionRecord(records []*behavior.SessionRecord) ([]*types.EventRecord, error) {
	if records == nil {
		return make([]*types.EventRecord, 0), nil
	}
	list := make([]*types.EventRecord, 0, len(records)) // 预分配容量
	for _, v := range records {
		list = append(list, &types.EventRecord{
			EventId:   v.EventId,
			EventType: v.EventType,
			EventTime: v.EventTime,
			Page:      v.Page,
			AppId:     v.Appid,
			Phone:     v.Phone,
			UserAgent: v.UserAgent.String,
			Screen:    v.Screen,
			Language:  v.Language,
			Data:      v.Data.String,
		})
	}
	return list, nil
}
