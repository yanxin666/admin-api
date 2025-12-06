package nightstar

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LiveCourseModel = (*customLiveCourseModel)(nil)

type (
	// LiveCourseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveCourseModel.
	LiveCourseModel interface {
		liveCourseModel
		withSession(session sqlx.Session) LiveCourseModel
		TableName() string
		FindOneByNo(ctx context.Context, no string) (*LiveCourse, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *LiveCourse) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *LiveCourse) error
	}

	customLiveCourseModel struct {
		*defaultLiveCourseModel
	}
)

// NewLiveCourseModel returns a model for the database table.
func NewLiveCourseModel(conn sqlx.SqlConn) LiveCourseModel {
	return &customLiveCourseModel{
		defaultLiveCourseModel: newLiveCourseModel(conn),
	}
}

func (m *customLiveCourseModel) withSession(session sqlx.Session) LiveCourseModel {
	return NewLiveCourseModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLiveCourseModel) FindOneByNo(ctx context.Context, no string) (*LiveCourse, error) {
	query := fmt.Sprintf("select %s from %s where `live_no` = ? order by id desc limit 1", liveCourseRows, m.table)
	var resp LiveCourse
	err := m.conn.QueryRowCtx(ctx, &resp, query, no)
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
func (m *customLiveCourseModel) InsertSession(ctx context.Context, session sqlx.Session, data *LiveCourse) (int64, error) {
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
func (m *customLiveCourseModel) UpdateSession(ctx context.Context, session sqlx.Session, data *LiveCourse) error {
	return m.withSession(session).Update(ctx, data)
}
