package write_ppt

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseOutlineModel = (*customCourseOutlineModel)(nil)

type (
	// CourseOutlineModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseOutlineModel.
	CourseOutlineModel interface {
		courseOutlineModel
		withSession(session sqlx.Session) CourseOutlineModel
		TableName() string
		FindOneByUnitAndSeries(ctx context.Context, unit int64, series string) (*CourseOutline, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *CourseOutline) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *CourseOutline) error
	}

	customCourseOutlineModel struct {
		*defaultCourseOutlineModel
	}
)

// NewCourseOutlineModel returns a model for the database table.
func NewCourseOutlineModel(conn sqlx.SqlConn) CourseOutlineModel {
	return &customCourseOutlineModel{
		defaultCourseOutlineModel: newCourseOutlineModel(conn),
	}
}

func (m *customCourseOutlineModel) withSession(session sqlx.Session) CourseOutlineModel {
	return NewCourseOutlineModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customCourseOutlineModel) FindOneByUnitAndSeries(ctx context.Context, unit int64, series string) (*CourseOutline, error) {
	query := fmt.Sprintf("select %s from %s where `unit` = ? and `series` = ? limit 1", courseOutlineRows, m.table)
	var resp CourseOutline
	err := m.conn.QueryRowCtx(ctx, &resp, query, unit, series)
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
func (m *customCourseOutlineModel) InsertSession(ctx context.Context, session sqlx.Session, data *CourseOutline) (int64, error) {
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
func (m *customCourseOutlineModel) UpdateSession(ctx context.Context, session sqlx.Session, data *CourseOutline) error {
	return m.withSession(session).Update(ctx, data)
}
