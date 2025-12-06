package nightstar

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ LiveCourseQuestionOptionModel = (*customLiveCourseQuestionOptionModel)(nil)

type (
	// LiveCourseQuestionOptionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveCourseQuestionOptionModel.
	LiveCourseQuestionOptionModel interface {
		liveCourseQuestionOptionModel
		withSession(session sqlx.Session) LiveCourseQuestionOptionModel
		TableName() string
		BatchDeleteByQuestionIds(ctx context.Context, questionIds []int64) error
	}

	customLiveCourseQuestionOptionModel struct {
		*defaultLiveCourseQuestionOptionModel
	}
)

// NewLiveCourseQuestionOptionModel returns a model for the database table.
func NewLiveCourseQuestionOptionModel(conn sqlx.SqlConn) LiveCourseQuestionOptionModel {
	return &customLiveCourseQuestionOptionModel{
		defaultLiveCourseQuestionOptionModel: newLiveCourseQuestionOptionModel(conn),
	}
}

func (m *customLiveCourseQuestionOptionModel) withSession(session sqlx.Session) LiveCourseQuestionOptionModel {
	return NewLiveCourseQuestionOptionModel(sqlx.NewSqlConnFromSession(session))
}

// BatchDeleteByQuestionIds 根据questionIds批量删除
func (m *customLiveCourseQuestionOptionModel) BatchDeleteByQuestionIds(ctx context.Context, questionIds []int64) error {
	idsStr := strings.Join(util.Int64ArrayToStringArray(questionIds), ",")
	query := fmt.Sprintf("update %s set deleted_at = NOW() where deleted_at is null and `live_course_question_id` in (%s)", m.table, idsStr)
	_, err := m.conn.ExecCtx(ctx, query)
	return err
}
