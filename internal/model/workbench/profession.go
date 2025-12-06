package workbench

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ProfessionModel = (*customProfessionModel)(nil)

type (
	// ProfessionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customProfessionModel.
	ProfessionModel interface {
		professionModel
		withSession(session sqlx.Session) ProfessionModel
		TableName() string

		FindAll(ctx context.Context) ([]*Profession, error)
		FindEnable(ctx context.Context) ([]*Profession, error)
		FindCount(ctx context.Context) (int64, error)
		FindPage(ctx context.Context, page int64, limit int64) ([]*Profession, error)
	}

	customProfessionModel struct {
		*defaultProfessionModel
	}
)

// NewProfessionModel returns a model for the database table.
func NewProfessionModel(conn sqlx.SqlConn) ProfessionModel {
	return &customProfessionModel{
		defaultProfessionModel: newProfessionModel(conn),
	}
}

func (m *customProfessionModel) withSession(session sqlx.Session) ProfessionModel {
	return NewProfessionModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customProfessionModel) FindAll(ctx context.Context) ([]*Profession, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE delete_time IS NULL ORDER BY order_num DESC", professionRows, m.table)
	var resp []*Profession
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customProfessionModel) FindEnable(ctx context.Context) ([]*Profession, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE status=1 AND delete_time IS NULL ORDER BY order_num DESC", professionRows, m.table)
	var resp []*Profession
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customProfessionModel) FindPage(ctx context.Context, page int64, limit int64) ([]*Profession, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT %s FROM %s WHERE delete_time IS NULL ORDER BY order_num DESC LIMIT %d,%d", professionRows, m.table, offset, limit)
	var resp []*Profession
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customProfessionModel) FindCount(ctx context.Context) (int64, error) {
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
