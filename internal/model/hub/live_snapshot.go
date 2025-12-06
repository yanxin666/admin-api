package hub

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LiveSnapshotModel = (*customLiveSnapshotModel)(nil)

type (
	// LiveSnapshotModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveSnapshotModel.
	LiveSnapshotModel interface {
		liveSnapshotModel
		withSession(session sqlx.Session) LiveSnapshotModel
		TableName() string
		FindVersionByLessonNo(ctx context.Context, lessonNo string) (int64, error)
	}

	customLiveSnapshotModel struct {
		*defaultLiveSnapshotModel
	}
)

// NewLiveSnapshotModel returns a model for the database table.
func NewLiveSnapshotModel(conn sqlx.SqlConn) LiveSnapshotModel {
	return &customLiveSnapshotModel{
		defaultLiveSnapshotModel: newLiveSnapshotModel(conn),
	}
}

func (m *customLiveSnapshotModel) withSession(session sqlx.Session) LiveSnapshotModel {
	return NewLiveSnapshotModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLiveSnapshotModel) FindVersionByLessonNo(ctx context.Context, lessonNo string) (int64, error) {
	var version int64
	query := fmt.Sprintf("select `version` from %s where `lesson_no` = ? order by id desc limit 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &version, query, lessonNo)
	switch {
	case err == nil:
		return version, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}
