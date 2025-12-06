package live_course

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatRoomModel = (*customChatRoomModel)(nil)

type (
	// ChatRoomModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatRoomModel.
	ChatRoomModel interface {
		chatRoomModel
		withSession(session sqlx.Session) ChatRoomModel
		TableName() string
		FindSignalRoomByStreamId(ctx context.Context, streamId int64) (*ChatRoom, error)
	}

	customChatRoomModel struct {
		*defaultChatRoomModel
	}
)

// NewChatRoomModel returns a model for the database table.
func NewChatRoomModel(conn sqlx.SqlConn) ChatRoomModel {
	return &customChatRoomModel{
		defaultChatRoomModel: newChatRoomModel(conn),
	}
}

func (m *customChatRoomModel) withSession(session sqlx.Session) ChatRoomModel {
	return NewChatRoomModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customChatRoomModel) FindSignalRoomByStreamId(ctx context.Context, streamId int64) (*ChatRoom, error) {
	query := fmt.Sprintf("select %s from %s where stream_id = ? and type = 1 limit 1", chatRoomRows, m.table)

	var resp ChatRoom
	err := m.conn.QueryRowCtx(ctx, &resp, query, streamId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
