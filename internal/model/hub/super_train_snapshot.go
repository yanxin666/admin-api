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

var _ SuperTrainSnapshotModel = (*customSuperTrainSnapshotModel)(nil)

type (
	// SuperTrainSnapshotModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSuperTrainSnapshotModel.
	SuperTrainSnapshotModel interface {
		superTrainSnapshotModel
		withSession(session sqlx.Session) SuperTrainSnapshotModel
		TableName() string
		FindVersionByNo(ctx context.Context, no string) (int64, error)

		FindPageByCondition(ctx context.Context, page int64, limit int64, keywordNo, keywordName string, condition map[string]interface{}) ([]SuperTrainSnapshot, int64, error)
		FindMaxVersion(ctx context.Context) (map[string]int64, error)
		UpdateOperateStatus(ctx context.Context, id, status, userId int64, remark string) error
		FindAllByBothId(ctx context.Context, startId, endId int64) ([]SuperTrainSnapshot, error)
	}

	customSuperTrainSnapshotModel struct {
		*defaultSuperTrainSnapshotModel
	}
)

// NewSuperTrainSnapshotModel returns a model for the database table.
func NewSuperTrainSnapshotModel(conn sqlx.SqlConn) SuperTrainSnapshotModel {
	return &customSuperTrainSnapshotModel{
		defaultSuperTrainSnapshotModel: newSuperTrainSnapshotModel(conn),
	}
}

func (m *customSuperTrainSnapshotModel) withSession(session sqlx.Session) SuperTrainSnapshotModel {
	return NewSuperTrainSnapshotModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customSuperTrainSnapshotModel) FindVersionByNo(ctx context.Context, no string) (int64, error) {
	var version int64
	query := fmt.Sprintf("select `version` from %s where `no` = ? order by id desc limit 1", m.table)
	err := m.conn.QueryRowCtx(ctx, &version, query, no)
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
func (m *customSuperTrainSnapshotModel) FindPageByCondition(ctx context.Context, page int64, limit int64, keywordNo, keywordName string, condition map[string]interface{}) ([]SuperTrainSnapshot, int64, error) {
	var resp []SuperTrainSnapshot
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc")
	if keywordNo != "" {
		query = query.Where("`no` like ?", "%"+cast.ToString(keywordNo)+"%")
	}
	if keywordName != "" {
		query = query.Where("`name` like ?", "%"+cast.ToString(keywordName)+"%")
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
func (m *customSuperTrainSnapshotModel) FindMaxVersion(ctx context.Context) (map[string]int64, error) {
	var resp []MaxVersion

	query := fmt.Sprintf("select no as number, MAX(version) as max_version from bk_super_train_snapshot where version > 1 GROUP BY no")
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

func (m *customSuperTrainSnapshotModel) UpdateOperateStatus(ctx context.Context, id, operateStatus, userId int64, remark string) error {
	query := fmt.Sprintf("update %s set `operate_status` = ?,`operate_id` = ?, `remark` = ?, updated_at = CURRENT_TIMESTAMP  where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, operateStatus, userId, remark, id)
	return err
}

func (m *customSuperTrainSnapshotModel) FindAllByBothId(ctx context.Context, startId, endId int64) ([]SuperTrainSnapshot, error) {
	query := fmt.Sprintf("select %s from %s where `id` >= ? and `id`<= ? order by id asc", lessonSnapshotRows, m.table)
	var resp []SuperTrainSnapshot
	err := m.conn.QueryRowsCtx(ctx, &resp, query, startId, endId)
	return resp, err
}
