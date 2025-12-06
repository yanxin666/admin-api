package supervisor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"muse-admin/internal/model/tools"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ScheduleModel = (*customScheduleModel)(nil)

type (
	// ScheduleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customScheduleModel.
	ScheduleModel interface {
		scheduleModel
		withSession(session sqlx.Session) ScheduleModel
		TableName() string
		FindPageByCondition(ctx context.Context, page, limit int64, condition map[string]interface{}) ([]Schedule, int64, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *Schedule) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *Schedule) error
		FindActiveSchedule(ctx context.Context) ([]Schedule, error)
		UpdateFillFieldsById(ctx context.Context, id int64, update *Schedule) (sql.Result, error)
	}

	customScheduleModel struct {
		*defaultScheduleModel
	}
)

// NewScheduleModel returns a model for the database table.
func NewScheduleModel(conn sqlx.SqlConn) ScheduleModel {
	return &customScheduleModel{
		defaultScheduleModel: newScheduleModel(conn),
	}
}

func (m *customScheduleModel) withSession(session sqlx.Session) ScheduleModel {
	return NewScheduleModel(sqlx.NewSqlConnFromSession(session))
}

// FindPageByCondition 分页查询
func (m *customScheduleModel) FindPageByCondition(ctx context.Context, page int64, limit int64, condition map[string]interface{}) ([]Schedule, int64, error) {
	var resp []Schedule
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*")
	total, err := query.WhereMap(condition).Order("id desc").Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}

// InsertSession 事务操作-新增
func (m *customScheduleModel) InsertSession(ctx context.Context, session sqlx.Session, data *Schedule) (int64, error) {
	result, err := m.withSession(session).Insert(ctx, data)
	if err != nil {
		return 0, err
	}

	// 获取新增ID
	aid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if aid == 0 {
		return 0, errors.New("新增事物失败，未生成自增ID")
	}

	return aid, nil
}

// UpdateSession 事务操作-更新
func (m *customScheduleModel) UpdateSession(ctx context.Context, session sqlx.Session, data *Schedule) error {
	return m.withSession(session).Update(ctx, data)
}

// FindActiveSchedule 查询进行中的督学课, 限制200条
func (m *customScheduleModel) FindActiveSchedule(ctx context.Context) ([]Schedule, error) {
	query := fmt.Sprintf("select %s from %s where status = 3 limit 200", scheduleRows, m.table)
	var resp []Schedule
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	return resp, err
}

// UpdateFillFieldsById 根据id 更新用户信息 支持自定义struct
func (m *customScheduleModel) UpdateFillFieldsById(ctx context.Context, id int64, data *Schedule) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, paramSet)
	args = append(args, id)

	result, err := m.conn.ExecCtx(ctx, query, args...)
	return result, err
}
