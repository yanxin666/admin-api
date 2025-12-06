package workbench

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ JobModel = (*customJobModel)(nil)

type (
	// JobModel is an interface to be customized, add more methods here,
	// and implement the added methods in customJobModel.
	JobModel interface {
		jobModel
		withSession(session sqlx.Session) JobModel
		TableName() string

		FindAll(ctx context.Context) ([]*Job, error)
		FindEnable(ctx context.Context) ([]*Job, error)
		FindPage(ctx context.Context, page int64, limit int64) ([]*Job, error)
		FindCount(ctx context.Context) (int64, error)
	}

	customJobModel struct {
		*defaultJobModel
	}
)

// NewJobModel returns a model for the database table.
func NewJobModel(conn sqlx.SqlConn) JobModel {
	return &customJobModel{
		defaultJobModel: newJobModel(conn),
	}
}

func (m *customJobModel) withSession(session sqlx.Session) JobModel {
	return NewJobModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customJobModel) FindAll(ctx context.Context) ([]*Job, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE delete_time IS NULL ORDER BY order_num DESC", jobRows, m.table)
	var resp []*Job
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customJobModel) FindEnable(ctx context.Context) ([]*Job, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE status=1 AND delete_time IS NULL ORDER BY order_num DESC", jobRows, m.table)
	var resp []*Job
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customJobModel) FindPage(ctx context.Context, page int64, limit int64) ([]*Job, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT %s FROM %s WHERE delete_time IS NULL ORDER BY order_num DESC LIMIT %d,%d", jobRows, m.table, offset, limit)
	var resp []*Job
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customJobModel) FindCount(ctx context.Context) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE delete_time IS NULL", m.table)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
