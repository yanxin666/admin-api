package knowledge

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ QuestionPointModel = (*customQuestionPointModel)(nil)

type (
	// QuestionPointModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionPointModel.
	QuestionPointModel interface {
		questionPointModel
		withSession(session sqlx.Session) QuestionPointModel
		TableName() string
	}

	customQuestionPointModel struct {
		*defaultQuestionPointModel
	}
)

// NewQuestionPointModel returns a model for the database table.
func NewQuestionPointModel(conn sqlx.SqlConn) QuestionPointModel {
	return &customQuestionPointModel{
		defaultQuestionPointModel: newQuestionPointModel(conn),
	}
}

func (m *customQuestionPointModel) withSession(session sqlx.Session) QuestionPointModel {
	return NewQuestionPointModel(sqlx.NewSqlConnFromSession(session))
}
