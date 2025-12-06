package member

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserOpinionModel = (*customUserOpinionModel)(nil)

type (
	// UserOpinionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserOpinionModel.
	UserOpinionModel interface {
		userOpinionModel
		withSession(session sqlx.Session) UserOpinionModel
		TableName() string
	}

	customUserOpinionModel struct {
		*defaultUserOpinionModel
	}
)

// NewUserOpinionModel returns a model for the database table.
func NewUserOpinionModel(conn sqlx.SqlConn) UserOpinionModel {
	return &customUserOpinionModel{
		defaultUserOpinionModel: newUserOpinionModel(conn),
	}
}

func (m *customUserOpinionModel) withSession(session sqlx.Session) UserOpinionModel {
	return NewUserOpinionModel(sqlx.NewSqlConnFromSession(session))
}
