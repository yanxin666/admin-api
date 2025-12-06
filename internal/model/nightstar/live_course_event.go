package nightstar

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ LiveCourseEventModel = (*customLiveCourseEventModel)(nil)

type (
	// LiveCourseEventModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveCourseEventModel.
	LiveCourseEventModel interface {
		liveCourseEventModel
		withSession(session sqlx.Session) LiveCourseEventModel
		TableName() string
		FindOneByCourseIdAndStartId(ctx context.Context, courseId int64, startId float64) (*LiveCourseEvent, error)
		FindAllByCourseId(ctx context.Context, courseId int64) ([]*LiveCourseEvent, error)
		BatchDeleteByIds(ctx context.Context, ids []int64) error
	}

	customLiveCourseEventModel struct {
		*defaultLiveCourseEventModel
	}
)

// NewLiveCourseEventModel returns a model for the database table.
func NewLiveCourseEventModel(conn sqlx.SqlConn) LiveCourseEventModel {
	return &customLiveCourseEventModel{
		defaultLiveCourseEventModel: newLiveCourseEventModel(conn),
	}
}

func (m *customLiveCourseEventModel) withSession(session sqlx.Session) LiveCourseEventModel {
	return NewLiveCourseEventModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLiveCourseEventModel) FindOneByCourseIdAndStartId(ctx context.Context, courseId int64, start float64) (*LiveCourseEvent, error) {
	query := fmt.Sprintf("select %s from %s where deleted_at is null and `live_course_id` = ? and `start_id` = ? limit 1", liveCourseEventRows, m.table)
	var resp LiveCourseEvent
	err := m.conn.QueryRowCtx(ctx, &resp, query, courseId, start)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customLiveCourseEventModel) FindAllByCourseId(ctx context.Context, courseId int64) ([]*LiveCourseEvent, error) {
	query := fmt.Sprintf("select %s from %s where deleted_at is null and `live_course_id` = ? order by id asc", liveCourseEventRows, m.table)
	var resp []*LiveCourseEvent
	err := m.conn.QueryRowsCtx(ctx, &resp, query, courseId)
	return resp, err
}

// BatchDeleteByIds 根据ids批量删除
func (m *customLiveCourseEventModel) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	idsStr := strings.Join(util.Int64ArrayToStringArray(ids), ",")
	query := fmt.Sprintf("update %s set deleted_at = NOW() where deleted_at is null and `id` in (%s)", m.table, idsStr)
	_, err := m.conn.ExecCtx(ctx, query)
	return err
}
