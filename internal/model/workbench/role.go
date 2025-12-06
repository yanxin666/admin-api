package workbench

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/define"
)

var _ RoleModel = (*customRoleModel)(nil)

type (
	// RoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleModel.
	RoleModel interface {
		roleModel
		withSession(session sqlx.Session) RoleModel
		TableName() string

		FindAll(ctx context.Context) ([]*Role, error)
		FindEnable(ctx context.Context) ([]*Role, error)
		FindByIds(ctx context.Context, ids string) ([]*Role, error)
		FindSubRole(ctx context.Context, id int64) ([]*Role, error)
	}

	customRoleModel struct {
		*defaultRoleModel
	}
)

// NewRoleModel returns a model for the database table.
func NewRoleModel(conn sqlx.SqlConn) RoleModel {
	return &customRoleModel{
		defaultRoleModel: newRoleModel(conn),
	}
}

func (m *customRoleModel) withSession(session sqlx.Session) RoleModel {
	return NewRoleModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customRoleModel) FindAll(ctx context.Context) ([]*Role, error) {
	query := fmt.Sprintf("SELECT %s FROM %s  WHERE id!=%d AND delete_time IS NULL ORDER BY order_num DESC", roleRows, m.table, define.SysSuperRoleId)
	var resp []*Role
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customRoleModel) FindEnable(ctx context.Context) ([]*Role, error) {
	query := fmt.Sprintf("SELECT %s FROM %s  WHERE id!=%d AND delete_time IS NULL AND status=1 ORDER BY order_num DESC", roleRows, m.table, define.SysSuperRoleId)
	var resp []*Role
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customRoleModel) FindByIds(ctx context.Context, ids string) ([]*Role, error) {
	query := fmt.Sprintf("SELECT %s FROM %s  WHERE id!=%d AND status=1 AND id IN(%s) AND delete_time IS NULL ORDER BY order_num DESC", roleRows, m.table, define.SysSuperRoleId, ids)
	var resp []*Role
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customRoleModel) FindSubRole(ctx context.Context, id int64) ([]*Role, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `parent_id`=%d AND delete_time IS NULL", roleRows, m.table, id)
	var resp []*Role
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
