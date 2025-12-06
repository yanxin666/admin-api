package lesson

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserScheduleDetailModel = (*customUserScheduleDetailModel)(nil)

type (
	// UserScheduleDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserScheduleDetailModel.
	UserScheduleDetailModel interface {
		userScheduleDetailModel
		withSession(session sqlx.Session) UserScheduleDetailModel
		TableName() string
	}

	customUserScheduleDetailModel struct {
		*defaultUserScheduleDetailModel
	}
)

// NewUserScheduleDetailModel returns a model for the database table.
func NewUserScheduleDetailModel(conn sqlx.SqlConn) UserScheduleDetailModel {
	return &customUserScheduleDetailModel{
		defaultUserScheduleDetailModel: newUserScheduleDetailModel(conn),
	}
}

func (m *customUserScheduleDetailModel) withSession(session sqlx.Session) UserScheduleDetailModel {
	return NewUserScheduleDetailModel(sqlx.NewSqlConnFromSession(session))
}
