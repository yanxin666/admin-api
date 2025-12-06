package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LiveVirtualUserModel = (*customLiveVirtualUserModel)(nil)

type (
	// LiveVirtualUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveVirtualUserModel.
	LiveVirtualUserModel interface {
		liveVirtualUserModel
		withSession(session sqlx.Session) LiveVirtualUserModel
		TableName() string
	}

	customLiveVirtualUserModel struct {
		*defaultLiveVirtualUserModel
	}
)

// NewLiveVirtualUserModel returns a model for the database table.
func NewLiveVirtualUserModel(conn sqlx.SqlConn) LiveVirtualUserModel {
	return &customLiveVirtualUserModel{
		defaultLiveVirtualUserModel: newLiveVirtualUserModel(conn),
	}
}

func (m *customLiveVirtualUserModel) withSession(session sqlx.Session) LiveVirtualUserModel {
	return NewLiveVirtualUserModel(sqlx.NewSqlConnFromSession(session))
}
