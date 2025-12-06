package system

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DictWhiteModel = (*customDictWhiteModel)(nil)

type (
	// DictWhiteModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDictWhiteModel.
	DictWhiteModel interface {
		dictWhiteModel
		withSession(session sqlx.Session) DictWhiteModel
		TableName() string

		FindPageByCondition(ctx context.Context, page, limit int64, key, remark string, condition map[string]interface{}) ([]DictWhite, int64, error)
		FindByKey(ctx context.Context, key string) (*DictWhite, error)
	}

	customDictWhiteModel struct {
		*defaultDictWhiteModel
	}
)

// NewDictWhiteModel returns a model for the database table.
func NewDictWhiteModel(conn sqlx.SqlConn) DictWhiteModel {
	return &customDictWhiteModel{
		defaultDictWhiteModel: newDictWhiteModel(conn),
	}
}

func (m *customDictWhiteModel) withSession(session sqlx.Session) DictWhiteModel {
	return NewDictWhiteModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customDictWhiteModel) FindPageByCondition(ctx context.Context, page, limit int64, key, remark string, condition map[string]interface{}) ([]DictWhite, int64, error) {
	var resp []DictWhite
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if key != "" {
		query = query.Where("`key` like ?", "%"+cast.ToString(key)+"%")
	}
	if remark != "" {
		query = query.Where("`remark` like ?", "%"+cast.ToString(remark)+"%")
	}
	total, err := query.Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}

func (m *customDictWhiteModel) FindByKey(ctx context.Context, key string) (*DictWhite, error) {
	query := fmt.Sprintf("select %s from %s where `key` = ? limit 1", dictWhiteRows, m.table)

	var resp DictWhite
	err := m.conn.QueryRowCtx(ctx, &resp, query, key)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
