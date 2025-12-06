package behavior

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SessionRecordModel = (*customSessionRecordModel)(nil)

type (
	// SessionRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSessionRecordModel.
	SessionRecordModel interface {
		sessionRecordModel
		withSession(session sqlx.Session) SessionRecordModel
		TableName() string
		FindListBySessionId(ctx context.Context, sessionId int64) ([]*SessionRecord, error)
	}

	customSessionRecordModel struct {
		*defaultSessionRecordModel
	}
)

// NewSessionRecordModel returns a model for the database table.
func NewSessionRecordModel(conn sqlx.SqlConn) SessionRecordModel {
	return &customSessionRecordModel{
		defaultSessionRecordModel: newSessionRecordModel(conn),
	}
}

func (m *customSessionRecordModel) withSession(session sqlx.Session) SessionRecordModel {
	return NewSessionRecordModel(sqlx.NewSqlConnFromSession(session))
}

// FindListBySessionId 根据sessionID批量获取会话记录
func (m *customSessionRecordModel) FindListBySessionId(ctx context.Context, sessionId int64) ([]*SessionRecord, error) {
	query := fmt.Sprintf("select %s from %s where `session_id` = ? ", sessionRecordRows, m.table)
	var resp []*SessionRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query, sessionId)
	switch {
	case err == nil:
		return resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
