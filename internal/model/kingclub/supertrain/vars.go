package supertrain

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var ErrNotFound = sqlx.ErrNotFound

type CourseJoinChapter struct {
	Id        int64          `db:"id"`          // 主键
	CourseNo  string         `db:"course_no"`   // 课程编号
	Type      int64          `db:"course_type"` // 课程类型 1.超练作文慢练 2.超练作文快练 3.超练阅读
	Name      string         `db:"name"`        // 课程名称
	Image     string         `db:"image"`       // 课程封面图
	Intro     string         `db:"intro"`       // 课程介绍
	Status    int64          `db:"status"`      // 状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架
	OpenTime  time.Time      `db:"open_time"`   // 课程开放时间
	Remark    string         `db:"remark"`      // 备注
	Extra     sql.NullString `db:"extra"`       // 拓展字段，用于存json
	CreatedAt time.Time      `db:"created_at"`  // 创建时间
	UpdatedAt time.Time      `db:"updated_at"`  // 更新时间

	ChapterId        int64          `db:"chapter_id"`         // 主键
	ChapterNo        string         `db:"chapter_chapter_no"` // 章节编号
	ChapterType      int64          `db:"chapter_type"`       // 1 任务 2 文章
	ChapterName      string         `db:"chapter_name"`       // 标题，后续不再维护
	ChapterTitle     string         `db:"chapter_title"`      // 精简版标题
	ChapterImage     string         `db:"chapter_image"`      // 课程封面图
	ChapterIntro     string         `db:"chapter_intro"`      // 课程介绍
	ChapterStatus    int64          `db:"chapter_status"`     // 状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架
	ChapterCanLearn  int64          `db:"chapter_can_learn"`  // 是否能学 1.可以 2.不可以
	ChapterIsNew     int64          `db:"chapter_is_new"`     // 是否上新 1.否 2.是
	ChapterExtra     sql.NullString `db:"chapter_extra"`      // 拓展字段，用于存json
	ChapterCreatedAt time.Time      `db:"chapter_created_at"` // 创建时间
	ChapterUpdatedAt time.Time      `db:"chapter_updated_at"` // 更新时间
}
