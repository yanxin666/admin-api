package nightstar

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LiveTeacherModel = (*customLiveTeacherModel)(nil)

type (
	// LiveTeacherModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveTeacherModel.
	LiveTeacherModel interface {
		liveTeacherModel
		withSession(session sqlx.Session) LiveTeacherModel
		TableName() string
		FindOneByTeacherCode(ctx context.Context, teacherCode string) (*LiveTeacher, error)
	}

	customLiveTeacherModel struct {
		*defaultLiveTeacherModel
	}
)

// NewLiveTeacherModel returns a model for the database table.
func NewLiveTeacherModel(conn sqlx.SqlConn) LiveTeacherModel {
	return &customLiveTeacherModel{
		defaultLiveTeacherModel: newLiveTeacherModel(conn),
	}
}

func (m *customLiveTeacherModel) withSession(session sqlx.Session) LiveTeacherModel {
	return NewLiveTeacherModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLiveTeacherModel) FindOneByTeacherCode(ctx context.Context, teacherCode string) (*LiveTeacher, error) {
	query := fmt.Sprintf("select %s from %s where `teacher_code` = ? limit 1", liveTeacherRows, m.table)
	var resp LiveTeacher
	err := m.conn.QueryRowCtx(ctx, &resp, query, teacherCode)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
