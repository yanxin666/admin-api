package supertrain

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MedalModel = (*customMedalModel)(nil)

type (
	// MedalModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMedalModel.
	MedalModel interface {
		medalModel
		withSession(session sqlx.Session) MedalModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *Medal) (int64, error)
	}

	customMedalModel struct {
		*defaultMedalModel
	}
)

// NewMedalModel returns a model for the database table.
func NewMedalModel(conn sqlx.SqlConn) MedalModel {
	return &customMedalModel{
		defaultMedalModel: newMedalModel(conn),
	}
}

func (m *customMedalModel) withSession(session sqlx.Session) MedalModel {
	return NewMedalModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customMedalModel) InsertSession(ctx context.Context, session sqlx.Session, data *Medal) (int64, error) {
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
