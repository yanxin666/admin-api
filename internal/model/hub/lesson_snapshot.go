package hub

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LessonSnapshotModel = (*customLessonSnapshotModel)(nil)

type (
	// LessonSnapshotModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonSnapshotModel.
	LessonSnapshotModel interface {
		lessonSnapshotModel
		withSession(session sqlx.Session) LessonSnapshotModel
		TableName() string
		FindOneByLessonNo(ctx context.Context, lessonNo int64) (*LessonSnapshot, error)
		FindVersionByLessonNo(ctx context.Context, lessonNo int64) (int64, error)

		FindPageByCondition(ctx context.Context, page int64, limit int64, keywordName, keywordPointName string, condition map[string]interface{}) ([]LessonSnapshot, int64, error)
		FindMaxVersion(ctx context.Context) (map[string]int64, error)
		UpdateOperateStatus(ctx context.Context, id, status, userId int64, remark string) error
		FindAllByBothId(ctx context.Context, startId, endId int64) ([]LessonSnapshot, error)
	}

	customLessonSnapshotModel struct {
		*defaultLessonSnapshotModel
	}
)

// NewLessonSnapshotModel returns a model for the database table.
func NewLessonSnapshotModel(conn sqlx.SqlConn) LessonSnapshotModel {
	return &customLessonSnapshotModel{
		defaultLessonSnapshotModel: newLessonSnapshotModel(conn),
	}
}

func (m *customLessonSnapshotModel) withSession(session sqlx.Session) LessonSnapshotModel {
	return NewLessonSnapshotModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLessonSnapshotModel) FindOneByLessonNo(ctx context.Context, lessonNo int64) (*LessonSnapshot, error) {
	query := fmt.Sprintf("select %s from %s where `lesson_no` = ? order by id desc limit 1", lessonSnapshotRows, m.table)
	var resp LessonSnapshot
	err := m.conn.QueryRowCtx(ctx, &resp, query, lessonNo)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customLessonSnapshotModel) FindVersionByLessonNo(ctx context.Context, lessonNo int64) (int64, error) {
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

// FindPageByCondition 分页查询
func (m *customLessonSnapshotModel) FindPageByCondition(ctx context.Context, page int64, limit int64, keywordName, keywordPointName string, condition map[string]interface{}) ([]LessonSnapshot, int64, error) {
	var resp []LessonSnapshot
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if keywordName != "" {
		query = query.Where("`name` like ?", "%"+cast.ToString(keywordName)+"%")
	}
	if keywordPointName != "" {
		query = query.Where("`point_name` like ?", "%"+cast.ToString(keywordPointName)+"%")
	}
	total, err := query.Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}

// FindMaxVersion 相同lessonNo的情况下，获取最大的version数据
func (m *customLessonSnapshotModel) FindMaxVersion(ctx context.Context) (map[string]int64, error) {
	var resp []MaxVersion
	query := fmt.Sprintf("select lesson_no as number, MAX(version) as max_version from bk_lesson_snapshot where version > 1 GROUP BY lesson_no")
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	inspectMap := make(map[string]int64, len(resp)) // map[课程编号][最大版本号]
	for _, v := range resp {
		inspectMap[v.Number] = v.MaxVersion
	}

	return inspectMap, nil
}

func (m *customLessonSnapshotModel) UpdateOperateStatus(ctx context.Context, id, operateStatus, userId int64, remark string) error {
	query := fmt.Sprintf("update %s set `operate_status` = ?,`operate_id` = ?, updated_at = CURRENT_TIMESTAMP, `remark` = ?  where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, operateStatus, userId, remark, id)
	return err
}

func (m *customLessonSnapshotModel) FindAllByBothId(ctx context.Context, startId, endId int64) ([]LessonSnapshot, error) {
	query := fmt.Sprintf("select %s from %s where `id` >= ? and `id`<= ? order by id asc", lessonSnapshotRows, m.table)
	var resp []LessonSnapshot
	err := m.conn.QueryRowsCtx(ctx, &resp, query, startId, endId)
	return resp, err
}
