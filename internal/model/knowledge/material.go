package knowledge

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MaterialModel = (*customMaterialModel)(nil)

type (
	// MaterialModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMaterialModel.
	MaterialModel interface {
		materialModel
		withSession(session sqlx.Session) MaterialModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *Material) (int64, error)
		FindMaterialId(ctx context.Context, data *Material) (int64, error)
	}

	customMaterialModel struct {
		*defaultMaterialModel
	}
)

// NewMaterialModel returns a model for the database table.
func NewMaterialModel(conn sqlx.SqlConn) MaterialModel {
	return &customMaterialModel{
		defaultMaterialModel: newMaterialModel(conn),
	}
}

func (m *customMaterialModel) withSession(session sqlx.Session) MaterialModel {
	return NewMaterialModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customMaterialModel) InsertSession(ctx context.Context, session sqlx.Session, data *Material) (int64, error) {
	result, err := m.withSession(session).Insert(ctx, data)
	if err != nil {
		return 0, err
	}

	// 获取新增ID
	aid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if aid == 0 {
		return 0, errors.New("新增事物失败，未生成自增ID")
	}

	return aid, nil
}

// FindMaterialId 根据data内容来获取素材ID
func (m *customMaterialModel) FindMaterialId(ctx context.Context, data *Material) (int64, error) {
	query := fmt.Sprintf("select `id` from %s where `title` = ? and `author` = ? and `source` = ? and `content` = ? and `background` = ? and `author_intro` = ?", m.table)
	var id int64
	err := m.conn.QueryRowPartialCtx(ctx, &id, query, data.Title, data.Author, data.Source, data.Content.String, data.Background.String, data.AuthorIntro.String)
	switch {
	case err == nil:
		return id, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}
