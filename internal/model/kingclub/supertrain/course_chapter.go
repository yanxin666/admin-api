package supertrain

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseChapterModel = (*customCourseChapterModel)(nil)

type (
	// CourseChapterModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseChapterModel.
	CourseChapterModel interface {
		courseChapterModel
		withSession(session sqlx.Session) CourseChapterModel
		TableName() string
		FindOneByCourseIdAndChapterNo(ctx context.Context, courseId int64, chapterNo string) (*CourseChapter, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *CourseChapter) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *CourseChapter) error

		BatchByCourseIds(ctx context.Context, courseIds []any) (map[int64][]CourseChapter, error)
		FindAllByCondition(ctx context.Context, condition map[string]interface{}) ([]CourseChapter, error)

		GetCourseIdByChapterList(ctx context.Context, courseId int64) ([]CourseChapter, error)
	}

	customCourseChapterModel struct {
		*defaultCourseChapterModel
	}
)

// NewCourseChapterModel returns a model for the database table.
func NewCourseChapterModel(conn sqlx.SqlConn) CourseChapterModel {
	return &customCourseChapterModel{
		defaultCourseChapterModel: newCourseChapterModel(conn),
	}
}

func (m *customCourseChapterModel) withSession(session sqlx.Session) CourseChapterModel {
	return NewCourseChapterModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customCourseChapterModel) FindOneByCourseIdAndChapterNo(ctx context.Context, courseId int64, chapterNo string) (*CourseChapter, error) {
	query := fmt.Sprintf("select %s from %s where `course_id` = ? and `chapter_no` = ? limit 1", courseChapterRows, m.table)
	var resp CourseChapter
	err := m.conn.QueryRowCtx(ctx, &resp, query, courseId, chapterNo)
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
func (m *customCourseChapterModel) InsertSession(ctx context.Context, session sqlx.Session, data *CourseChapter) (int64, error) {
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
func (m *customCourseChapterModel) UpdateSession(ctx context.Context, session sqlx.Session, data *CourseChapter) error {
	return m.withSession(session).Update(ctx, data)
}

func (m *customCourseChapterModel) BatchByCourseIds(ctx context.Context, courseIds []any) (map[int64][]CourseChapter, error) {
	var list []CourseChapter
	err := sqld.NewModel(ctx, m.conn, m.table).
		Select("*").
		WhereIn("course_id", courseIds...).
		Order("sequence,id").
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64][]CourseChapter)
	for _, v := range list {
		data[v.CourseId] = append(data[v.CourseId], v)
	}

	return data, nil
}

func (m *customCourseChapterModel) FindAllByCondition(ctx context.Context, condition map[string]interface{}) ([]CourseChapter, error) {
	var list []CourseChapter
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").Order("sequence,id")
	for k, v := range condition {
		if k == "name" || k == "chapter_no" {
			query = query.Where(k+" like ?", "%"+cast.ToString(v)+"%")
		} else {
			query = query.Where(k+" = ?", v)
		}
	}

	err := query.Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *customCourseChapterModel) GetCourseIdByChapterList(ctx context.Context, courseId int64) ([]CourseChapter, error) {
	query := fmt.Sprintf("select %s from %s where `course_id` = ? and `status` = ? ", courseChapterRows, m.table)
	var resp []CourseChapter
	err := m.conn.QueryRowsCtx(ctx, &resp, query, courseId, 3)
	return resp, err
}
