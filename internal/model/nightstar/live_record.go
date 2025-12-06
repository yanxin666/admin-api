package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LiveRecordModel = (*customLiveRecordModel)(nil)

type (
	// LiveRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveRecordModel.
	LiveRecordModel interface {
		liveRecordModel
		withSession(session sqlx.Session) LiveRecordModel
		TableName() string
	}

	customLiveRecordModel struct {
		*defaultLiveRecordModel
	}
)

// NewLiveRecordModel returns a model for the database table.
func NewLiveRecordModel(conn sqlx.SqlConn) LiveRecordModel {
	return &customLiveRecordModel{
		defaultLiveRecordModel: newLiveRecordModel(conn),
	}
}

func (m *customLiveRecordModel) withSession(session sqlx.Session) LiveRecordModel {
	return NewLiveRecordModel(sqlx.NewSqlConnFromSession(session))
}
