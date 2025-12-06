package teacher

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ TeacherModel = (*customTeacherModel)(nil)

type (
	// TeacherModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTeacherModel.
	TeacherModel interface {
		teacherModel
		withSession(session sqlx.Session) TeacherModel
		TableName() string
		FindAllRole(ctx context.Context, condition map[string]interface{}) ([]Teacher, error)
	}

	customTeacherModel struct {
		*defaultTeacherModel
	}
)

// NewTeacherModel returns a model for the database table.
func NewTeacherModel(conn sqlx.SqlConn) TeacherModel {
	return &customTeacherModel{
		defaultTeacherModel: newTeacherModel(conn),
	}
}

func (m *customTeacherModel) withSession(session sqlx.Session) TeacherModel {
	return NewTeacherModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customTeacherModel) FindAllRole(ctx context.Context, condition map[string]interface{}) ([]Teacher, error) {
	var resp []Teacher
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").WhereMap(condition).Order("id desc").Find(&resp)
	return resp, err
}
