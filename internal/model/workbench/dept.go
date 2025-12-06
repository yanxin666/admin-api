package workbench

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DeptModel = (*customDeptModel)(nil)

type (
	// DeptModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDeptModel.
	DeptModel interface {
		deptModel
		TableName() string

		FindAll(ctx context.Context) ([]*Dept, error)
		FindSubDept(ctx context.Context, id int64) ([]*Dept, error)
		FindEnable(ctx context.Context) ([]*Dept, error)
		withSession(session sqlx.Session) DeptModel
	}

	customDeptModel struct {
		*defaultDeptModel
	}
)

// NewDeptModel returns a model for the database table.
func NewDeptModel(conn sqlx.SqlConn) DeptModel {
	return &customDeptModel{
		defaultDeptModel: newDeptModel(conn),
	}
}

func (m *customDeptModel) withSession(session sqlx.Session) DeptModel {
	return NewDeptModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customDeptModel) FindAll(ctx context.Context) ([]*Dept, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE delete_time IS NULL ORDER BY order_num DESC", deptRows, m.table)
	var resp []*Dept
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customDeptModel) FindEnable(ctx context.Context) ([]*Dept, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE status=1 AND delete_time IS NULL ORDER BY order_num DESC", deptRows, m.table)
	var resp []*Dept
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customDeptModel) FindSubDept(ctx context.Context, id int64) ([]*Dept, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `parent_id`=%d AND delete_time IS NULL", deptRows, m.table, id)
	var resp []*Dept
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
