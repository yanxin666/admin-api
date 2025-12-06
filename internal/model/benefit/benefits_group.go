package benefit

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BenefitsGroupModel = (*customBenefitsGroupModel)(nil)

type (
	// BenefitsGroupModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBenefitsGroupModel.
	BenefitsGroupModel interface {
		benefitsGroupModel
		withSession(session sqlx.Session) BenefitsGroupModel
		TableName() string
		BatchByGroupIds(ctx context.Context, GroupIds []any) (map[int64]BenefitsGroup, error)
		FindGroupAll(ctx context.Context) ([]BenefitsGroup, error)
	}

	customBenefitsGroupModel struct {
		*defaultBenefitsGroupModel
	}
)

// NewBenefitsGroupModel returns a model for the database table.
func NewBenefitsGroupModel(conn sqlx.SqlConn) BenefitsGroupModel {
	return &customBenefitsGroupModel{
		defaultBenefitsGroupModel: newBenefitsGroupModel(conn),
	}
}

func (m *customBenefitsGroupModel) withSession(session sqlx.Session) BenefitsGroupModel {
	return NewBenefitsGroupModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBenefitsGroupModel) BatchByGroupIds(ctx context.Context, GroupIds []any) (map[int64]BenefitsGroup, error) {
	var resp []BenefitsGroup
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereIn("id", GroupIds...).Order("id desc").Find(&resp)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]BenefitsGroup)
	for _, v := range resp {
		data[v.Id] = v
	}

	return data, nil
}

func (m *customBenefitsGroupModel) FindGroupAll(ctx context.Context) ([]BenefitsGroup, error) {
	var list []BenefitsGroup
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		// Where("status = 1").
		Order("id desc").
		Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
