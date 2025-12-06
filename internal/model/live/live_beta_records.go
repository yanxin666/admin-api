package live

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ LiveBetaRecordsModel = (*customLiveBetaRecordsModel)(nil)

type (
	// LiveBetaRecordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveBetaRecordsModel.
	LiveBetaRecordsModel interface {
		liveBetaRecordsModel
		withSession(session sqlx.Session) LiveBetaRecordsModel
		TableName() string
		FindPageByCondition(ctx context.Context, page, limit int64, phoneKeyword, usernameKeyword, companyKeyword string, condition map[string]interface{}) ([]LiveBetaRecords, int64, error)
		UpdateOperateStatus(ctx context.Context, id, status, userId int64) error
	}

	customLiveBetaRecordsModel struct {
		*defaultLiveBetaRecordsModel
	}
)

// NewLiveBetaRecordsModel returns a model for the database table.
func NewLiveBetaRecordsModel(conn sqlx.SqlConn) LiveBetaRecordsModel {
	return &customLiveBetaRecordsModel{
		defaultLiveBetaRecordsModel: newLiveBetaRecordsModel(conn),
	}
}

func (m *customLiveBetaRecordsModel) withSession(session sqlx.Session) LiveBetaRecordsModel {
	return NewLiveBetaRecordsModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLiveBetaRecordsModel) FindPageByCondition(ctx context.Context, page, limit int64, phoneKeyword, usernameKeyword, companyKeyword string, condition map[string]interface{}) ([]LiveBetaRecords, int64, error) {
	var resp []LiveBetaRecords
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if phoneKeyword != "" {
		query = query.Where("`phone` like ?", "%"+cast.ToString(phoneKeyword)+"%")
	}
	if usernameKeyword != "" {
		query = query.Where("`user_name` like ?", "%"+cast.ToString(usernameKeyword)+"%")
	}
	if companyKeyword != "" {
		query = query.Where("`company` like ?", "%"+cast.ToString(companyKeyword)+"%")
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

func (m *customLiveBetaRecordsModel) UpdateOperateStatus(ctx context.Context, id, status, userId int64) error {
	query := fmt.Sprintf("update %s set `status` = ?,`operate_id` = ?, operate_at = ? where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, status, userId, time.Now().Unix(), id)
	return err
}
