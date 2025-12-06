package hub

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WritePptSnapshotModel = (*customWritePptSnapshotModel)(nil)

type (
	// WritePptSnapshotModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWritePptSnapshotModel.
	WritePptSnapshotModel interface {
		writePptSnapshotModel
		withSession(session sqlx.Session) WritePptSnapshotModel
		TableName() string
		FindVersionByLessonNo(ctx context.Context, lessonNo int64) (int64, error)

		UpdateOperateStatus(ctx context.Context, id, status, userId int64, remark string) error
		FindAllByBothId(ctx context.Context, startId, endId int64) ([]WritePptSnapshot, error)
	}

	customWritePptSnapshotModel struct {
		*defaultWritePptSnapshotModel
	}
)

// NewWritePptSnapshotModel returns a model for the database table.
func NewWritePptSnapshotModel(conn sqlx.SqlConn) WritePptSnapshotModel {
	return &customWritePptSnapshotModel{
		defaultWritePptSnapshotModel: newWritePptSnapshotModel(conn),
	}
}

func (m *customWritePptSnapshotModel) withSession(session sqlx.Session) WritePptSnapshotModel {
	return NewWritePptSnapshotModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customWritePptSnapshotModel) FindVersionByLessonNo(ctx context.Context, lessonNo int64) (int64, error) {
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

func (m *customWritePptSnapshotModel) UpdateOperateStatus(ctx context.Context, id, operateStatus, userId int64, remark string) error {
	query := fmt.Sprintf("update %s set `operate_status` = ?,`operate_id` = ?, `remark` = ?, updated_at = CURRENT_TIMESTAMP  where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, operateStatus, userId, remark, id)
	return err
}

func (m *customWritePptSnapshotModel) FindAllByBothId(ctx context.Context, startId, endId int64) ([]WritePptSnapshot, error) {
	query := fmt.Sprintf("select %s from %s where `id` >= ? and `id`<= ? order by id asc", writePptSnapshotRows, m.table)
	var resp []WritePptSnapshot
	err := m.conn.QueryRowsCtx(ctx, &resp, query, startId, endId)
	return resp, err
}
