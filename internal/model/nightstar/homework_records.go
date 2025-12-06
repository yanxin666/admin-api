package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ HomeworkRecordsModel = (*customHomeworkRecordsModel)(nil)

type (
	// HomeworkRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customHomeworkRecordsModel.
	HomeworkRecordsModel interface {
		homeworkRecordsModel
		withSession(session sqlx.Session) HomeworkRecordsModel
		TableName() string
	}

	customHomeworkRecordsModel struct {
		*defaultHomeworkRecordsModel
	}
)

// NewHomeworkRecordsModel returns a model for the database table.
func NewHomeworkRecordsModel(conn sqlx.SqlConn) HomeworkRecordsModel {
	return &customHomeworkRecordsModel{
		defaultHomeworkRecordsModel: newHomeworkRecordsModel(conn),
	}
}

func (m *customHomeworkRecordsModel) withSession(session sqlx.Session) HomeworkRecordsModel {
	return NewHomeworkRecordsModel(sqlx.NewSqlConnFromSession(session))
}
