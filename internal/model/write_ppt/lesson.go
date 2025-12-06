package write_ppt

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LessonModel = (*customLessonModel)(nil)

type (
	// LessonModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonModel.
	LessonModel interface {
		lessonModel
		withSession(session sqlx.Session) LessonModel
		TableName() string
		FindByLessonNo(ctx context.Context, lessonNo int64) (*Lesson, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *Lesson) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *Lesson) error
	}

	customLessonModel struct {
		*defaultLessonModel
	}
)

// NewLessonModel returns a model for the database table.
func NewLessonModel(conn sqlx.SqlConn) LessonModel {
	return &customLessonModel{
		defaultLessonModel: newLessonModel(conn),
	}
}

func (m *customLessonModel) withSession(session sqlx.Session) LessonModel {
	return NewLessonModel(sqlx.NewSqlConnFromSession(session))
}

// FindByLessonNo 根据编号查询
func (m *customLessonModel) FindByLessonNo(ctx context.Context, lessonNo int64) (*Lesson, error) {
	query := fmt.Sprintf("select %s from %s where `lesson_no` = ? limit 1", lessonRows, m.table)
	var resp Lesson
	err := m.conn.QueryRowCtx(ctx, &resp, query, lessonNo)
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
func (m *customLessonModel) InsertSession(ctx context.Context, session sqlx.Session, data *Lesson) (int64, error) {
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
func (m *customLessonModel) UpdateSession(ctx context.Context, session sqlx.Session, data *Lesson) error {
	return m.withSession(session).Update(ctx, data)
}
