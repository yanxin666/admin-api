package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/pkg/errs"
)

var _ WhiteListModel = (*customWhiteListModel)(nil)

type (
	// WhiteListModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWhiteListModel.
	WhiteListModel interface {
		whiteListModel
		withSession(session sqlx.Session) WhiteListModel
		TableName() string

		FindPageByCondition(ctx context.Context, page, limit int64, phone, startTime, endTime, remark string, condition map[string]interface{}) ([]WhiteList, int64, error)
		BatchPhonesByUnique(ctx context.Context, phones []string, product, source int64) (map[string]WhiteList, error)
		BatchInsert(ctx context.Context, list []*WhiteList) error
		UpdateColumn(ctx context.Context, id []any, data *WhiteList) (int64, error)
	}

	customWhiteListModel struct {
		*defaultWhiteListModel
	}
)

// NewWhiteListModel returns a model for the database table.
func NewWhiteListModel(conn sqlx.SqlConn) WhiteListModel {
	return &customWhiteListModel{
		defaultWhiteListModel: newWhiteListModel(conn),
	}
}

func (m *customWhiteListModel) withSession(session sqlx.Session) WhiteListModel {
	return NewWhiteListModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customWhiteListModel) FindPageByCondition(ctx context.Context, page, limit int64, phone, startTime, endTime, remark string, condition map[string]interface{}) ([]WhiteList, int64, error) {
	var resp []WhiteList
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if phone != "" {
		query = query.Where("`phone` like ?", "%"+cast.ToString(phone)+"%")
	}
	if startTime != "" {
		query = query.Where("`start_time` >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("`end_time` <= ?", endTime)
	}
	if remark != "" {
		query = query.Where("`remark` like ?", "%"+cast.ToString(remark)+"%")
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

func (m *customWhiteListModel) BatchPhonesByUnique(ctx context.Context, phones []string, product, source int64) (map[string]WhiteList, error) {
	var batch []any
	for _, v := range phones {
		batch = append(batch, v)
	}

	var list []WhiteList
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		Where("product = ?", product).
		Where("source = ?", source).
		WhereIn("phone", batch...).
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[string]WhiteList)
	for _, v := range list {
		data[v.Phone] = v
	}

	return data, nil
}

func (m *customWhiteListModel) BatchInsert(ctx context.Context, list []*WhiteList) error {
	inserts := sqld.NewModel(ctx, m.conn, m.table).Inserts(whiteListRowsExpectAutoSet)
	for _, v := range list {
		inserts.StrongAppend(v)
	}
	_, err := inserts.Execute()
	return err
}

// UpdateColumn 更新指定列
func (m *customWhiteListModel) UpdateColumn(ctx context.Context, id []any, data *WhiteList) (int64, error) {
	result, err := sqld.NewModel(ctx, m.conn, m.table).
		WhereIn("id", id...).
		UpdateColumn("`start_time`, `end_time`, `status`, `remark`", data.StartTime, data.EndTime, data.Status, data.Remark)
	if err != nil {
		return 0, errs.WithErr(err)
	}

	rid, err := result.RowsAffected()
	if err != nil {
		return 0, errs.WithErr(err)
	}

	return rid, err
}
