package knowledge

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ QuestionOptionModel = (*customQuestionOptionModel)(nil)

type (
	// QuestionOptionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionOptionModel.
	QuestionOptionModel interface {
		questionOptionModel
		withSession(session sqlx.Session) QuestionOptionModel
		TableName() string
	}

	customQuestionOptionModel struct {
		*defaultQuestionOptionModel
	}
)

// NewQuestionOptionModel returns a model for the database table.
func NewQuestionOptionModel(conn sqlx.SqlConn) QuestionOptionModel {
	return &customQuestionOptionModel{
		defaultQuestionOptionModel: newQuestionOptionModel(conn),
	}
}

func (m *customQuestionOptionModel) withSession(session sqlx.Session) QuestionOptionModel {
	return NewQuestionOptionModel(sqlx.NewSqlConnFromSession(session))
}
