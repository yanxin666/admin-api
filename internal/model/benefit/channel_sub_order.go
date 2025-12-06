package benefit

import (
	"context"
	"fmt"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChannelSubOrderModel = (*customChannelSubOrderModel)(nil)

type (
	// ChannelSubOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChannelSubOrderModel.
	ChannelSubOrderModel interface {
		channelSubOrderModel
		withSession(session sqlx.Session) ChannelSubOrderModel
		TableName() string
		FindByChannelOrderId(ctx context.Context, channelOrderId int64) ([]ChannelSubOrder, error)
		BatchByChannelOrderIds(ctx context.Context, channelOrderIds []any) (map[int64][]ChannelSubOrder, error)
		UpdateUserIdWithTx(ctx context.Context, session sqlx.Session, newUserId, oldUserId, channelOrderId, benefitsResourceId int64) error
	}

	customChannelSubOrderModel struct {
		*defaultChannelSubOrderModel
	}
)

// NewChannelSubOrderModel returns a model for the database table.
func NewChannelSubOrderModel(conn sqlx.SqlConn) ChannelSubOrderModel {
	return &customChannelSubOrderModel{
		defaultChannelSubOrderModel: newChannelSubOrderModel(conn),
	}
}

func (m *customChannelSubOrderModel) withSession(session sqlx.Session) ChannelSubOrderModel {
	return NewChannelSubOrderModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customChannelSubOrderModel) FindByChannelOrderId(ctx context.Context, channelOrderId int64) ([]ChannelSubOrder, error) {
	var resp []ChannelSubOrder
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").Where("channel_order_id = ?", channelOrderId).Order("id desc").Find(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *customChannelSubOrderModel) BatchByChannelOrderIds(ctx context.Context, channelOrderIds []any) (map[int64][]ChannelSubOrder, error) {
	var list []ChannelSubOrder
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereIn("channel_order_id", channelOrderIds...).Order("id desc").Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64][]ChannelSubOrder)
	for _, v := range list {
		data[v.ChannelOrderId] = append(data[v.ChannelOrderId], v)
	}

	return data, nil
}

// UpdateUserIdWithTx 更改用户权益所属用户ID，使用事务
func (m *customChannelSubOrderModel) UpdateUserIdWithTx(ctx context.Context, session sqlx.Session, newUserId, oldUserId, channelOrderId, benefitsResourceId int64) error {
	query := fmt.Sprintf("update %s set `user_id` = ? where `user_id` = ? and `channel_order_id` = ? and `benefits_resource_id` = ?", m.table)
	_, err := session.ExecCtx(ctx, query, newUserId, oldUserId, channelOrderId, benefitsResourceId)
	if err != nil {
		return err
	}
	return nil
}
