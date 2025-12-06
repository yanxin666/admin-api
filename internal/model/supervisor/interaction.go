package supervisor

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/model/tools"
)

var _ InteractionModel = (*customInteractionModel)(nil)

type (
	// InteractionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInteractionModel.
	InteractionModel interface {
		interactionModel
		withSession(session sqlx.Session) InteractionModel
		TableName() string
		FindAllByCondition(ctx context.Context, nameLike string) ([]Interaction, error)
		UpdateFillFieldsById(ctx context.Context, id int64, data *Interaction) (sql.Result, error)
	}

	customInteractionModel struct {
		*defaultInteractionModel
	}
)

// NewInteractionModel returns a model for the database table.
func NewInteractionModel(conn sqlx.SqlConn) InteractionModel {
	return &customInteractionModel{
		defaultInteractionModel: newInteractionModel(conn),
	}
}

func (m *customInteractionModel) withSession(session sqlx.Session) InteractionModel {
	return NewInteractionModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customInteractionModel) FindAllByCondition(ctx context.Context, nameLike string) ([]Interaction, error) {
	var resp []Interaction
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").Order("id desc")
	if nameLike != "" {
		query = query.Where("`name` like ?", "%"+cast.ToString(nameLike)+"%")
	}
	err := query.Find(&resp)
	return resp, err
}

func (m *customInteractionModel) UpdateFillFieldsById(ctx context.Context, id int64, data *Interaction) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, paramSet)
	args = append(args, id)

	result, err := m.conn.ExecCtx(ctx, query, args...)
	return result, err
}
