package knowledge

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LessonQuestionModel = (*customLessonQuestionModel)(nil)

type (
	// LessonQuestionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonQuestionModel.
	LessonQuestionModel interface {
		lessonQuestionModel
		withSession(session sqlx.Session) LessonQuestionModel
		TableName() string
		FindByQuestionNo(ctx context.Context, questionNo string) (*LessonQuestion, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *LessonQuestion) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *LessonQuestion) error

		FindAll(ctx context.Context) ([]*LessonQuestion, error)
	}

	customLessonQuestionModel struct {
		*defaultLessonQuestionModel
	}
)

// NewLessonQuestionModel returns a model for the database table.
func NewLessonQuestionModel(conn sqlx.SqlConn) LessonQuestionModel {
	return &customLessonQuestionModel{
		defaultLessonQuestionModel: newLessonQuestionModel(conn),
	}
}

func (m *customLessonQuestionModel) withSession(session sqlx.Session) LessonQuestionModel {
	return NewLessonQuestionModel(sqlx.NewSqlConnFromSession(session))
}

// FindByQuestionNo 根据编号查询
func (m *defaultLessonQuestionModel) FindByQuestionNo(ctx context.Context, questionNo string) (*LessonQuestion, error) {
	query := fmt.Sprintf("select %s from %s where `question_no` = ? limit 1", lessonQuestionRows, m.table)
	var resp LessonQuestion
	err := m.conn.QueryRowCtx(ctx, &resp, query, questionNo)
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
func (m *customLessonQuestionModel) InsertSession(ctx context.Context, session sqlx.Session, data *LessonQuestion) (int64, error) {
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
func (m *customLessonQuestionModel) UpdateSession(ctx context.Context, session sqlx.Session, data *LessonQuestion) error {
	return m.withSession(session).Update(ctx, data)
}

func (m *customLessonQuestionModel) FindAll(ctx context.Context) ([]*LessonQuestion, error) {
	query := fmt.Sprintf("SELECT %s FROM %s ", lessonQuestionRows, m.table)
	var resp []*LessonQuestion
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
