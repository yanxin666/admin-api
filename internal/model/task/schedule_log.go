package task

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ScheduleLogModel = (*customScheduleLogModel)(nil)

type (
	// ScheduleLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customScheduleLogModel.
	ScheduleLogModel interface {
		scheduleLogModel
		withSession(session sqlx.Session) ScheduleLogModel
		TableName() string
	}

	customScheduleLogModel struct {
		*defaultScheduleLogModel
	}
)

// NewScheduleLogModel returns a model for the database table.
func NewScheduleLogModel(conn sqlx.SqlConn) ScheduleLogModel {
	return &customScheduleLogModel{
		defaultScheduleLogModel: newScheduleLogModel(conn),
	}
}

func (m *customScheduleLogModel) withSession(session sqlx.Session) ScheduleLogModel {
	return NewScheduleLogModel(sqlx.NewSqlConnFromSession(session))
}
