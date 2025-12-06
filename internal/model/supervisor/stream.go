package supervisor

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)
import (
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
)

var _ StreamModel = (*customStreamModel)(nil)

type (
	// StreamModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStreamModel.
	StreamModel interface {
		streamModel
		withSession(session sqlx.Session) StreamModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *Stream) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *Stream) error
		BatchMapByIds(ctx context.Context, Ids []any) (map[int64]Stream, error)
	}

	customStreamModel struct {
		*defaultStreamModel
	}
)

// NewStreamModel returns a model for the database table.
func NewStreamModel(conn sqlx.SqlConn) StreamModel {
	return &customStreamModel{
		defaultStreamModel: newStreamModel(conn),
	}
}

func (m *customStreamModel) withSession(session sqlx.Session) StreamModel {
	return NewStreamModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customStreamModel) InsertSession(ctx context.Context, session sqlx.Session, data *Stream) (int64, error) {
	result, err := m.withSession(session).Insert(ctx, data)
	if err != nil {
		return 0, err
	}

	// 获取新增ID
	aid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if aid == 0 {
		return 0, errors.New("新增事物失败，未生成自增ID")
	}

	return aid, nil
}

// UpdateSession 事务操作-更新
func (m *customStreamModel) UpdateSession(ctx context.Context, session sqlx.Session, data *Stream) error {
	return m.withSession(session).Update(ctx, data)
}

// BatchMapByIds 批量根据IDs查询数据并返回map
func (m *customStreamModel) BatchMapByIds(ctx context.Context, ids []any) (map[int64]Stream, error) {
	var list []Stream
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("id", ids...).Order("id asc").
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]Stream)
	for _, v := range list {
		data[v.Id] = v
	}
	return data, nil
}
