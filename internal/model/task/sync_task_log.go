package task

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/model/tools"
)

var _ SyncTaskLogModel = (*customSyncTaskLogModel)(nil)

type (
	// SyncTaskLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSyncTaskLogModel.
	SyncTaskLogModel interface {
		syncTaskLogModel
		withSession(session sqlx.Session) SyncTaskLogModel
		TableName() string
		UpdateFillFields(ctx context.Context, taskId, index int64, update *SyncTaskLog) (sql.Result, error)
		FindPageByCondition(ctx context.Context, page, limit int64, keyword string, condition map[string]interface{}) ([]SyncTaskLog, int64, error)
		FindOneByMd5(ctx context.Context, md5 string) (*SyncTaskLog, error)
		FindOneByTaskIdAndIndex(ctx context.Context, taskId, index int64) (*SyncTaskLog, error)
		FindOneByMd5AndStatus(ctx context.Context, md5 string, status int64) (*SyncTaskLog, error)
		UpdateStatusSuc(ctx context.Context, taskId int64) error
	}

	customSyncTaskLogModel struct {
		*defaultSyncTaskLogModel
	}
)

// NewSyncTaskLogModel returns a model for the database table.
func NewSyncTaskLogModel(conn sqlx.SqlConn) SyncTaskLogModel {
	return &customSyncTaskLogModel{
		defaultSyncTaskLogModel: newSyncTaskLogModel(conn),
	}
}

func (m *customSyncTaskLogModel) withSession(session sqlx.Session) SyncTaskLogModel {
	return NewSyncTaskLogModel(sqlx.NewSqlConnFromSession(session))
}

// UpdateFillFields 根据taskId、index 更新用户信息 支持自定义struct
func (m *customSyncTaskLogModel) UpdateFillFields(ctx context.Context, taskId, index int64, data *SyncTaskLog) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE task_id = ? and `index` = ?", m.table, paramSet)
	args = append(args, taskId, index)

	result, err := m.conn.ExecCtx(ctx, query, args...)
	return result, err
}

func (m *customSyncTaskLogModel) FindPageByCondition(ctx context.Context, page, limit int64, keyword string, condition map[string]interface{}) ([]SyncTaskLog, int64, error) {
	var resp []SyncTaskLog
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if keyword != "" {
		query = query.Where("`data` like ?", "%"+cast.ToString(keyword)+"%")
	}
	total, err := query.Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}

func (m *customSyncTaskLogModel) FindOneByMd5(ctx context.Context, md5 string) (*SyncTaskLog, error) {
	query := fmt.Sprintf("select %s from %s where `md5` = ? limit 1", syncTaskLogRows, m.table)
	var resp SyncTaskLog
	err := m.conn.QueryRowCtx(ctx, &resp, query, md5)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customSyncTaskLogModel) FindOneByTaskIdAndIndex(ctx context.Context, taskId, index int64) (*SyncTaskLog, error) {
	query := fmt.Sprintf("select %s from %s where `task_id` = ? and `index` = ? limit 1", syncTaskLogRows, m.table)
	var resp SyncTaskLog
	err := m.conn.QueryRowCtx(ctx, &resp, query, taskId, index)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customSyncTaskLogModel) FindOneByMd5AndStatus(ctx context.Context, md5 string, status int64) (*SyncTaskLog, error) {
	query := fmt.Sprintf("select %s from %s where `md5` = ? and `status` = ? limit 1", syncTaskLogRows, m.table)
	var resp SyncTaskLog
	err := m.conn.QueryRowCtx(ctx, &resp, query, md5, status)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

// UpdateStatusSuc 根据taskId更新data的成功状态
func (m *customSyncTaskLogModel) UpdateStatusSuc(ctx context.Context, taskId int64) error {
	query := fmt.Sprintf("UPDATE %s SET `status` = 2 WHERE task_id = ? and `status` IN (1, 3)", m.table)
	_, err := m.conn.ExecCtx(ctx, query, taskId)
	return err
}
