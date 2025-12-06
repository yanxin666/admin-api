package lesson

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserLessonDetailModel = (*customUserLessonDetailModel)(nil)

type (
	// UserLessonDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLessonDetailModel.
	UserLessonDetailModel interface {
		userLessonDetailModel
		withSession(session sqlx.Session) UserLessonDetailModel
		TableName() string
	}

	customUserLessonDetailModel struct {
		*defaultUserLessonDetailModel
	}
)

// NewUserLessonDetailModel returns a model for the database table.
func NewUserLessonDetailModel(conn sqlx.SqlConn) UserLessonDetailModel {
	return &customUserLessonDetailModel{
		defaultUserLessonDetailModel: newUserLessonDetailModel(conn),
	}
}

func (m *customUserLessonDetailModel) withSession(session sqlx.Session) UserLessonDetailModel {
	return NewUserLessonDetailModel(sqlx.NewSqlConnFromSession(session))
}
