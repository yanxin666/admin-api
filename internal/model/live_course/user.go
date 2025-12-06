package live_course

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		withSession(session sqlx.Session) UserModel
		TableName() string
		GetStudents(ctx context.Context, streamId, page, limit int64, userName string) ([]User, error)
		GetStudentsCnt(ctx context.Context, streamId int64, userName string) (int64, error)
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUserModel) GetStudents(ctx context.Context, streamId, page, limit int64, userName string) ([]User, error) {
	offset := (page - 1) * limit
	var resp []User
	var err error
	if userName != "" {
		likeKeyword := "%" + userName + "%"
		query := fmt.Sprintf("select %s from %s where stream_id = ? and role_type =2 and  user_name like ? limit ?,? ", userRows, m.table)
		err = m.conn.QueryRowsCtx(ctx, &resp, query, streamId, likeKeyword, offset, limit)
	} else {
		query := fmt.Sprintf("select %s from %s where stream_id = ? and role_type =2 limit ?,? ", userRows, m.table)
		err = m.conn.QueryRowsCtx(ctx, &resp, query, streamId, offset, limit)
	}

	return resp, err
}

func (m *customUserModel) GetStudentsCnt(ctx context.Context, streamId int64, userName string) (int64, error) {
	var resp int64
	var err error
	if userName != "" {
		likeKeyword := "%" + userName + "%"
		query := fmt.Sprintf("select count(*) from %s where stream_id = ? and role_type =2 and user_name like ?", m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, streamId, likeKeyword)
	} else {
		query := fmt.Sprintf("select count(*) from %s where stream_id = ? and role_type =2 ", m.table)
		err = m.conn.QueryRowCtx(ctx, &resp, query, streamId)
	}
	return resp, err
}
