package workbench

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ LogModel = (*customLogModel)(nil)

type LoginLog struct {
	Id         int64     `db:"id"`          // 编号
	Account    string    `db:"account"`     // 操作账号
	Ip         string    `db:"ip"`          // ip
	Uri        string    `db:"uri"`         // 请求路径
	Status     int64     `db:"status"`      // 0.失败 1.成功
	CreateTime time.Time `db:"create_time"` // 创建时间
}

type (
	// LogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLogModel.
	LogModel interface {
		logModel
		withSession(session sqlx.Session) LogModel
		TableName() string

		FindPage(ctx context.Context, t int64, page int64, limit int64) ([]*LoginLog, error)
		FindCount(ctx context.Context, t int64) (int64, error)
	}

	customLogModel struct {
		*defaultLogModel
	}
)

// NewLogModel returns a model for the database table.
func NewLogModel(conn sqlx.SqlConn) LogModel {
	return &customLogModel{
		defaultLogModel: newLogModel(conn),
	}
}

func (m *customLogModel) withSession(session sqlx.Session) LogModel {
	return NewLogModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLogModel) FindPage(ctx context.Context, t int64, page int64, limit int64) ([]*LoginLog, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT l.id,IFNULL(u.account,'NULL') as account,l.ip,l.uri,l.status,l.create_time FROM (SELECT * FROM wk_log WHERE type=%d ORDER BY id DESC LIMIT %d,%d) l LEFT JOIN wk_user u ON l.user_id=u.id", t, offset, limit)
	var resp []*LoginLog
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customLogModel) FindCount(ctx context.Context, t int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE type=%d", m.table, t)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
