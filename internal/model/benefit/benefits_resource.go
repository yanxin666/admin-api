package benefit

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BenefitsResourceModel = (*customBenefitsResourceModel)(nil)

type (
	// BenefitsResourceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBenefitsResourceModel.
	BenefitsResourceModel interface {
		benefitsResourceModel
		withSession(session sqlx.Session) BenefitsResourceModel
		TableName() string
		BatchByIds(ctx context.Context, ids []any) (map[int64]BenefitsResource, error)
		FindResourceAll(ctx context.Context) ([]BenefitsResource, error)
	}

	customBenefitsResourceModel struct {
		*defaultBenefitsResourceModel
	}
)

// NewBenefitsResourceModel returns a model for the database table.
func NewBenefitsResourceModel(conn sqlx.SqlConn) BenefitsResourceModel {
	return &customBenefitsResourceModel{
		defaultBenefitsResourceModel: newBenefitsResourceModel(conn),
	}
}

func (m *customBenefitsResourceModel) withSession(session sqlx.Session) BenefitsResourceModel {
	return NewBenefitsResourceModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBenefitsResourceModel) BatchByIds(ctx context.Context, ids []any) (map[int64]BenefitsResource, error) {
	var list []BenefitsResource
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("id", ids...).
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]BenefitsResource)
	for _, v := range list {
		data[v.Id] = v
	}

	return data, nil
}

func (m *customBenefitsResourceModel) FindResourceAll(ctx context.Context) ([]BenefitsResource, error) {
	var list []BenefitsResource
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		Where("status = 1").
		Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
