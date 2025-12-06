package live

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LiveRoleModel = (*customLiveRoleModel)(nil)

type (
	// LiveRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveRoleModel.
	LiveRoleModel interface {
		liveRoleModel
		withSession(session sqlx.Session) LiveRoleModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *LiveRole) (int64, error)
	}

	customLiveRoleModel struct {
		*defaultLiveRoleModel
	}
)

// NewLiveRoleModel returns a model for the database table.
func NewLiveRoleModel(conn sqlx.SqlConn) LiveRoleModel {
	return &customLiveRoleModel{
		defaultLiveRoleModel: newLiveRoleModel(conn),
	}
}

func (m *customLiveRoleModel) withSession(session sqlx.Session) LiveRoleModel {
	return NewLiveRoleModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customLiveRoleModel) InsertSession(ctx context.Context, session sqlx.Session, data *LiveRole) (int64, error) {
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
