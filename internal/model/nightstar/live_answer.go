package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LiveAnswerModel = (*customLiveAnswerModel)(nil)

type (
	// LiveAnswerModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveAnswerModel.
	LiveAnswerModel interface {
		liveAnswerModel
		withSession(session sqlx.Session) LiveAnswerModel
		TableName() string
	}

	customLiveAnswerModel struct {
		*defaultLiveAnswerModel
	}
)

// NewLiveAnswerModel returns a model for the database table.
func NewLiveAnswerModel(conn sqlx.SqlConn) LiveAnswerModel {
	return &customLiveAnswerModel{
		defaultLiveAnswerModel: newLiveAnswerModel(conn),
	}
}

func (m *customLiveAnswerModel) withSession(session sqlx.Session) LiveAnswerModel {
	return NewLiveAnswerModel(sqlx.NewSqlConnFromSession(session))
}
