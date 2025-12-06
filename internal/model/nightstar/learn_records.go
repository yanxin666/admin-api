package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LearnRecordsModel = (*customLearnRecordsModel)(nil)

type (
	// LearnRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLearnRecordsModel.
	LearnRecordsModel interface {
		learnRecordsModel
		withSession(session sqlx.Session) LearnRecordsModel
		TableName() string
	}

	customLearnRecordsModel struct {
		*defaultLearnRecordsModel
	}
)

// NewLearnRecordsModel returns a model for the database table.
func NewLearnRecordsModel(conn sqlx.SqlConn) LearnRecordsModel {
	return &customLearnRecordsModel{
		defaultLearnRecordsModel: newLearnRecordsModel(conn),
	}
}

func (m *customLearnRecordsModel) withSession(session sqlx.Session) LearnRecordsModel {
	return NewLearnRecordsModel(sqlx.NewSqlConnFromSession(session))
}
