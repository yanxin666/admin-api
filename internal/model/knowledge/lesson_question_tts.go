package knowledge

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LessonQuestionTtsModel = (*customLessonQuestionTtsModel)(nil)

type (
	// LessonQuestionTtsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonQuestionTtsModel.
	LessonQuestionTtsModel interface {
		lessonQuestionTtsModel
		withSession(session sqlx.Session) LessonQuestionTtsModel
		TableName() string
		FindByQuestionId(ctx context.Context, questionId int64) (*LessonQuestionTts, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *LessonQuestionTts) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *LessonQuestionTts) error
	}

	customLessonQuestionTtsModel struct {
		*defaultLessonQuestionTtsModel
	}
)

// NewLessonQuestionTtsModel returns a model for the database table.
func NewLessonQuestionTtsModel(conn sqlx.SqlConn) LessonQuestionTtsModel {
	return &customLessonQuestionTtsModel{
		defaultLessonQuestionTtsModel: newLessonQuestionTtsModel(conn),
	}
}

func (m *customLessonQuestionTtsModel) withSession(session sqlx.Session) LessonQuestionTtsModel {
	return NewLessonQuestionTtsModel(sqlx.NewSqlConnFromSession(session))
}

// FindByQuestionId 根据题目ID查询
func (m *customLessonQuestionTtsModel) FindByQuestionId(ctx context.Context, questionId int64) (*LessonQuestionTts, error) {
	query := fmt.Sprintf("select %s from %s where `question_id` = ? limit 1", lessonQuestionTtsRows, m.table)
	var resp LessonQuestionTts
	err := m.conn.QueryRowCtx(ctx, &resp, query, questionId)
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
func (m *customLessonQuestionTtsModel) InsertSession(ctx context.Context, session sqlx.Session, data *LessonQuestionTts) (int64, error) {
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
func (m *customLessonQuestionTtsModel) UpdateSession(ctx context.Context, session sqlx.Session, data *LessonQuestionTts) error {
	return m.withSession(session).Update(ctx, data)
}
