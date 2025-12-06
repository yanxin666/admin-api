package member

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/model/tools"
)

var _ UserInfoModel = (*customUserInfoModel)(nil)

type (
	// UserInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserInfoModel.
	UserInfoModel interface {
		userInfoModel
		withSession(session sqlx.Session) UserInfoModel
		TableName() string
		FindPage(ctx context.Context, page int64, limit int64) ([]UserInfo, error)
		FindOneByUserId(ctx context.Context, userId int64) (*UserInfo, error)
		BatchByUserIds(ctx context.Context, ids []any) (map[int64]UserInfo, error)
		UpdateFillFieldsById(ctx context.Context, uid int64, data *UserInfo) (sql.Result, error)
	}

	customUserInfoModel struct {
		*defaultUserInfoModel
	}
)

// NewUserInfoModel returns a model for the database table.
func NewUserInfoModel(conn sqlx.SqlConn) UserInfoModel {
	return &customUserInfoModel{
		defaultUserInfoModel: newUserInfoModel(conn),
	}
}

func (m *customUserInfoModel) withSession(session sqlx.Session) UserInfoModel {
	return NewUserInfoModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUserInfoModel) FindPage(ctx context.Context, page int64, limit int64) ([]UserInfo, error) {
	offset := (page - 1) * limit
	query := fmt.Sprintf("select %s from %s ORDER BY id DESC LIMIT %d,%d", userInfoRows, m.table, offset, limit)
	var resp []UserInfo
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	return resp, err
}

func (m *customUserInfoModel) FindOneByUserId(ctx context.Context, userId int64) (*UserInfo, error) {
	var resp UserInfo
	query := fmt.Sprintf("select %s from %s where `base_user_id` = ? limit 1", userInfoRows, m.table)
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUserInfoModel) BatchByUserIds(ctx context.Context, ids []any) (map[int64]UserInfo, error) {
	var list []UserInfo
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("id", ids...).
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]UserInfo, len(list))
	for _, v := range list {
		data[v.BaseUserId] = v
	}

	return data, nil
}

func (m *customUserInfoModel) UpdateFillFieldsById(ctx context.Context, uid int64, data *UserInfo) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE base_user_id = ?", m.table, paramSet)
	args = append(args, uid)

	result, err := m.conn.ExecCtx(ctx, query, args...)
	return result, err
}
