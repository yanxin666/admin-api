package workbench

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DictionaryModel = (*customDictionaryModel)(nil)

type (
	// DictionaryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDictionaryModel.
	DictionaryModel interface {
		dictionaryModel
		withSession(session sqlx.Session) DictionaryModel
		TableName() string

		FindDictionaryList(ctx context.Context) ([]*Dictionary, error)
		FindPageByParentId(ctx context.Context, id int64, page int64, limit int64) ([]*Dictionary, error)
		FindCountByParentId(ctx context.Context, id int64) (int64, error)
	}

	customDictionaryModel struct {
		*defaultDictionaryModel
	}
)

// NewDictionaryModel returns a model for the database table.
func NewDictionaryModel(conn sqlx.SqlConn) DictionaryModel {
	return &customDictionaryModel{
		defaultDictionaryModel: newDictionaryModel(conn),
	}
}

func (m *customDictionaryModel) withSession(session sqlx.Session) DictionaryModel {
	return NewDictionaryModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customDictionaryModel) FindDictionaryList(ctx context.Context) ([]*Dictionary, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE parent_id=0 AND delete_time IS NULL ORDER BY order_num DESC", dictionaryRows, m.table)
	var resp []*Dictionary
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customDictionaryModel) FindPageByParentId(ctx context.Context, id int64, page int64, limit int64) ([]*Dictionary, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT %s FROM %s WHERE parent_id=%d AND delete_time IS NULL ORDER BY order_num DESC LIMIT %d,%d", dictionaryRows, m.table, id, offset, limit)
	var resp []*Dictionary
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customDictionaryModel) FindCountByParentId(ctx context.Context, id int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE parent_id=%d AND delete_time IS NULL", m.table, id)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
