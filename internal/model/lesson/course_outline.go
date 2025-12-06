package lesson

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseOutlineModel = (*customCourseOutlineModel)(nil)

type (
	// CourseOutlineModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseOutlineModel.
	CourseOutlineModel interface {
		courseOutlineModel
		withSession(session sqlx.Session) CourseOutlineModel
		TableName() string
		FindAllByRemark(ctx context.Context, remark string) ([]CourseOutline, error)
		FindAllByMarkId(ctx context.Context, markId int64) ([]CourseOutline, error)
	}

	customCourseOutlineModel struct {
		*defaultCourseOutlineModel
	}
)

// NewCourseOutlineModel returns a model for the database table.
func NewCourseOutlineModel(conn sqlx.SqlConn) CourseOutlineModel {
	return &customCourseOutlineModel{
		defaultCourseOutlineModel: newCourseOutlineModel(conn),
	}
}

func (m *customCourseOutlineModel) withSession(session sqlx.Session) CourseOutlineModel {
	return NewCourseOutlineModel(sqlx.NewSqlConnFromSession(session))
}

// FindAllByRemark 根据remark来获大纲信息
func (m *customCourseOutlineModel) FindAllByRemark(ctx context.Context, remark string) ([]CourseOutline, error) {
	query := fmt.Sprintf("select %s from %s where `remark` = ?", courseOutlineRows, m.table)
	var resp []CourseOutline
	err := m.conn.QueryRowsCtx(ctx, &resp, query, remark)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// FindAllByMarkId 根据mark_id来获大纲信息
func (m *customCourseOutlineModel) FindAllByMarkId(ctx context.Context, markId int64) ([]CourseOutline, error) {
	query := fmt.Sprintf("select %s from %s where `mark_id` = ?", courseOutlineRows, m.table)
	var resp []CourseOutline
	err := m.conn.QueryRowsCtx(ctx, &resp, query, markId)
	if err != nil {
		return nil, err
	}

	return resp, err
}
