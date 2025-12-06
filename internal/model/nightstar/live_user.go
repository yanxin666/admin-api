package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LiveUserModel = (*customLiveUserModel)(nil)

type (
	// LiveUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveUserModel.
	LiveUserModel interface {
		liveUserModel
		withSession(session sqlx.Session) LiveUserModel
		TableName() string
	}

	customLiveUserModel struct {
		*defaultLiveUserModel
	}
)

// NewLiveUserModel returns a model for the database table.
func NewLiveUserModel(conn sqlx.SqlConn) LiveUserModel {
	return &customLiveUserModel{
		defaultLiveUserModel: newLiveUserModel(conn),
	}
}

func (m *customLiveUserModel) withSession(session sqlx.Session) LiveUserModel {
	return NewLiveUserModel(sqlx.NewSqlConnFromSession(session))
}
