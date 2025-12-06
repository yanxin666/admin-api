package lesson

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserLessonAnswerModel = (*customUserLessonAnswerModel)(nil)

type (
	// UserLessonAnswerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLessonAnswerModel.
	UserLessonAnswerModel interface {
		userLessonAnswerModel
		withSession(session sqlx.Session) UserLessonAnswerModel
		TableName() string
	}

	customUserLessonAnswerModel struct {
		*defaultUserLessonAnswerModel
	}
)

// NewUserLessonAnswerModel returns a model for the database table.
func NewUserLessonAnswerModel(conn sqlx.SqlConn) UserLessonAnswerModel {
	return &customUserLessonAnswerModel{
		defaultUserLessonAnswerModel: newUserLessonAnswerModel(conn),
	}
}

func (m *customUserLessonAnswerModel) withSession(session sqlx.Session) UserLessonAnswerModel {
	return NewUserLessonAnswerModel(sqlx.NewSqlConnFromSession(session))
}
