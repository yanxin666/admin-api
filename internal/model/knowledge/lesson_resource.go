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
	"time"
)

var _ LessonResourceModel = (*customLessonResourceModel)(nil)

type (
	// LessonResourceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLessonResourceModel.
	LessonResourceModel interface {
		lessonResourceModel
		withSession(session sqlx.Session) LessonResourceModel
		TableName() string
		InsertSession(ctx context.Context, session sqlx.Session, data *LessonResource) (int64, error)
		FindOneByLessonIdAndQuestionId(ctx context.Context, lessonId, questionId int64) (int64, error)
		DeleteSession(ctx context.Context, session sqlx.Session, id int64) error
		BatchLessonIdsByCondition(ctx context.Context, keyword string, condition map[string]interface{}) ([]int64, error)
		BatchByLessonIds(ctx context.Context, lessonIds []any) (map[int64][]LessonResourceLeftJoinQuestion, error)
	}

	customLessonResourceModel struct {
		*defaultLessonResourceModel
	}

	LessonResourceLeftJoinQuestion struct {
		ResourceId   int64   `db:"resource_id"`   // 资源ID           // 主键ID
		LessonId     int64   `db:"lesson_id"`     // 课节ID
		QuestionId   int64   `db:"id"`            // 题目ID
		Sequence     float64 `db:"sequence"`      // 排序，正序，从1开始，默认跨度1
		QuestionNo   string  `db:"question_no"`   // 题目编号
		NodeType     int64   `db:"node_type"`     // 节点类型 1.小灶课 2.小语文 3.大语文
		Type         int64   `db:"type"`          // 试题类型 1：单选，2：多选，3：填空，4：判断，5：简答，6：阅读题 7：作文
		UsageType    int64   `db:"usage_type"`    // 使用类型，1：例题；2：练习题；3：候补题
		GradePhase   int64   `db:"grade_phase"`   // 年级学段，1：小学；2：初中；3：高中
		ReviewStatus int64   `db:"review_status"` // 审核状态 1.初始录入 2.审核中 3.审核通过 4.审核未通过 5.下架
		Level        int64   `db:"level"`         // 等级
		ExampleId    int64   `db:"example_id"`    // 例题ID
		// ExampleContent  string         `db:"example_content"`  // 例题内容
		MaterialId      int64          `db:"material_id"`      // 素材ID
		MaterialContent sql.NullString `db:"material_content"` // 素材内容，比如选取的阅读素材片段
		Ask             string         `db:"ask"`              // 问题
		Answer          sql.NullString `db:"answer"`           // 答案，针对简答题，用于答案判别
		Analysis        sql.NullString `db:"analysis"`         // 题目解析，用于答案判别
		StartTts        sql.NullString `db:"start_tts"`        // 开场白TTS
		Duration        int64          `db:"duration"`         // 预计用时
		Source          string         `db:"source"`           // 来源
		Version         int64          `db:"version"`          // 版本
		CreatedBy       int64          `db:"created_by"`       // 创建者，程序导入为0，其他为操作者ID
		UpdatedBy       int64          `db:"updated_by"`       // 更新者，程序导入为0，其他为操作者ID
		CreatedAt       time.Time      `db:"created_at"`       // 创建时间
		UpdatedAt       time.Time      `db:"updated_at"`       // 更新时间
	}
)

// NewLessonResourceModel returns a model for the database table.
func NewLessonResourceModel(conn sqlx.SqlConn) LessonResourceModel {
	return &customLessonResourceModel{
		defaultLessonResourceModel: newLessonResourceModel(conn),
	}
}

func (m *customLessonResourceModel) withSession(session sqlx.Session) LessonResourceModel {
	return NewLessonResourceModel(sqlx.NewSqlConnFromSession(session))
}

// InsertSession 事务操作-新增
func (m *customLessonResourceModel) InsertSession(ctx context.Context, session sqlx.Session, data *LessonResource) (int64, error) {
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

func (m *customLessonResourceModel) FindOneByLessonIdAndQuestionId(ctx context.Context, lessonId, questionId int64) (int64, error) {
	query := fmt.Sprintf("select id from %s where `lesson_id` = ? and `question_id` = ? limit 1", m.table)
	var id int64
	err := m.conn.QueryRowPartialCtx(ctx, &id, query, lessonId, questionId)
	switch {
	case err == nil:
		return id, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return 0, nil
	default:
		return 0, err
	}
}

func (m *customLessonResourceModel) DeleteSession(ctx context.Context, session sqlx.Session, id int64) error {
	return m.withSession(session).Delete(ctx, id)
}

func (m *customLessonResourceModel) BatchLessonIdsByCondition(ctx context.Context, keyword string, condition map[string]interface{}) ([]int64, error) {
	// SELECT DISTINCT
	// r.lesson_id
	// FROM
	// `kn_lesson_question` q
	// INNER JOIN kn_lesson_resource r ON r.question_id = q.id
	// WHERE
	// r.lesson_id IN ( 18 )
	var lessonIds []int64
	query := sqld.NewModel(ctx, m.conn, m.table).
		Select("DISTINCT r.lesson_id").
		From("kn_lesson_question q").
		Join("inner join kn_lesson_resource r ON r.question_id = q.id")
	if len(condition) != 0 {
		query = query.WhereMap(condition)
	}
	if keyword != "" {
		query = query.Where("q.`ask` like ?", "%"+cast.ToString(keyword)+"%")
	}
	err := query.Find(&lessonIds)
	if err != nil {
		return nil, err
	}

	return lessonIds, nil
}

func (m *customLessonResourceModel) BatchByLessonIds(ctx context.Context, lessonIds []any) (map[int64][]LessonResourceLeftJoinQuestion, error) {
	// SELECT
	//	r.id AS resource_id,
	//	r.lesson_id,
	//	IF(q.usage_type = 1, e.title, q.ask) AS ask,
	//	IF(q.usage_type = 1, e.sub_title, q.answer) AS answer,
	//	IF(q.usage_type = 1, e.lesson_notes, q.analysis) AS analysis,
	// -- 	e.`explain` AS example_content,
	//	m.content AS material_content,
	//	q.*,
	//	r.sequence
	// FROM
	// `kn_lesson_resource` r
	// LEFT JOIN kn_lesson_question q ON r.question_id = q.id
	// LEFT JOIN kn_example e ON q.example_id = e.id
	// LEFT JOIN kn_material m ON q.material_id = m.id
	// WHERE
	// r.lesson_id IN ( 18 )
	var list []LessonResourceLeftJoinQuestion
	err := sqld.NewModel(ctx, m.conn, m.table).
		Select("r.id AS resource_id, r.lesson_id, q.*, "+
			// "IF(q.usage_type = 1, e.title, q.ask) AS ask, "+
			// "IF(q.usage_type = 1, e.sub_title, q.answer) AS answer, "+
			"IF(q.usage_type = 1, e.lesson_notes, q.analysis) AS analysis, "+
			// "e.`explain` AS example_content, "+
			"m.content AS material_content, r.sequence").
		From(m.TableName()+" r").
		Join("left join kn_lesson_question q ON r.question_id = q.id").
		Join("left join kn_example e ON q.example_id = e.id").
		Join("left join kn_material m ON q.material_id = m.id").
		WhereIn("r.lesson_id", lessonIds...).
		Order("r.sequence,r.id").
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64][]LessonResourceLeftJoinQuestion)
	for _, v := range list {
		data[v.LessonId] = append(data[v.LessonId], v)
	}

	return data, nil
}
