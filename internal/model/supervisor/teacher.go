package supervisor

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)
import (
	"github.com/spf13/cast"
)

var _ TeacherModel = (*customTeacherModel)(nil)

type (
	// TeacherModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTeacherModel.
	TeacherModel interface {
		teacherModel
		withSession(session sqlx.Session) TeacherModel
		TableName() string
		FindAllByCondition(ctx context.Context, nameLike string) ([]Teacher, error)
		BatchMapByIds(ctx context.Context, Ids []any) (map[int64]Teacher, error)
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

func (m *customTeacherModel) FindAllByCondition(ctx context.Context, nameLike string) ([]Teacher, error) {
	var resp []Teacher
	query := sqld.NewModel(ctx, m.conn, m.table).Select("*").Order("id desc")
	if nameLike != "" {
		query = query.Where("`teacher_name` like ?", "%"+cast.ToString(nameLike)+"%")
	}
	err := query.Find(&resp)
	return resp, err
}

// BatchMapByIds 批量根据IDs查询数据并返回map
func (m *customTeacherModel) BatchMapByIds(ctx context.Context, ids []any) (map[int64]Teacher, error) {
	var list []Teacher
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("id", ids...).Order("id asc").
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]Teacher)
	for _, v := range list {
		data[v.Id] = v
	}
	return data, nil
}
