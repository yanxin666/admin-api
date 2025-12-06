package benefit

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BenefitsNotifyLogModel = (*customBenefitsNotifyLogModel)(nil)

type (
	// BenefitsNotifyLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBenefitsNotifyLogModel.
	BenefitsNotifyLogModel interface {
		benefitsNotifyLogModel
		withSession(session sqlx.Session) BenefitsNotifyLogModel
		TableName() string
	}

	customBenefitsNotifyLogModel struct {
		*defaultBenefitsNotifyLogModel
	}
)

// NewBenefitsNotifyLogModel returns a model for the database table.
func NewBenefitsNotifyLogModel(conn sqlx.SqlConn) BenefitsNotifyLogModel {
	return &customBenefitsNotifyLogModel{
		defaultBenefitsNotifyLogModel: newBenefitsNotifyLogModel(conn),
	}
}

func (m *customBenefitsNotifyLogModel) withSession(session sqlx.Session) BenefitsNotifyLogModel {
	return NewBenefitsNotifyLogModel(sqlx.NewSqlConnFromSession(session))
}
