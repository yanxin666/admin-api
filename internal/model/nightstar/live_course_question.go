package nightstar

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ LiveCourseQuestionModel = (*customLiveCourseQuestionModel)(nil)

type (
	// LiveCourseQuestionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveCourseQuestionModel.
	LiveCourseQuestionModel interface {
		liveCourseQuestionModel
		withSession(session sqlx.Session) LiveCourseQuestionModel
		TableName() string
		FindAllByCourseEventId(ctx context.Context, eventId int64) ([]*LiveCourseQuestion, error)
		BatchInsert(ctx context.Context, list []*LiveCourseQuestion) error
		BatchDeleteByIds(ctx context.Context, ids []int64) error
	}

	customLiveCourseQuestionModel struct {
		*defaultLiveCourseQuestionModel
	}
)

// NewLiveCourseQuestionModel returns a model for the database table.
func NewLiveCourseQuestionModel(conn sqlx.SqlConn) LiveCourseQuestionModel {
	return &customLiveCourseQuestionModel{
		defaultLiveCourseQuestionModel: newLiveCourseQuestionModel(conn),
	}
}

func (m *customLiveCourseQuestionModel) withSession(session sqlx.Session) LiveCourseQuestionModel {
	return NewLiveCourseQuestionModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customLiveCourseQuestionModel) FindAllByCourseEventId(ctx context.Context, eventId int64) ([]*LiveCourseQuestion, error) {
	query := fmt.Sprintf("select %s from %s where deleted_at is null and `live_course_event_id` = ? order by id asc", liveCourseQuestionRows, m.table)
	var resp []*LiveCourseQuestion
	err := m.conn.QueryRowsCtx(ctx, &resp, query, eventId)
	return resp, err
}

func (m *customLiveCourseQuestionModel) BatchInsert(ctx context.Context, list []*LiveCourseQuestion) error {
	inserts := sqld.NewModel(ctx, m.conn, m.table).Inserts(liveCourseQuestionRowsExpectAutoSet)
	for _, v := range list {
		inserts.Append(v)
	}
	_, err := inserts.Execute()
	return err
}

// BatchDeleteByIds 批量删除题目
func (m *customLiveCourseQuestionModel) BatchDeleteByIds(ctx context.Context, ids []int64) error {
	idsStr := strings.Join(util.Int64ArrayToStringArray(ids), ",")
	query := fmt.Sprintf("update %s set deleted_at = NOW() where deleted_at is null and `id` in (%s)", m.table, idsStr)
	_, err := m.conn.ExecCtx(ctx, query)
	return err
}
