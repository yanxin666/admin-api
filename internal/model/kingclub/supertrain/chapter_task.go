package supertrain

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChapterTaskModel = (*customChapterTaskModel)(nil)

type (
	// ChapterTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChapterTaskModel.
	ChapterTaskModel interface {
		chapterTaskModel
		withSession(session sqlx.Session) ChapterTaskModel
		TableName() string
		FindOneByCourseIdAndChapterIdAndTaskNo(ctx context.Context, courseId, chapterId int64, taskNo string) (*ChapterTask, error)
		InsertSession(ctx context.Context, session sqlx.Session, data *ChapterTask) (int64, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *ChapterTask) error
		GetChapterIdByTaskList(ctx context.Context, chapterId int64) ([]ChapterTask, error)
	}

	customChapterTaskModel struct {
		*defaultChapterTaskModel
	}
)

// NewChapterTaskModel returns a model for the database table.
func NewChapterTaskModel(conn sqlx.SqlConn) ChapterTaskModel {
	return &customChapterTaskModel{
		defaultChapterTaskModel: newChapterTaskModel(conn),
	}
}

func (m *customChapterTaskModel) withSession(session sqlx.Session) ChapterTaskModel {
	return NewChapterTaskModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customChapterTaskModel) FindOneByCourseIdAndChapterIdAndTaskNo(ctx context.Context, courseId, chapterId int64, taskNo string) (*ChapterTask, error) {
	query := fmt.Sprintf("select %s from %s where `course_id` = ? and `course_chapter_id` = ? and `task_no` = ? limit 1", chapterTaskRows, m.table)
	var resp ChapterTask
	err := m.conn.QueryRowCtx(ctx, &resp, query, courseId, chapterId, taskNo)
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
func (m *customChapterTaskModel) InsertSession(ctx context.Context, session sqlx.Session, data *ChapterTask) (int64, error) {
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
func (m *customChapterTaskModel) UpdateSession(ctx context.Context, session sqlx.Session, data *ChapterTask) error {
	return m.withSession(session).Update(ctx, data)
}

func (m *customChapterTaskModel) GetChapterIdByTaskList(ctx context.Context, chapterId int64) ([]ChapterTask, error) {
	query := fmt.Sprintf("select %s from %s where `course_chapter_id` = ? AND deleted = ?", chapterTaskRows, m.table)
	var resp []ChapterTask
	err := m.conn.QueryRowsCtx(ctx, &resp, query, chapterId, 1)
	return resp, err
}
