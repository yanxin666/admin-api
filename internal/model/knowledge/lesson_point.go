package knowledge

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LessonPointModel = (*customLessonPointModel)(nil)

type (
	// LessonPointModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonPointModel.
	LessonPointModel interface {
		lessonPointModel
		withSession(session sqlx.Session) LessonPointModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *LessonPoint) (int64, error)
		FindResourceId(ctx context.Context, data *LessonPoint) (int64, error)
	}

	customLessonPointModel struct {
		*defaultLessonPointModel
	}
)

// NewLessonPointModel returns a model for the database table.
func NewLessonPointModel(conn sqlx.SqlConn) LessonPointModel {
	return &customLessonPointModel{
		defaultLessonPointModel: newLessonPointModel(conn),
	}
}

func (m *customLessonPointModel) withSession(session sqlx.Session) LessonPointModel {
	return NewLessonPointModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customLessonPointModel) InsertSession(ctx context.Context, session sqlx.Session, data *LessonPoint) (int64, error) {
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

// FindResourceId 根据data内容来获资源ID
func (m *customLessonPointModel) FindResourceId(ctx context.Context, data *LessonPoint) (int64, error) {
	query := fmt.Sprintf("select `id` from %s where `lesson_id` = ? and `name` = ? ", m.table)
	var id int64
	err := m.conn.QueryRowPartialCtx(ctx, &id, query, data.LessonId, data.Name)
	switch {
	case err == nil:
		return id, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}
