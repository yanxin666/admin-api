package benefit

import (
	"context"
	"fmt"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChannelOrderModel = (*customChannelOrderModel)(nil)

type (
	// ChannelOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChannelOrderModel.
	ChannelOrderModel interface {
		channelOrderModel
		withSession(session sqlx.Session) ChannelOrderModel
		TableName() string
		FindPageByCondition(ctx context.Context, page, limit int64, keywordOrderNo string, condition map[string]interface{}) ([]ChannelOrder, int64, error)
		FindPageByConditionCursor(ctx context.Context, lastId, limit int64, keywordOrderNo string, condition map[string]interface{}) ([]ChannelOrder, int64, error)
		FindOneByOrderNo(ctx context.Context, orderNo string) (*ChannelOrder, error)
	}

	customChannelOrderModel struct {
		*defaultChannelOrderModel
	}
)

// NewChannelOrderModel returns a model for the database table.
func NewChannelOrderModel(conn sqlx.SqlConn) ChannelOrderModel {
	return &customChannelOrderModel{
		defaultChannelOrderModel: newChannelOrderModel(conn),
	}
}

func (m *customChannelOrderModel) withSession(session sqlx.Session) ChannelOrderModel {
	return NewChannelOrderModel(sqlx.NewSqlConnFromSession(session))
}

// FindPageByCondition 分页查询
func (m *customChannelOrderModel) FindPageByCondition(ctx context.Context, page, limit int64, keywordOrderNo string, condition map[string]interface{}) ([]ChannelOrder, int64, error) {
	var resp []ChannelOrder
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if keywordOrderNo != "" {
		query = query.Where("`order_no` like ?", cast.ToString(keywordOrderNo)+"%")
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

// FindPageByConditionCursor 游标分页查询
func (m *customChannelOrderModel) FindPageByConditionCursor(ctx context.Context, lastId, limit int64, keywordOrderNo string, condition map[string]interface{}) ([]ChannelOrder, int64, error) {
	var resp []ChannelOrder
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	query = query.UsePageCursor("id > ?", lastId)

	if keywordOrderNo != "" {
		query = query.Where("`order_no` like ?", cast.ToString(keywordOrderNo)+"%")
	}
	total, err := query.Page(int(lastId), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}

func (m *customChannelOrderModel) FindOneByOrderNo(ctx context.Context, orderNo string) (*ChannelOrder, error) {
	query := fmt.Sprintf("select %s from %s where order_no = ? limit 1", channelOrderRows, m.table)
	var resp ChannelOrder
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderNo)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
