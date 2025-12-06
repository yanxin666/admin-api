package live

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseModel = (*customCourseModel)(nil)

type (
	// CourseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseModel.
	CourseModel interface {
		courseModel
		withSession(session sqlx.Session) CourseModel
		TableName() string
		FindOneByNo(ctx context.Context, no int64) (*Course, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *Course) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *Course) error
	}

	customCourseModel struct {
		*defaultCourseModel
	}
)

// NewCourseModel returns a model for the database table.
func NewCourseModel(conn sqlx.SqlConn) CourseModel {
	return &customCourseModel{
		defaultCourseModel: newCourseModel(conn),
	}
}

func (m *customCourseModel) withSession(session sqlx.Session) CourseModel {
	return NewCourseModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customCourseModel) FindOneByNo(ctx context.Context, no int64) (*Course, error) {
	query := fmt.Sprintf("select %s from %s where `course_no` = ? limit 1", courseRows, m.table)
	var resp Course
	err := m.conn.QueryRowCtx(ctx, &resp, query, no)
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
func (m *customCourseModel) InsertSession(ctx context.Context, session sqlx.Session, data *Course) (int64, error) {
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
func (m *customCourseModel) UpdateSession(ctx context.Context, session sqlx.Session, data *Course) error {
	return m.withSession(session).Update(ctx, data)
}
