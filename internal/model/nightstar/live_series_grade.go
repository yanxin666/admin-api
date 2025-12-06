package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LiveSeriesGradeModel = (*customLiveSeriesGradeModel)(nil)

type (
	// LiveSeriesGradeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveSeriesGradeModel.
	LiveSeriesGradeModel interface {
		liveSeriesGradeModel
		withSession(session sqlx.Session) LiveSeriesGradeModel
		TableName() string
	}

	customLiveSeriesGradeModel struct {
		*defaultLiveSeriesGradeModel
	}
)

// NewLiveSeriesGradeModel returns a model for the database table.
func NewLiveSeriesGradeModel(conn sqlx.SqlConn) LiveSeriesGradeModel {
	return &customLiveSeriesGradeModel{
		defaultLiveSeriesGradeModel: newLiveSeriesGradeModel(conn),
	}
}

func (m *customLiveSeriesGradeModel) withSession(session sqlx.Session) LiveSeriesGradeModel {
	return NewLiveSeriesGradeModel(sqlx.NewSqlConnFromSession(session))
}
