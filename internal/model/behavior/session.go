package behavior

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SessionModel = (*customSessionModel)(nil)

type (
	// SessionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSessionModel.
	SessionModel interface {
		sessionModel
		withSession(session sqlx.Session) SessionModel
		TableName() string
		FindListByUserId(ctx context.Context, userId int64) ([]*Session, error)
	}

	customSessionModel struct {
		*defaultSessionModel
	}
)

// NewSessionModel returns a model for the database table.
func NewSessionModel(conn sqlx.SqlConn) SessionModel {
	return &customSessionModel{
		defaultSessionModel: newSessionModel(conn),
	}
}

func (m *customSessionModel) withSession(session sqlx.Session) SessionModel {
	return NewSessionModel(sqlx.NewSqlConnFromSession(session))
}

// FindListByUserId 根据用户id查询会话列表
func (m *customSessionModel) FindListByUserId(ctx context.Context, userId int64) ([]*Session, error) {
	var result []*Session
	query := sqld.NewModel(ctx, m.conn, m.table).
		Select("us.*").
		From(fmt.Sprintf("%s us", m.table)).
		Join("LEFT JOIN us_session_user usu ON us.id = usu.session_id").
		Where("usu.user_id = ?", userId).
		Order("us.id desc")
	err := query.Find(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
