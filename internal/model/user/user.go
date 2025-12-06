package user

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/model/tools"
	"strings"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		withSession(session sqlx.Session) UserModel
		TableName() string
		FindOneByBaseUserId(ctx context.Context, baseId int64) ([]*User, error)
		UpdateUserById(ctx context.Context, uid int64, data *User) (sql.Result, error)
		UpdateById(ctx context.Context, id int64, updateFields map[string]interface{}) error
		BatchMapByBaseUserIds(ctx context.Context, baseUserIds []any) (map[int64][]User, error)
		BatchListByBaseUserIds(ctx context.Context, baseUserIds []any) ([]User, error)
		BatchByUserIds(ctx context.Context, userIds []any) (map[int64]UserJoinBaseUser, error)
	}

	customUserModel struct {
		*defaultUserModel
	}

	UserJoinBaseUser struct {
		User
		UserNo    string `db:"user_no"`     // 用户编号[只做展示用途]
		MaskPhone string `db:"mask_phone"`  // 加密手机号
		Phone     string `db:"phone"`       // 脱敏手机号
		Product   int64  `db:"product"`     // 所属产品 0.辞源
		Source    int64  `db:"source"`      // 首次来源 100.提审标记 0.APP 1.豆伴匠 2.风的颜色 1001.听力熊 1002.微软
		Status    int64  `db:"base_status"` // 状态 0.正常 1.冻结 2.注销
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUserModel) FindOneByBaseUserId(ctx context.Context, baseId int64) ([]*User, error) {
	query := fmt.Sprintf("select %s from %s where `base_user_id` = ? and `status` = 1 order by id asc", userRows, m.table)
	var resp []*User
	err := m.conn.QueryRowsCtx(ctx, &resp, query, baseId)
	return resp, err
}

// UpdateUserById 根据id 更新用户信息 支持自定义struct
func (m *customUserModel) UpdateUserById(ctx context.Context, uid int64, data *User) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)
	// 构建 UPDATE 语句的 SQL 语句和参数
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, paramSet)
	args = append(args, uid)
	return m.conn.ExecCtx(ctx, query, args...)
}

func (m *customUserModel) UpdateById(ctx context.Context, id int64, updateFields map[string]interface{}) error {
	// 构建更新的SQL语句
	var placeholders []string
	var values []interface{}

	for field, value := range updateFields {
		placeholders = append(placeholders, fmt.Sprintf("%s=?", field))
		values = append(values, value)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE `id`=?", m.table, strings.Join(placeholders, ", "))
	values = append(values, id)

	_, err := m.conn.ExecCtx(ctx, query, values...)
	return err
}

func (m *customUserModel) BatchMapByBaseUserIds(ctx context.Context, baseUserIds []any) (map[int64][]User, error) {
	var list []User
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("base_user_id", baseUserIds...).Order("id asc").
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64][]User)
	for _, v := range list {
		data[v.BaseUserId] = append(data[v.BaseUserId], v)
	}

	return data, nil
}

func (m *customUserModel) BatchListByBaseUserIds(ctx context.Context, baseUserIds []any) ([]User, error) {
	var list []User
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("base_user_id", baseUserIds...).Order("id asc").
		Find(&list)

	return list, err
}

func (m *customUserModel) BatchByUserIds(ctx context.Context, userIds []any) (map[int64]UserJoinBaseUser, error) {
	var list []UserJoinBaseUser
	err := sqld.NewModel(ctx, m.conn, m.table).Select("u.*,bu.id as base_user_id,bu.user_no,bu.mask_phone,bu.phone,bu.product,bu.source,bu.`status` as base_status").From(m.TableName()+" u").
		Join("left join us_base_user as bu ON u.base_user_id = bu.id").
		WhereIn("u.id", userIds...).
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]UserJoinBaseUser)
	for _, v := range list {
		data[v.Id] = v
	}

	return data, nil
}
