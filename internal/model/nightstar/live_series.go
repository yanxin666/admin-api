package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LiveSeriesModel = (*customLiveSeriesModel)(nil)

type (
	// LiveSeriesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveSeriesModel.
	LiveSeriesModel interface {
		liveSeriesModel
		withSession(session sqlx.Session) LiveSeriesModel
		TableName() string
	}

	customLiveSeriesModel struct {
		*defaultLiveSeriesModel
	}
)

// NewLiveSeriesModel returns a model for the database table.
func NewLiveSeriesModel(conn sqlx.SqlConn) LiveSeriesModel {
	return &customLiveSeriesModel{
		defaultLiveSeriesModel: newLiveSeriesModel(conn),
	}
}

func (m *customLiveSeriesModel) withSession(session sqlx.Session) LiveSeriesModel {
	return NewLiveSeriesModel(sqlx.NewSqlConnFromSession(session))
}
