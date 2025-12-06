package system

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SmsRecordModel = (*customSmsRecordModel)(nil)

type (
	// SmsRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSmsRecordModel.
	SmsRecordModel interface {
		smsRecordModel
		withSession(session sqlx.Session) SmsRecordModel
		TableName() string
		FindTempAndPhone(ctx context.Context, temp, phone string) (*SmsRecord, error)
		FindByPhones(ctx context.Context, temp string, phoneArr []any) ([]string, error)
	}

	customSmsRecordModel struct {
		*defaultSmsRecordModel
	}
)

// NewSmsRecordModel returns a model for the database table.
func NewSmsRecordModel(conn sqlx.SqlConn) SmsRecordModel {
	return &customSmsRecordModel{
		defaultSmsRecordModel: newSmsRecordModel(conn),
	}
}

func (m *customSmsRecordModel) withSession(session sqlx.Session) SmsRecordModel {
	return NewSmsRecordModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSmsRecordModel) FindTempAndPhone(ctx context.Context, temp, phone string) (*SmsRecord, error) {
	query := fmt.Sprintf("select %s from %s where `template` = ? and `phone` = ? limit 1", smsRecordRows, m.table)
	var resp SmsRecord
	err := m.conn.QueryRowCtx(ctx, &resp, query, temp, phone)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

// FindByPhones 根据phones查询多条记录
func (m *customSmsRecordModel) FindByPhones(ctx context.Context, temp string, phoneArr []any) (arr []string, err error) {
	var resp []SmsRecord
	err = sqld.NewModel(ctx, m.conn, m.table).Select("*").Where("template = ?", temp).Where("code = ?", "ok").WhereIn("phone", phoneArr...).Order("id desc").Find(&resp)
	if err != nil {
		return nil, err
	}

	for _, v := range resp {
		arr = append(arr, v.Phone)
	}

	return arr, nil
}
