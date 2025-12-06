package supertrain

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseModel = (*customCourseModel)(nil)

type (
	// CourseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseModel.
	CourseModel interface {
		courseModel
		withSession(session sqlx.Session) CourseModel
		TableName() string
		FindOneByCourseNo(ctx context.Context, courseNo string) (*Course, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *Course) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *Course) error
		FindPageByCondition(ctx context.Context, page, limit int64, condition map[string]interface{}, conditionLike map[string]string, courseIds []any) ([]Course, int64, error)
		BatchMapByIds(ctx context.Context, Ids []any) (map[int64]Course, error)
		FindListByConds(ctx context.Context, status int64, condition map[string]interface{}) ([]Course, error)
	}

	customCourseModel struct {
		*defaultCourseModel
	}
)

// NewCourseModel returns a model for the database table.
func NewCourseModel(conn sqlx.SqlConn) CourseModel {
	return &customCourseModel{
		defaultCourseModel: newCourseModel(conn),
	}
}

func (m *customCourseModel) withSession(session sqlx.Session) CourseModel {
	return NewCourseModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customCourseModel) FindOneByCourseNo(ctx context.Context, courseNo string) (*Course, error) {
	query := fmt.Sprintf("select %s from %s where `course_no` = ? limit 1", courseRows, m.table)
	var resp Course
	err := m.conn.QueryRowCtx(ctx, &resp, query, courseNo)
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
func (m *customCourseModel) InsertSession(ctx context.Context, session sqlx.Session, data *Course) (int64, error) {
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
func (m *customCourseModel) UpdateSession(ctx context.Context, session sqlx.Session, data *Course) error {
	return m.withSession(session).Update(ctx, data)
}

func (m *customCourseModel) FindPageByCondition(ctx context.Context, page, limit int64, condition map[string]interface{}, conditionLike map[string]string, courseIds []any) ([]Course, int64, error) {
	var resp []Course
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereMap(condition).
		Order("id desc")
	for k, v := range conditionLike {
		if v != "" {
			query = query.Where(k+" like ?", "%"+v+"%")
		}
	}
	if len(courseIds) > 0 {
		query = query.WhereIn("id", courseIds...)
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

// BatchMapByIds 批量根据IDs查询数据并返回map
func (m *customCourseModel) BatchMapByIds(ctx context.Context, ids []any) (map[int64]Course, error) {
	var list []Course
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("id", ids...).Order("id asc").
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]Course)
	for _, v := range list {
		data[v.Id] = v
	}

	return data, nil
}

// FindListByConds 根据条件查询列表
func (m *customCourseModel) FindListByConds(ctx context.Context, status int64, condition map[string]interface{}) ([]Course, error) {
	db := sqld.NewModel(ctx, m.conn, m.table)
	var resp []Course
	err := db.Select("*").WhereMap(condition).Where("status = ?", status).Find(&resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
