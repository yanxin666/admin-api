package supertrain

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SubTaskModel = (*customSubTaskModel)(nil)

type (
	// SubTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSubTaskModel.
	SubTaskModel interface {
		subTaskModel
		withSession(session sqlx.Session) SubTaskModel
		TableName() string
		FindOneByTaskIdAndSubTaskNo(ctx context.Context, taskId int64, subTaskNo string) (*SubTask, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *SubTask) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *SubTask) error
	}

	customSubTaskModel struct {
		*defaultSubTaskModel
	}
)

// NewSubTaskModel returns a model for the database table.
func NewSubTaskModel(conn sqlx.SqlConn) SubTaskModel {
	return &customSubTaskModel{
		defaultSubTaskModel: newSubTaskModel(conn),
	}
}

func (m *customSubTaskModel) withSession(session sqlx.Session) SubTaskModel {
	return NewSubTaskModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSubTaskModel) FindOneByTaskIdAndSubTaskNo(ctx context.Context, taskId int64, subTaskNo string) (*SubTask, error) {
	query := fmt.Sprintf("select %s from %s where `chapter_task_id` = ? and `sub_task_no` = ? limit 1", subTaskRows, m.table)
	var resp SubTask
	err := m.conn.QueryRowCtx(ctx, &resp, query, taskId, subTaskNo)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

// InsertSession 事务操作-新增
func (m *customSubTaskModel) InsertSession(ctx context.Context, session sqlx.Session, data *SubTask) (int64, error) {
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
func (m *customSubTaskModel) UpdateSession(ctx context.Context, session sqlx.Session, data *SubTask) error {
	return m.withSession(session).Update(ctx, data)
}
