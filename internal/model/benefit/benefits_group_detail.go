package benefit

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ BenefitsGroupDetailModel = (*customBenefitsGroupDetailModel)(nil)

type (
	// BenefitsGroupDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customBenefitsGroupDetailModel.
	BenefitsGroupDetailModel interface {
		benefitsGroupDetailModel
		withSession(session sqlx.Session) BenefitsGroupDetailModel
		TableName() string
	}

	customBenefitsGroupDetailModel struct {
		*defaultBenefitsGroupDetailModel
	}
)

// NewBenefitsGroupDetailModel returns a model for the database table.
func NewBenefitsGroupDetailModel(conn sqlx.SqlConn) BenefitsGroupDetailModel {
	return &customBenefitsGroupDetailModel{
		defaultBenefitsGroupDetailModel: newBenefitsGroupDetailModel(conn),
	}
}

func (m *customBenefitsGroupDetailModel) withSession(session sqlx.Session) BenefitsGroupDetailModel {
	return NewBenefitsGroupDetailModel(sqlx.NewSqlConnFromSession(session))
}
