package benefit

import (
	"context"
	"database/sql"
	"fmt"
	"muse-admin/internal/model/tools"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserBenefitsDetailModel = (*customUserBenefitsDetailModel)(nil)

type (
	// UserBenefitsDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserBenefitsDetailModel.
	UserBenefitsDetailModel interface {
		userBenefitsDetailModel
		withSession(session sqlx.Session) UserBenefitsDetailModel
		TableName() string
		FindPageByCondition(ctx context.Context, page int64, limit int64, condition map[string]interface{}, fromTime, endTime string) ([]UserBenefitsDetail, int64, error)
		FindPageByConditionAndUserIds(ctx context.Context, page int64, limit int64, userIds []any, condition map[string]interface{}, fromTime, endTime string) ([]UserBenefitsDetail, int64, error)
		FindByIds(ctx context.Context, ids []any) ([]UserBenefitsDetail, error)
		UpdateFillFieldsById(ctx context.Context, id int64, data *UserBenefitsDetail) (sql.Result, error)
		UpdateUserIdWithTx(ctx context.Context, session sqlx.Session, userId, benefitId int64) error
	}

	customUserBenefitsDetailModel struct {
		*defaultUserBenefitsDetailModel
	}
)

// NewUserBenefitsDetailModel returns a model for the database table.
func NewUserBenefitsDetailModel(conn sqlx.SqlConn) UserBenefitsDetailModel {
	return &customUserBenefitsDetailModel{
		defaultUserBenefitsDetailModel: newUserBenefitsDetailModel(conn),
	}
}

func (m *customUserBenefitsDetailModel) withSession(session sqlx.Session) UserBenefitsDetailModel {
	return NewUserBenefitsDetailModel(sqlx.NewSqlConnFromSession(session))
}

// FindPageByCondition 分页查询
func (m *customUserBenefitsDetailModel) FindPageByCondition(ctx context.Context, page int64, limit int64, condition map[string]interface{}, fromTime, endTime string) ([]UserBenefitsDetail, int64, error) {
	var resp []UserBenefitsDetail
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*")
	if fromTime != "" {
		query = query.Where("from_time > ?", fromTime)
	}
	if endTime != "" {
		query = query.Where("end_time < ?", endTime)
	}
	total, err := query.WhereMap(condition).Order("id desc").Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}

func (m *customUserBenefitsDetailModel) FindPageByConditionAndUserIds(ctx context.Context, page int64, limit int64, userIds []any, condition map[string]interface{}, fromTime, endTime string) ([]UserBenefitsDetail, int64, error) {
	var resp []UserBenefitsDetail
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*")
	if len(userIds) > 0 {
		query = query.WhereIn("user_id", userIds...)
	}
	if fromTime != "" {
		query = query.Where("from_time > ?", fromTime)
	}
	if endTime != "" {
		query = query.Where("end_time < ?", endTime)
	}
	total, err := query.WhereMap(condition).Order("id desc").Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}

// FindByIds 根据Ids查询多条记录
func (m *customUserBenefitsDetailModel) FindByIds(ctx context.Context, ids []any) ([]UserBenefitsDetail, error) {
	var resp []UserBenefitsDetail
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereIn("id", ids).Order("id desc").Find(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateFillFieldsById 根据userId 进行更新 支持自定义struct
func (m *customUserBenefitsDetailModel) UpdateFillFieldsById(ctx context.Context, id int64, data *UserBenefitsDetail) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, paramSet)
	args = append(args, id)
	return m.conn.ExecCtx(ctx, query, args...)
}

// UpdateUserIdWithTx 更改用户权益所属用户ID，使用事务
func (m *customUserBenefitsDetailModel) UpdateUserIdWithTx(ctx context.Context, session sqlx.Session, userId, benefitId int64) error {
	query := fmt.Sprintf("update %s set `user_id` = ? where `id` = ?", m.table)
	_, err := session.ExecCtx(ctx, query, userId, benefitId)
	if err != nil {
		return err
	}
	return nil
}
