package benefit

import (
	"context"
	"database/sql"
	"fmt"
	"muse-admin/internal/model/tools"

	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ BenefitsUserRecordModel = (*customBenefitsUserRecordModel)(nil)

type (
	// BenefitsUserRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBenefitsUserRecordModel.
	BenefitsUserRecordModel interface {
		benefitsUserRecordModel
		withSession(session sqlx.Session) BenefitsUserRecordModel
		TableName() string
		FindAllByUserId(ctx context.Context, userId int64) ([]BenefitsUserRecord, error)
		BatchByUserIds(ctx context.Context, userIds []any) ([]BenefitsUserRecord, error)
		GroupResourceByUserId(ctx context.Context, userId int64) ([]GroupResourceByUserIdResult, error)
		FindPageByCondition(ctx context.Context, page, limit int64, condition map[string]interface{}) ([]BenefitsUserRecord, int64, error)
		FindOneByUserIdAndBenefitId(ctx context.Context, userId, benefitId int64) (*BenefitsUserRecord, error)
		UpdateFillFieldsByUserIdAndResourceId(ctx context.Context, userId, benefitsResourceId int64, data *BenefitsUserRecord) (sql.Result, error)
		UpdateUserIdWithTx(ctx context.Context, session sqlx.Session, targetUserId, oldUserId, benefitId int64) error
	}

	GroupResourceByUserIdResult struct {
		BenefitsResourceId int64 `db:"benefits_resource_id"` // 权益资源ID
		GrantCount         int64 `db:"grant_count"`          // 发放数量
		ReclaimCount       int64 `db:"reclaim_count"`        // 回收数量
	}

	customBenefitsUserRecordModel struct {
		*defaultBenefitsUserRecordModel
	}
)

// NewBenefitsUserRecordModel returns a model for the database table.
func NewBenefitsUserRecordModel(conn sqlx.SqlConn) BenefitsUserRecordModel {
	return &customBenefitsUserRecordModel{
		defaultBenefitsUserRecordModel: newBenefitsUserRecordModel(conn),
	}
}

func (m *customBenefitsUserRecordModel) withSession(session sqlx.Session) BenefitsUserRecordModel {
	return NewBenefitsUserRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customBenefitsUserRecordModel) FindAllByUserId(ctx context.Context, userId int64) ([]BenefitsUserRecord, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", benefitsUserRecordRows, m.table)
	var resp []BenefitsUserRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	return resp, err
}

func (m *customBenefitsUserRecordModel) FindOneByUserIdAndBenefitId(ctx context.Context, userId, benefitId int64) (*BenefitsUserRecord, error) {
	query := fmt.Sprintf("select %s from %s where user_id = ? and user_benefits_id = ? limit 1", benefitsUserRecordRows, m.table)
	var resp BenefitsUserRecord
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId, benefitId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customBenefitsUserRecordModel) BatchByUserIds(ctx context.Context, userIds []any) ([]BenefitsUserRecord, error) {
	var list []BenefitsUserRecord
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").From(m.TableName()).WhereIn("user_id", userIds...).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customBenefitsUserRecordModel) GroupResourceByUserId(ctx context.Context, userId int64) ([]GroupResourceByUserIdResult, error) {
	// SELECT
	//  benefits_resource_id,
	// 	COUNT(DISTINCT CASE WHEN type = 1 THEN order_no END) AS grant_count,
	// 	COUNT(DISTINCT CASE WHEN type = 2 THEN order_no END) AS reclaim_count
	// FROM
	// bt_benefits_user_record where user_id = 2
	// GROUP BY
	// benefits_resource_id;

	query := fmt.Sprintf("select benefits_resource_id," +
		"COUNT(DISTINCT CASE WHEN type = 1 THEN order_no END) AS grant_count," +
		"COUNT(DISTINCT CASE WHEN type = 2 THEN order_no END) AS reclaim_count " +
		"FROM bt_benefits_user_record where user_id = ? GROUP BY benefits_resource_id")
	var resp []GroupResourceByUserIdResult
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	return resp, err
}

func (m *customBenefitsUserRecordModel) FindPageByCondition(ctx context.Context, page, limit int64, condition map[string]interface{}) ([]BenefitsUserRecord, int64, error) {
	var resp []BenefitsUserRecord
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	total, err := query.Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}

// UpdateFillFieldsByUserIdAndResourceId 支持更新自定义struct
func (m *customBenefitsUserRecordModel) UpdateFillFieldsByUserIdAndResourceId(ctx context.Context, userId, benefitsResourceId int64, data *BenefitsUserRecord) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE user_id = ? and benefits_resource_id = ? ", m.table, paramSet)
	args = append(args, userId, benefitsResourceId)
	return m.conn.ExecCtx(ctx, query, args...)
}

// UpdateUserIdWithTx 更改用户权益所属用户ID，使用事务
func (m *customBenefitsUserRecordModel) UpdateUserIdWithTx(ctx context.Context, session sqlx.Session, targetUserId, oldUserId, benefitId int64) error {
	query := fmt.Sprintf("update %s set `user_id` = ? where `user_id` = ? and `user_benefits_id` = ?", m.table)
	_, err := session.ExecCtx(ctx, query, targetUserId, oldUserId, benefitId)
	if err != nil {
		return err
	}
	return nil
}
