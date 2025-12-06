package supervisor

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ScheduleInteractionLinkModel = (*customScheduleInteractionLinkModel)(nil)

type (
	// ScheduleInteractionLinkModel is an interface to be customized, add more methods here,
	// and implement the added methods in customScheduleInteractionLinkModel.
	ScheduleInteractionLinkModel interface {
		scheduleInteractionLinkModel
		withSession(session sqlx.Session) ScheduleInteractionLinkModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *ScheduleInteractionLink) (int64, error)
	}

	customScheduleInteractionLinkModel struct {
		*defaultScheduleInteractionLinkModel
	}
)

// NewScheduleInteractionLinkModel returns a model for the database table.
func NewScheduleInteractionLinkModel(conn sqlx.SqlConn) ScheduleInteractionLinkModel {
	return &customScheduleInteractionLinkModel{
		defaultScheduleInteractionLinkModel: newScheduleInteractionLinkModel(conn),
	}
}

func (m *customScheduleInteractionLinkModel) withSession(session sqlx.Session) ScheduleInteractionLinkModel {
	return NewScheduleInteractionLinkModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customScheduleInteractionLinkModel) InsertSession(ctx context.Context, session sqlx.Session, data *ScheduleInteractionLink) (int64, error) {
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
