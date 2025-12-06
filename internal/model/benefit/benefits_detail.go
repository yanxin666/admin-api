package benefit

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BenefitsDetailModel = (*customBenefitsDetailModel)(nil)

type (
	// BenefitsDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBenefitsDetailModel.
	BenefitsDetailModel interface {
		benefitsDetailModel
		withSession(session sqlx.Session) BenefitsDetailModel
		TableName() string
	}

	customBenefitsDetailModel struct {
		*defaultBenefitsDetailModel
	}
)

// NewBenefitsDetailModel returns a model for the database table.
func NewBenefitsDetailModel(conn sqlx.SqlConn) BenefitsDetailModel {
	return &customBenefitsDetailModel{
		defaultBenefitsDetailModel: newBenefitsDetailModel(conn),
	}
}

func (m *customBenefitsDetailModel) withSession(session sqlx.Session) BenefitsDetailModel {
	return NewBenefitsDetailModel(sqlx.NewSqlConnFromSession(session))
}
