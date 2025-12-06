package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomeworkDetailModel = (*customHomeworkDetailModel)(nil)

type (
	// HomeworkDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomeworkDetailModel.
	HomeworkDetailModel interface {
		homeworkDetailModel
		withSession(session sqlx.Session) HomeworkDetailModel
		TableName() string
	}

	customHomeworkDetailModel struct {
		*defaultHomeworkDetailModel
	}
)

// NewHomeworkDetailModel returns a model for the database table.
func NewHomeworkDetailModel(conn sqlx.SqlConn) HomeworkDetailModel {
	return &customHomeworkDetailModel{
		defaultHomeworkDetailModel: newHomeworkDetailModel(conn),
	}
}

func (m *customHomeworkDetailModel) withSession(session sqlx.Session) HomeworkDetailModel {
	return NewHomeworkDetailModel(sqlx.NewSqlConnFromSession(session))
}
