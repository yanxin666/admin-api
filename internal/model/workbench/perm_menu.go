package workbench

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PermMenuModel = (*customPermMenuModel)(nil)

type (
	// PermMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPermMenuModel.
	PermMenuModel interface {
		permMenuModel
		withSession(session sqlx.Session) PermMenuModel
		TableName() string

		FindByIds(ctx context.Context, ids string) ([]*PermMenu, error)
		FindCountByParentId(ctx context.Context, id int64) (int64, error)
		FindAll(ctx context.Context) ([]*PermMenu, error)
		FindAllToSort(ctx context.Context) ([]*PermMenu, error)
		FindSubPermMenu(ctx context.Context, id int64) ([]*PermMenu, error)
	}

	customPermMenuModel struct {
		*defaultPermMenuModel
	}
)

// NewPermMenuModel returns a model for the database table.
func NewPermMenuModel(conn sqlx.SqlConn) PermMenuModel {
	return &customPermMenuModel{
		defaultPermMenuModel: newPermMenuModel(conn),
	}
}

func (m *customPermMenuModel) withSession(session sqlx.Session) PermMenuModel {
	return NewPermMenuModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customPermMenuModel) FindByIds(ctx context.Context, ids string) ([]*PermMenu, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `id` IN(%s) AND delete_time IS NULL", permMenuRows, m.table, ids)
	var resp []*PermMenu
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customPermMenuModel) FindCountByParentId(ctx context.Context, id int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE `parent_id`=%d AND delete_time IS NULL", m.table, id)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *customPermMenuModel) FindAll(ctx context.Context) ([]*PermMenu, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE delete_time IS NULL ORDER BY order_num DESC", permMenuRows, m.table)
	var resp []*PermMenu
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customPermMenuModel) FindAllToSort(ctx context.Context) ([]*PermMenu, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE delete_time IS NULL ORDER BY order_num DESC,id asc", permMenuRows, m.table)
	var resp []*PermMenu
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customPermMenuModel) FindSubPermMenu(ctx context.Context, id int64) ([]*PermMenu, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `parent_id`=%d AND delete_time IS NULL", permMenuRows, m.table, id)
	var resp []*PermMenu
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
