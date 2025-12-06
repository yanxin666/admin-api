package knowledge

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LessonQuestionOptionModel = (*customLessonQuestionOptionModel)(nil)

type (
	// LessonQuestionOptionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonQuestionOptionModel.
	LessonQuestionOptionModel interface {
		lessonQuestionOptionModel
		withSession(session sqlx.Session) LessonQuestionOptionModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *LessonQuestionOption) (int64, error)
	}

	customLessonQuestionOptionModel struct {
		*defaultLessonQuestionOptionModel
	}
)

// NewLessonQuestionOptionModel returns a model for the database table.
func NewLessonQuestionOptionModel(conn sqlx.SqlConn) LessonQuestionOptionModel {
	return &customLessonQuestionOptionModel{
		defaultLessonQuestionOptionModel: newLessonQuestionOptionModel(conn),
	}
}

func (m *customLessonQuestionOptionModel) withSession(session sqlx.Session) LessonQuestionOptionModel {
	return NewLessonQuestionOptionModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customLessonQuestionOptionModel) InsertSession(ctx context.Context, session sqlx.Session, data *LessonQuestionOption) (int64, error) {
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
