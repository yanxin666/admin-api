package knowledge

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ExampleModel = (*customExampleModel)(nil)

type (
	// ExampleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customExampleModel.
	ExampleModel interface {
		exampleModel
		withSession(session sqlx.Session) ExampleModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *Example) (int64, error)
		FindExampleId(ctx context.Context, data *Example) (int64, error)
	}

	customExampleModel struct {
		*defaultExampleModel
	}
)

// NewExampleModel returns a model for the database table.
func NewExampleModel(conn sqlx.SqlConn) ExampleModel {
	return &customExampleModel{
		defaultExampleModel: newExampleModel(conn),
	}
}

func (m *customExampleModel) withSession(session sqlx.Session) ExampleModel {
	return NewExampleModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customExampleModel) InsertSession(ctx context.Context, session sqlx.Session, data *Example) (int64, error) {
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

// FindExampleId 根据data内容来获取例题ID
func (m *customExampleModel) FindExampleId(ctx context.Context, data *Example) (int64, error) {
	query := fmt.Sprintf("select `id` from %s where `explain` = ? and `plain_text` = ? and `title` = ? and `sub_title` = ? and `lesson_notes` = ?", m.table)
	var id int64
	err := m.conn.QueryRowPartialCtx(ctx, &id, query, data.Explain, data.PlainText, data.Title, data.SubTitle, data.LessonNotes)
	switch {
	case err == nil:
		return id, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}
