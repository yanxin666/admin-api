package write_ppt

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ GradeCourseSeriesModel = (*customGradeCourseSeriesModel)(nil)

type (
	// GradeCourseSeriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGradeCourseSeriesModel.
	GradeCourseSeriesModel interface {
		gradeCourseSeriesModel
		withSession(session sqlx.Session) GradeCourseSeriesModel
		TableName() string
	}

	customGradeCourseSeriesModel struct {
		*defaultGradeCourseSeriesModel
	}
)

// NewGradeCourseSeriesModel returns a model for the database table.
func NewGradeCourseSeriesModel(conn sqlx.SqlConn) GradeCourseSeriesModel {
	return &customGradeCourseSeriesModel{
		defaultGradeCourseSeriesModel: newGradeCourseSeriesModel(conn),
	}
}

func (m *customGradeCourseSeriesModel) withSession(session sqlx.Session) GradeCourseSeriesModel {
	return NewGradeCourseSeriesModel(sqlx.NewSqlConnFromSession(session))
}
