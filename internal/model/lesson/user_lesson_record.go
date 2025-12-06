package lesson

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserLessonRecordModel = (*customUserLessonRecordModel)(nil)

type (
	// UserLessonRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLessonRecordModel.
	UserLessonRecordModel interface {
		userLessonRecordModel
		withSession(session sqlx.Session) UserLessonRecordModel
		TableName() string
	}

	customUserLessonRecordModel struct {
		*defaultUserLessonRecordModel
	}
)

// NewUserLessonRecordModel returns a model for the database table.
func NewUserLessonRecordModel(conn sqlx.SqlConn) UserLessonRecordModel {
	return &customUserLessonRecordModel{
		defaultUserLessonRecordModel: newUserLessonRecordModel(conn),
	}
}

func (m *customUserLessonRecordModel) withSession(session sqlx.Session) UserLessonRecordModel {
	return NewUserLessonRecordModel(sqlx.NewSqlConnFromSession(session))
}
