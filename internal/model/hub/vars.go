package hub

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var ErrNotFound = sqlx.ErrNotFound

// MaxVersion 相同 number 的情况下，获取最大的version数据
type MaxVersion struct {
	Number     string `db:"number"`      // 来源编号
	MaxVersion int64  `db:"max_version"` // 最大版本号
}
