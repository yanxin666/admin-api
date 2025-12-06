package behavior

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SessionUserModel = (*customSessionUserModel)(nil)

type (
	// SessionUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSessionUserModel.
	SessionUserModel interface {
		sessionUserModel
		withSession(session sqlx.Session) SessionUserModel
		TableName() string
		FindUidList(ctx context.Context) ([]int64, error)
	}

	customSessionUserModel struct {
		*defaultSessionUserModel
	}
)

// NewSessionUserModel returns a model for the database table.
func NewSessionUserModel(conn sqlx.SqlConn) SessionUserModel {
	return &customSessionUserModel{
		defaultSessionUserModel: newSessionUserModel(conn),
	}
}

func (m *customSessionUserModel) withSession(session sqlx.Session) SessionUserModel {
	return NewSessionUserModel(sqlx.NewSqlConnFromSession(session))
}

// FindUidList 获取用户userId列表 取前10名活跃用户
func (m *customSessionUserModel) FindUidList(ctx context.Context) ([]int64, error) {
	var userIds []int64
	query := fmt.Sprintf("select `user_id` from %s GROUP BY `user_id` ORDER BY MAX(`updated_at`) DESC limit 10", m.table)
	err := m.conn.QueryRowsCtx(ctx, &userIds, query)
	if err != nil {
		return nil, err
	}
	return userIds, nil
}
