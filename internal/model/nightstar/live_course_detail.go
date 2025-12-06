package nightstar

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LiveCourseDetailModel = (*customLiveCourseDetailModel)(nil)

type (
	// LiveCourseDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveCourseDetailModel.
	LiveCourseDetailModel interface {
		liveCourseDetailModel
		withSession(session sqlx.Session) LiveCourseDetailModel
		TableName() string
		FindOneByCourseId(ctx context.Context, courseId int64) (*LiveCourseDetail, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *LiveCourseDetail) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *LiveCourseDetail) error
	}

	customLiveCourseDetailModel struct {
		*defaultLiveCourseDetailModel
	}
)

// NewLiveCourseDetailModel returns a model for the database table.
func NewLiveCourseDetailModel(conn sqlx.SqlConn) LiveCourseDetailModel {
	return &customLiveCourseDetailModel{
		defaultLiveCourseDetailModel: newLiveCourseDetailModel(conn),
	}
}

func (m *customLiveCourseDetailModel) withSession(session sqlx.Session) LiveCourseDetailModel {
	return NewLiveCourseDetailModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLiveCourseDetailModel) FindOneByCourseId(ctx context.Context, courseId int64) (*LiveCourseDetail, error) {
	query := fmt.Sprintf("select %s from %s where `live_course_id` = ? limit 1", liveCourseDetailRows, m.table)
	var resp LiveCourseDetail
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
func (m *customLiveCourseDetailModel) InsertSession(ctx context.Context, session sqlx.Session, data *LiveCourseDetail) (int64, error) {
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
func (m *customLiveCourseDetailModel) UpdateSession(ctx context.Context, session sqlx.Session, data *LiveCourseDetail) error {
	return m.withSession(session).Update(ctx, data)
}
