package nightstar

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LiveAppointmentModel = (*customLiveAppointmentModel)(nil)

type (
	// LiveAppointmentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLiveAppointmentModel.
	LiveAppointmentModel interface {
		liveAppointmentModel
		withSession(session sqlx.Session) LiveAppointmentModel
		TableName() string
	}

	customLiveAppointmentModel struct {
		*defaultLiveAppointmentModel
	}
)

// NewLiveAppointmentModel returns a model for the database table.
func NewLiveAppointmentModel(conn sqlx.SqlConn) LiveAppointmentModel {
	return &customLiveAppointmentModel{
		defaultLiveAppointmentModel: newLiveAppointmentModel(conn),
	}
}

func (m *customLiveAppointmentModel) withSession(session sqlx.Session) LiveAppointmentModel {
	return NewLiveAppointmentModel(sqlx.NewSqlConnFromSession(session))
}
