package member

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DoushenVipModel = (*customDoushenVipModel)(nil)

type (
	// DoushenVipModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDoushenVipModel.
	DoushenVipModel interface {
		doushenVipModel
		withSession(session sqlx.Session) DoushenVipModel
		TableName() string
		BatchMapByUserIds(ctx context.Context, userIds []any) (map[int64]DoushenVip, error)
	}

	customDoushenVipModel struct {
		*defaultDoushenVipModel
	}
)

// NewDoushenVipModel returns a model for the database table.
func NewDoushenVipModel(conn sqlx.SqlConn) DoushenVipModel {
	return &customDoushenVipModel{
		defaultDoushenVipModel: newDoushenVipModel(conn),
	}
}

func (m *customDoushenVipModel) withSession(session sqlx.Session) DoushenVipModel {
	return NewDoushenVipModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customDoushenVipModel) BatchMapByUserIds(ctx context.Context, userIds []any) (map[int64]DoushenVip, error) {
	var list []DoushenVip
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("user_id", userIds...).Order("id asc").
		Where("`status` = 1").
		Find(&list)

	data := make(map[int64]DoushenVip)
	for _, v := range list {
		// 是否已经存在当前 userId 的记录，todo 目前只收第一条，后续看策略
		if _, exists := data[v.UserId]; !exists {
			// 如果不存在，则将当前记录存入 map
			data[v.UserId] = v
		}
	}
	return data, err
}
