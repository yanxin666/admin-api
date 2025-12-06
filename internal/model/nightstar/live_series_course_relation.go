package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LiveSeriesCourseRelationModel = (*customLiveSeriesCourseRelationModel)(nil)

type (
	// LiveSeriesCourseRelationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveSeriesCourseRelationModel.
	LiveSeriesCourseRelationModel interface {
		liveSeriesCourseRelationModel
		withSession(session sqlx.Session) LiveSeriesCourseRelationModel
		TableName() string
	}

	customLiveSeriesCourseRelationModel struct {
		*defaultLiveSeriesCourseRelationModel
	}
)

// NewLiveSeriesCourseRelationModel returns a model for the database table.
func NewLiveSeriesCourseRelationModel(conn sqlx.SqlConn) LiveSeriesCourseRelationModel {
	return &customLiveSeriesCourseRelationModel{
		defaultLiveSeriesCourseRelationModel: newLiveSeriesCourseRelationModel(conn),
	}
}

func (m *customLiveSeriesCourseRelationModel) withSession(session sqlx.Session) LiveSeriesCourseRelationModel {
	return NewLiveSeriesCourseRelationModel(sqlx.NewSqlConnFromSession(session))
}
