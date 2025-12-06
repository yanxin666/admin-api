package code

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RedeemCodeModel = (*customRedeemCodeModel)(nil)

type (
	// RedeemCodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRedeemCodeModel.
	RedeemCodeModel interface {
		redeemCodeModel
		withSession(session sqlx.Session) RedeemCodeModel
		TableName() string
		Inserts(ctx context.Context, redeemCode []*RedeemCode) error
		FindPageByCondition(ctx context.Context, page int64, limit int64, condition map[string]interface{}) ([]RedeemCode, int64, error)
	}

	customRedeemCodeModel struct {
		*defaultRedeemCodeModel
	}
)

// NewRedeemCodeModel returns a model for the database table.
func NewRedeemCodeModel(conn sqlx.SqlConn) RedeemCodeModel {
	return &customRedeemCodeModel{
		defaultRedeemCodeModel: newRedeemCodeModel(conn),
	}
}

func (m *customRedeemCodeModel) withSession(session sqlx.Session) RedeemCodeModel {
	return NewRedeemCodeModel(sqlx.NewSqlConnFromSession(session))
}

const batchSize = 3000 // 每批插入的数量
func (m *customRedeemCodeModel) Inserts(ctx context.Context, redeemCode []*RedeemCode) error {
	total := len(redeemCode)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}

		inserts := sqld.NewModel(ctx, m.conn, m.table).Inserts(redeemCodeRowsExpectAutoSet)
		for _, v := range redeemCode[i:end] {
			inserts.Append(v.Code, v.Status, v.BenefitsGroupId, v.Batch, v.ValidDate, v.Source, v.Remark)
		}

		_, err := inserts.Execute()
		if err != nil {
			return err // 如果插入出错，直接返回
		}
	}
	return nil
}

// FindPageByCondition 分页查询
func (m *customRedeemCodeModel) FindPageByCondition(ctx context.Context, page int64, limit int64, condition map[string]interface{}) ([]RedeemCode, int64, error) {
	var resp []RedeemCode
	db := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	total, err := db.Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return resp, 0, nil
	}

	return resp, total, nil
}
