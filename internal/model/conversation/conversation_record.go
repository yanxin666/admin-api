package conversation

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ConversationRecordModel = (*customConversationRecordModel)(nil)

type (
	// ConversationRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationRecordModel.
	ConversationRecordModel interface {
		conversationRecordModel
		withSession(session sqlx.Session) ConversationRecordModel
		TableName() string
	}

	customConversationRecordModel struct {
		*defaultConversationRecordModel
	}
)

// NewConversationRecordModel returns a model for the database table.
func NewConversationRecordModel(conn sqlx.SqlConn) ConversationRecordModel {
	return &customConversationRecordModel{
		defaultConversationRecordModel: newConversationRecordModel(conn),
	}
}

func (m *customConversationRecordModel) withSession(session sqlx.Session) ConversationRecordModel {
	return NewConversationRecordModel(sqlx.NewSqlConnFromSession(session))
}
