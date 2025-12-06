package lesson

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserScheduleModel = (*customUserScheduleModel)(nil)

type (
	// UserScheduleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserScheduleModel.
	UserScheduleModel interface {
		userScheduleModel
		withSession(session sqlx.Session) UserScheduleModel
		TableName() string
	}

	customUserScheduleModel struct {
		*defaultUserScheduleModel
	}
)

// NewUserScheduleModel returns a model for the database table.
func NewUserScheduleModel(conn sqlx.SqlConn) UserScheduleModel {
	return &customUserScheduleModel{
		defaultUserScheduleModel: newUserScheduleModel(conn),
	}
}

func (m *customUserScheduleModel) withSession(session sqlx.Session) UserScheduleModel {
	return NewUserScheduleModel(sqlx.NewSqlConnFromSession(session))
}
