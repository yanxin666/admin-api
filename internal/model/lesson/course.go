package lesson

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CourseModel = (*customCourseModel)(nil)

type (
	// CourseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseModel.
	CourseModel interface {
		courseModel
		withSession(session sqlx.Session) CourseModel
		TableName() string
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
