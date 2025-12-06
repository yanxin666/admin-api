package knowledge

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/model/tools"
	"time"
)

var _ LessonModel = (*customLessonModel)(nil)

type (
	// LessonModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonModel.
	LessonModel interface {
		lessonModel
		withSession(session sqlx.Session) LessonModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *Lesson) (int64, error)
		FindLessonIdByName(ctx context.Context, data *Lesson) (int64, error)
		FindOneByNo(ctx context.Context, lessonNo string) (*Lesson, error)
		UpdateSession(ctx context.Context, session sqlx.Session, data *Lesson) error
		UpdateFieldsWithTx(ctx context.Context, session sqlx.Session, lessonGroupNo string, data *Lesson) error
		FindOneByParentLesson(ctx context.Context, lessonGroupNo string, level int64) (*Lesson, error)
		FindOneByNoAndGroupNo(ctx context.Context, lessonNo, lessonGroupNo string) (*Lesson, error)

		FindPageByCondition(ctx context.Context, page, limit int64, lessonIds []int64, name, title, subTile string, condition map[string]interface{}) ([]LessonJoinCourse, int64, error)
	}

	customLessonModel struct {
		*defaultLessonModel
	}

	LessonJoinCourse struct {
		Id            int64     `db:"id"`              // 主键ID
		LessonNo      string    `db:"lesson_no"`       // 课节来源编号
		LessonGroupNo string    `db:"lesson_group_no"` // 课程组编号
		NodeType      int64     `db:"node_type"`       // 导入数据时使用，节点类型 1.小灶课 2.小语文 3.大语文
		ParentId      int64     `db:"parent_id"`       // 导入数据时使用，只有大语文才会用到，用作关联父子级
		Level         int64     `db:"level"`           // 难度等级 只有小灶课才会使用，1-4个等级
		LessonType    int64     `db:"lesson_type"`     // 课程类型 1.主线课 2.月考  3.小灶课
		ReviewStatus  int64     `db:"review_status"`   // 审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架
		Remark        string    `db:"remark"`          // 备注
		CreatedAt     time.Time `db:"created_at"`      // 创建时间
		UpdatedAt     time.Time `db:"updated_at"`      // 更新时间

		CourseId    int64        `db:"course_id"`    // 课程ID
		Unit        string       `db:"unit"`         // 单元名称
		Name        string       `db:"name"`         // 章节名称
		Grade       int64        `db:"grade"`        // 所属年级
		ExpectMonth int64        `db:"expect_month"` // 预期月份
		ExpectDate  sql.NullTime `db:"expect_date"`  // 预期上课日期
		Title       string       `db:"title"`        // 标题
		SubTitle    string       `db:"sub_title"`    // 副标题
		LearnTarget string       `db:"learn_target"` // 学习目标
		Sequence    float64      `db:"sequence"`     // 排序，正序，月份相同排序
	}
)

// NewLessonModel returns a model for the database table.
func NewLessonModel(conn sqlx.SqlConn) LessonModel {
	return &customLessonModel{
		defaultLessonModel: newLessonModel(conn),
	}
}

func (m *customLessonModel) withSession(session sqlx.Session) LessonModel {
	return NewLessonModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customLessonModel) InsertSession(ctx context.Context, session sqlx.Session, data *Lesson) (int64, error) {
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

// FindLessonIdByName 根据data内容来获课节ID
func (m *customLessonModel) FindLessonIdByName(ctx context.Context, data *Lesson) (int64, error) {
	query := fmt.Sprintf("select `id` from %s where `name` = ? ", m.table)
	var id int64
	err := m.conn.QueryRowPartialCtx(ctx, &id, query, data.Name)
	switch {
	case err == nil:
		return id, nil
	case errors.Is(err, sqlx.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customLessonModel) FindOneByNo(ctx context.Context, lessonNo string) (*Lesson, error) {
	query := fmt.Sprintf("select %s from %s where `lesson_no` = ? limit 1", lessonRows, m.table)
	var resp Lesson
	err := m.conn.QueryRowCtx(ctx, &resp, query, lessonNo)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customLessonModel) FindOneByNoAndGroupNo(ctx context.Context, lessonNo, lessonGroupNo string) (*Lesson, error) {
	query := fmt.Sprintf("select %s from %s where `lesson_no` = ? and `lesson_group_no` = ?  limit 1", lessonRows, m.table)
	var resp Lesson
	err := m.conn.QueryRowCtx(ctx, &resp, query, lessonNo, lessonGroupNo)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

// UpdateSession 事务操作-更新
func (m *customLessonModel) UpdateSession(ctx context.Context, session sqlx.Session, data *Lesson) error {
	return m.withSession(session).Update(ctx, data)
}

// UpdateFieldsWithTx 事务操作-更新
func (m *customLessonModel) UpdateFieldsWithTx(ctx context.Context, session sqlx.Session, lessonGroupNo string, data *Lesson) error {
	paramSet, args := tools.BuildUpdateSet(data)
	// 构建 UPDATE 语句的 SQL 语句和参数
	query := fmt.Sprintf("UPDATE %s SET %s WHERE lesson_group_no = ?", m.table, paramSet)
	args = append(args, lessonGroupNo)
	_, err := session.ExecCtx(ctx, query, args...)
	return err
}

func (m *customLessonModel) FindOneByParentLesson(ctx context.Context, lessonGroupNo string, level int64) (*Lesson, error) {
	query := fmt.Sprintf("select %s from %s where `lesson_group_no` = ? and `parent_id` = 0 and `level` = ? limit 1", lessonRows, m.table)
	var resp Lesson
	err := m.conn.QueryRowCtx(ctx, &resp, query, lessonGroupNo, level)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (m *customLessonModel) FindPageByCondition(ctx context.Context, page, limit int64, lessonIds []int64, name, title, subTile string, condition map[string]interface{}) ([]LessonJoinCourse, int64, error) {
	// int64变any
	var anyArr []any
	for _, v := range lessonIds {
		anyArr = append(anyArr, v)
	}

	// SELECT
	//	l.*,
	//	c.*
	//	FROM
	// `kn_lesson` l
	// LEFT JOIN ls_course_outline c ON c.lesson_group_no = l.lesson_group_no
	// WHERE
	// l.parent_id = 0

	var resp []LessonJoinCourse
	query := sqld.NewModel(ctx, m.conn, m.table).Select("l.*,c.course_id,c.unit,c.name,c.grade,c.expect_month,c.expect_date,c.title,c.sub_title,c.learn_target,c.sequence").
		From(m.TableName() + " l").
		Join("left join ls_course_outline c ON c.lesson_group_no = l.lesson_group_no").
		WhereMap(condition).
		Order("l.id desc")
	if len(anyArr) != 0 {
		query = query.WhereIn("l.id", anyArr...)
	}
	if name != "" {
		query = query.Where("c.`name` like ?", "%"+cast.ToString(name)+"%")
	}
	if title != "" {
		query = query.Where("c.`title` like ?", "%"+cast.ToString(title)+"%")
	}
	if subTile != "" {
		query = query.Where("c.`sub_title` like ?", "%"+cast.ToString(subTile)+"%")
	}
	total, err := query.Page(int(page), int(limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}

	return resp, total, nil
}
