package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BaseUserModel = (*customBaseUserModel)(nil)

type (
	// BaseUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBaseUserModel.
	BaseUserModel interface {
		baseUserModel
		withSession(session sqlx.Session) BaseUserModel
		TableName() string
		FindOneByPhone(ctx context.Context, phoneEncrypt string) (*BaseUser, error)
		FindPageByCondition(ctx context.Context, page, limit int64, keyword string, condition map[string]interface{}) ([]BaseUser, int64, error)
		FindAllByCondition(ctx context.Context, keyword string, condition map[string]interface{}) ([]BaseUser, error)
		BatchByBaseUserIds(ctx context.Context, baseUserIds []any) (map[int64]BaseUser, error)
	}

	customBaseUserModel struct {
		*defaultBaseUserModel
	}
)

// NewBaseUserModel returns a model for the database table.
func NewBaseUserModel(conn sqlx.SqlConn) BaseUserModel {
	return &customBaseUserModel{
		defaultBaseUserModel: newBaseUserModel(conn),
	}
}

func (m *customBaseUserModel) withSession(session sqlx.Session) BaseUserModel {
	return NewBaseUserModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBaseUserModel) FindOneByPhone(ctx context.Context, phoneEncrypt string) (*BaseUser, error) {
	query := fmt.Sprintf("select %s from %s where `mask_phone` = ? limit 1", baseUserRows, m.table)
	var resp BaseUser
	err := m.conn.QueryRowCtx(ctx, &resp, query, phoneEncrypt)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customBaseUserModel) FindPageByCondition(ctx context.Context, page, limit int64, keyword string, condition map[string]interface{}) ([]BaseUser, int64, error) {
	var resp []BaseUser
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if keyword != "" {
		query = query.Where("`phone` like ?", "%"+cast.ToString(keyword)+"%")
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

func (m *customBaseUserModel) FindAllByCondition(ctx context.Context, keyword string, condition map[string]interface{}) ([]BaseUser, error) {
	var resp []BaseUser
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if keyword != "" {
		query = query.Where("`phone` like ?", "%"+cast.ToString(keyword)+"%")
	}
	err := query.Find(&resp)
	return resp, err
}

func (m *customBaseUserModel) BatchByBaseUserIds(ctx context.Context, baseUserIds []any) (map[int64]BaseUser, error) {
	var list []BaseUser
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("id", baseUserIds...).
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]BaseUser)
	for _, v := range list {
		data[v.Id] = v
	}

	return data, nil
}
