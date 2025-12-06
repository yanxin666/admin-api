package live_course

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SignalLogModel = (*customSignalLogModel)(nil)

type (
	// SignalLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSignalLogModel.
	SignalLogModel interface {
		signalLogModel
		withSession(session sqlx.Session) SignalLogModel
		TableName() string
	}

	customSignalLogModel struct {
		*defaultSignalLogModel
	}
)

// NewSignalLogModel returns a model for the database table.
func NewSignalLogModel(conn sqlx.SqlConn) SignalLogModel {
	return &customSignalLogModel{
		defaultSignalLogModel: newSignalLogModel(conn),
	}
}

func (m *customSignalLogModel) withSession(session sqlx.Session) SignalLogModel {
	return NewSignalLogModel(sqlx.NewSqlConnFromSession(session))
}
