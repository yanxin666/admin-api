package supertrain

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SubTaskMedalModel = (*customSubTaskMedalModel)(nil)

type (
	// SubTaskMedalModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSubTaskMedalModel.
	SubTaskMedalModel interface {
		subTaskMedalModel
		DeleteBySubTaskId(ctx context.Context, subTaskId int64) error
		withSession(session sqlx.Session) SubTaskMedalModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *SubTaskMedal) (int64, error)
		DeleteSessionBySubTaskId(ctx context.Context, session sqlx.Session, subTaskId int64) error
	}

	customSubTaskMedalModel struct {
		*defaultSubTaskMedalModel
	}
)

// NewSubTaskMedalModel returns a model for the database table.
func NewSubTaskMedalModel(conn sqlx.SqlConn) SubTaskMedalModel {
	return &customSubTaskMedalModel{
		defaultSubTaskMedalModel: newSubTaskMedalModel(conn),
	}
}

func (m *customSubTaskMedalModel) DeleteBySubTaskId(ctx context.Context, subTaskId int64) error {
	query := fmt.Sprintf("delete from %s where `sub_task_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, subTaskId)
	return err
}

func (m *customSubTaskMedalModel) withSession(session sqlx.Session) SubTaskMedalModel {
	return NewSubTaskMedalModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customSubTaskMedalModel) InsertSession(ctx context.Context, session sqlx.Session, data *SubTaskMedal) (int64, error) {
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

func (m *customSubTaskMedalModel) DeleteSessionBySubTaskId(ctx context.Context, session sqlx.Session, subTaskId int64) error {
	return m.withSession(session).DeleteBySubTaskId(ctx, subTaskId)
}
