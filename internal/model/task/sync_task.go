package task

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/model/tools"
	"muse-admin/pkg/errs"
)

var _ SyncTaskModel = (*customSyncTaskModel)(nil)

type (
	// SyncTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSyncTaskModel.
	SyncTaskModel interface {
		syncTaskModel
		withSession(session sqlx.Session) SyncTaskModel
		TableName() string
		UpdateById(ctx context.Context, id int64, update map[string]interface{}) error
		UpdateFillFieldsById(ctx context.Context, id int64, update *SyncTask) (sql.Result, error)
		UpdateFailById(ctx context.Context, id int64, errMsg string) error
		FindPageByCondition(ctx context.Context, page int64, limit int64, condition map[string]interface{}) ([]SyncTask, int64, error)
	}

	customSyncTaskModel struct {
		*defaultSyncTaskModel
	}
)

// NewSyncTaskModel returns a model for the database table.
func NewSyncTaskModel(conn sqlx.SqlConn) SyncTaskModel {
	return &customSyncTaskModel{
		defaultSyncTaskModel: newSyncTaskModel(conn),
	}
}

func (m *customSyncTaskModel) withSession(session sqlx.Session) SyncTaskModel {
	return NewSyncTaskModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSyncTaskModel) UpdateById(ctx context.Context, id int64, update map[string]interface{}) error {
	db := sqld.NewModel(ctx, m.conn, m.table)
	_, err := db.Where("id =?", id).UpdateColumnMap(update)
	return err
}

// UpdateFillFieldsById 根据id 更新用户信息 支持自定义struct
func (m *customSyncTaskModel) UpdateFillFieldsById(ctx context.Context, id int64, data *SyncTask) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, paramSet)
	args = append(args, id)

	result, err := m.conn.ExecCtx(ctx, query, args...)
	return result, err
}

// UpdateFailById 根据id 更新任务失败时的信息
func (m *customSyncTaskModel) UpdateFailById(ctx context.Context, id int64, errMsg string) error {
	_, err := m.UpdateFillFieldsById(ctx, id, &SyncTask{
		Status:  3,
		EndTime: util.GetStandardNowDatetime(),
		ErrorMsg: sql.NullString{
			String: errMsg,
			Valid:  true,
		},
	})
	if err != nil {
		return errs.WithMsg(err, errs.ErrCodeAbnormal, "更新任务结果失败")
	}

	return nil
}

// FindPageByCondition 分页查询
func (m *customSyncTaskModel) FindPageByCondition(ctx context.Context, page int64, limit int64, condition map[string]interface{}) ([]SyncTask, int64, error) {
	var resp []SyncTask
	db := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	total, err := db.Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}
