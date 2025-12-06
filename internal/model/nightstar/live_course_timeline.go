package nightstar

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LiveCourseTimelineModel = (*customLiveCourseTimelineModel)(nil)

type (
	// LiveCourseTimelineModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveCourseTimelineModel.
	LiveCourseTimelineModel interface {
		liveCourseTimelineModel
		withSession(session sqlx.Session) LiveCourseTimelineModel
		TableName() string
		FindOneByCourseId(ctx context.Context, courseId int64) (*LiveCourseTimeline, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *LiveCourseTimeline) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *LiveCourseTimeline) error
	}

	customLiveCourseTimelineModel struct {
		*defaultLiveCourseTimelineModel
	}
)

// NewLiveCourseTimelineModel returns a model for the database table.
func NewLiveCourseTimelineModel(conn sqlx.SqlConn) LiveCourseTimelineModel {
	return &customLiveCourseTimelineModel{
		defaultLiveCourseTimelineModel: newLiveCourseTimelineModel(conn),
	}
}

func (m *customLiveCourseTimelineModel) withSession(session sqlx.Session) LiveCourseTimelineModel {
	return NewLiveCourseTimelineModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLiveCourseTimelineModel) FindOneByCourseId(ctx context.Context, courseId int64) (*LiveCourseTimeline, error) {
	query := fmt.Sprintf("select %s from %s where `live_course_id` = ? limit 1", liveCourseTimelineRows, m.table)
	var resp LiveCourseTimeline
	err := m.conn.QueryRowCtx(ctx, &resp, query, courseId)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

// InsertSession 事务操作-新增
func (m *customLiveCourseTimelineModel) InsertSession(ctx context.Context, session sqlx.Session, data *LiveCourseTimeline) (int64, error) {
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
func (m *customLiveCourseTimelineModel) UpdateSession(ctx context.Context, session sqlx.Session, data *LiveCourseTimeline) error {
	return m.withSession(session).Update(ctx, data)
}
