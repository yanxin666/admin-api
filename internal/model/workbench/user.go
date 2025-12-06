package workbench

import (
	"context"
	"database/sql"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/arch/sqld"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"muse-admin/internal/define"
	"muse-admin/internal/model/tools"
	"muse-admin/internal/types"
	"time"
)

var _ UserModel = (*customUserModel)(nil)

type UserDetail struct {
	Id           int64     `db:"id"`            // 编号
	Account      string    `db:"account"`       // 账号
	Username     string    `db:"username"`      // 姓名
	Nickname     string    `db:"nickname"`      // 昵称
	Avatar       string    `db:"avatar"`        // 头像
	Gender       int64     `db:"gender"`        // 0=保密 1=女 2=男
	Profession   string    `db:"profession"`    // 职称
	ProfessionId int64     `db:"profession_id"` // 职称id
	Job          string    `db:"job"`           // 岗位
	JobId        int64     `db:"job_id"`        // 岗位id
	Dept         string    `db:"dept"`          // 部门
	DeptId       int64     `db:"dept_id"`       // 部门id
	Roles        string    `db:"roles"`         // 角色集
	RoleIds      string    `db:"role_ids"`      // 角色集id
	Email        string    `db:"email"`         // 邮件
	Mobile       string    `db:"mobile"`        // 手机号
	Remark       string    `db:"remark"`        // 备注
	OrderNum     int64     `db:"order_num"`     // 排序值
	Status       int64     `db:"status"`        // 0=禁用 1=开启
	CreateTime   time.Time `db:"create_time"`   // 创建时间
	UpdateTime   time.Time `db:"update_time"`   // 更新时间
}

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		withSession(session sqlx.Session) UserModel
		TableName() string
		UpdateFillFieldsById(ctx context.Context, id int64, data *User) (sql.Result, error)

		FindPage(ctx context.Context, page int64, limit int64, deptIds string) ([]*UserDetail, error)
		FindCountByCondition(ctx context.Context, condition string, value int64) (int64, error)
		FindCountByDeptIds(ctx context.Context, deptIds string) (int64, error)
		FindCountByRoleId(ctx context.Context, roleId int64) (int64, error)
		FindCountByJobId(ctx context.Context, jobId int64) (int64, error)
		FindCountByProfessionId(ctx context.Context, professionId int64) (int64, error)
		FindUserListByRoleIdWithPage(ctx context.Context, username string, page types.PageReq, roleId int64) ([]*UserDetail, int64, error)
		FindAllIncludeDeleted(ctx context.Context) ([]*User, error)
		BatchByUserIds(ctx context.Context, ids []any) (map[int64]User, error)
		GetArrangeClassTeacherList(ctx context.Context, username string, status int64, page types.PageReq, roleId int64) ([]*User, int64, error)
		FindUserByRoleID(ctx context.Context, uid, roleId int64) (*User, error)
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

// UpdateFillFieldsById 根据id 更新用户信息 支持自定义struct
func (m *customUserModel) UpdateFillFieldsById(ctx context.Context, id int64, data *User) (sql.Result, error) {
	paramSet, args := tools.BuildUpdateSet(data)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", m.table, paramSet)
	args = append(args, id)

	result, err := m.conn.ExecCtx(ctx, query, args...)
	return result, err
}

func (m *customUserModel) FindPage(ctx context.Context, page int64, limit int64, deptIds string) ([]*UserDetail, error) {
	// SELECT u.id,u.dept_id,u.job_id,u.profession_id,u.account,u.username,u.nickname,u.avatar,u.gender,IFNULL(p.name,'NULL') as profession,IFNULL(j.name,'NULL') as job,IFNULL(d.name,'NULL') as dept,IFNULL(GROUP_CONCAT(r.name),'NULL') as roles,IFNULL(GROUP_CONCAT(r.id),0) as role_ids,u.email,u.mobile,u.remark,u.order_num,u.status,u.create_time,u.update_time FROM (SELECT * FROM wk_user WHERE id!=1 AND dept_id IN(0,1,3,2) AND delete_time IS NULL ORDER BY order_num DESC LIMIT 0,20) u LEFT JOIN wk_profession p ON u.profession_id=p.id LEFT JOIN wk_dept d ON u.dept_id=d.id LEFT JOIN wk_job j ON u.job_id=j.id LEFT JOIN wk_role r ON JSON_CONTAINS(u.role_ids,JSON_ARRAY(r.id)) GROUP BY u.id order by u.id desc
	offset := (page - 1) * limit
	query := fmt.Sprintf("SELECT u.id,u.dept_id,u.job_id,u.profession_id,u.account,u.username,u.nickname,u.avatar,u.gender,IFNULL(p.name,'NULL') as profession,IFNULL(j.name,'NULL') as job,IFNULL(d.name,'NULL') as dept,IFNULL(GROUP_CONCAT(r.name),'NULL') as roles,IFNULL(GROUP_CONCAT(r.id),0) as role_ids,u.email,u.mobile,u.remark,u.order_num,u.status,u.create_time,u.update_time FROM (SELECT * FROM wk_user WHERE id!=%d AND dept_id IN(%s) AND delete_time IS NULL) u LEFT JOIN wk_profession p ON u.profession_id=p.id LEFT JOIN wk_dept d ON u.dept_id=d.id LEFT JOIN wk_job j ON u.job_id=j.id LEFT JOIN wk_role r ON JSON_CONTAINS(u.role_ids,JSON_ARRAY(r.id)) GROUP BY u.id ORDER BY order_num DESC,u.id DESC LIMIT %d,%d", define.SysSuperUserId, deptIds, offset, limit)
	var resp []*UserDetail
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customUserModel) FindCountByCondition(ctx context.Context, condition string, value int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE %s=%d AND delete_time IS NULL", m.table, condition, value)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *customUserModel) FindCountByDeptIds(ctx context.Context, deptIds string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE id!=%d AND dept_id IN(%s) AND delete_time IS NULL", m.table, define.SysSuperUserId, deptIds)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *customUserModel) FindCountByRoleId(ctx context.Context, roleId int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s u WHERE JSON_CONTAINS(u.role_ids,JSON_ARRAY(%d)) AND u.delete_time IS NULL", m.table, roleId)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *customUserModel) FindCountByJobId(ctx context.Context, jobId int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE job_id=%d AND delete_time IS NULL", m.table, jobId)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *customUserModel) FindCountByProfessionId(ctx context.Context, jobId int64) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE profession_id=%d AND delete_time IS NULL", m.table, jobId)
	var resp int64
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

// FindUserListByRoleIdWithPage 根据角色id查询用户列表(带分页)
func (m *customUserModel) FindUserListByRoleIdWithPage(ctx context.Context, username string, page types.PageReq, roleId int64) ([]*UserDetail, int64, error) {
	query := sqld.NewModel(ctx, m.conn, m.table).
		Select("u.*,IFNULL(p.NAME,'暂无') AS profession,IFNULL(j.NAME,'暂无') AS job,IFNULL(d.NAME,'暂无') AS dept").
		From(fmt.Sprintf("%v %v", m.table, "u")).
		Join("LEFT JOIN wk_profession p ON u.profession_id = p.id").
		Join("LEFT JOIN wk_dept d ON u.dept_id = d.id").
		Join("LEFT JOIN wk_job j ON u.job_id = j.id").
		Where("u.delete_time is null AND JSON_CONTAINS(u.role_ids,JSON_ARRAY(?))", roleId)
	if len(username) > 0 {
		// LIKE 查询
		query = query.Where("u.username LIKE ?", "%"+username+"%")
	}
	var resp []*UserDetail
	total, err := query.Page(int(page.Page), int(page.Limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}
	return resp, total, nil
}

// FindAllIncludeDeleted 查询所有用户列表（包含已删除）
func (m *customUserModel) FindAllIncludeDeleted(ctx context.Context) ([]*User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", m.table)
	var resp []*User
	err := m.conn.QueryRowsCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, sqlx.ErrNotFound
	}
	return resp, nil
}

// BatchByUserIds 批量获取用户信息数据
func (m *customUserModel) BatchByUserIds(ctx context.Context, ids []any) (map[int64]User, error) {
	var list []User
	err := sqld.NewModel(ctx, m.conn, m.table).Select("*").
		WhereIn("id", ids...).
		Find(&list)
	if err != nil {
		return nil, err
	}

	data := make(map[int64]User, len(list))
	for _, v := range list {
		data[v.Id] = v
	}

	return data, nil
}

// GetArrangeClassTeacherList 排班获取教师列表
func (m *customUserModel) GetArrangeClassTeacherList(ctx context.Context, username string, status int64, page types.PageReq, roleId int64) ([]*User, int64, error) {
	query := sqld.NewModel(ctx, m.conn, m.table).
		Select("u.id,u.username,if(c.id is not null,1,2) status").
		From(fmt.Sprintf("%v %v", m.table, "u")).
		Join("LEFT JOIN ch_user_appointment c ON c.teacher_id = u.id AND c.status = ?", define.AppointUserStatus.Available).
		Where("u.delete_time is null AND JSON_CONTAINS(u.role_ids,JSON_ARRAY(?))", roleId)
	if len(username) > 0 {
		// LIKE 查询
		query = query.Where("u.username LIKE ?", "%"+username+"%")
	}
	if status > 0 {
		if status == 1 {
			query = query.Where("c.id is not null")
		} else {
			query = query.Where("c.id is null")
		}
	}
	var resp []*User
	total, err := query.Page(int(page.Page), int(page.Limit), &resp)
	if err != nil {
		return nil, 0, err
	}
	if len(resp) == 0 {
		return nil, 0, nil
	}
	return resp, total, nil
}

// FindUserByRoleID 根据当前UID和角色ID 查询 是否满足
func (m *customUserModel) FindUserByRoleID(ctx context.Context, uid, roleId int64) (*User, error) {
	query := fmt.Sprintf("SELECT * FROM %s  WHERE `id`= %d AND JSON_CONTAINS(role_ids,JSON_ARRAY(%d)) AND delete_time IS NULL", m.table, uid, roleId)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
