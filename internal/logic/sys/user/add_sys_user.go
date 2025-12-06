package user

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/third/feishu"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"

	"encoding/json"
	"errors"
	"muse-admin/internal/define"
	"muse-admin/internal/model/workbench"
	"muse-admin/internal/svc"
	"muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"muse-admin/pkg/other"
	"slices"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddSysUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSysUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSysUserLogic {
	return &AddSysUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSysUserLogic) AddSysUser(req *types.AddSysUserReq) error {
	_, err := l.svcCtx.SysUserModel.FindOneByAccount(l.ctx, req.Account)
	if errors.Is(err, workbench.ErrNotFound) {
		currentUserId := tools.GetUserIdByCtx(l.ctx)
		var currentUserRoleIds []int64
		var roleIds []int64
		if currentUserId == define.SysSuperUserId {
			sysRoleList, _ := l.svcCtx.SysRoleModel.FindAll(l.ctx)
			for _, role := range sysRoleList {
				currentUserRoleIds = append(currentUserRoleIds, role.Id)
				roleIds = append(roleIds, role.Id)
			}

		} else {
			currentUser, _ := l.svcCtx.SysUserModel.FindOne(l.ctx, currentUserId)
			err := json.Unmarshal([]byte(currentUser.RoleIds), &currentUserRoleIds)
			if err != nil {
				return errs.WithCode(err, errs.ServerErrorCode)
			}

			roleIds = append(roleIds, currentUserRoleIds...)
		}

		for _, id := range req.RoleIds {
			if !slices.Contains(roleIds, id) {
				return errs.NewCode(errs.AssigningRolesErrorCode)
			}
		}

		_, err := l.svcCtx.SysDeptModel.FindOne(l.ctx, req.DeptId)
		if err != nil {
			return errs.NewCode(errs.DeptIdErrorCode)
		}

		_, err = l.svcCtx.SysProfessionModel.FindOne(l.ctx, req.ProfessionId)
		if err != nil {
			return errs.NewCode(errs.ProfessionIdErrorCode)
		}

		_, err = l.svcCtx.SysJobModel.FindOne(l.ctx, req.JobId)
		if err != nil {
			return errs.NewCode(errs.JobIdErrorCode)
		}

		var sysUser = new(workbench.User)
		err = copier.Copy(sysUser, req)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		bytes, err := json.Marshal(req.RoleIds)
		sysUser.RoleIds = string(bytes)
		dictionary, err := l.svcCtx.SysDictionaryModel.FindOneByUniqueKey(l.ctx, "sys_pwd")
		var password string
		if dictionary.Status == define.SysEnable {
			password = dictionary.Value
		} else {
			password = define.SysNewUserDefaultPassword
		}

		sysUser.Password = util.GenerateMD5Str(password + l.svcCtx.Config.Salt)
		sysUser.Avatar = other.AvatarUrl()

		// 获取用户OpenId
		openIdMap := feishu.GetOpenId(req.Mobile, "open_id")
		if openId, ok := openIdMap[sysUser.Mobile]; ok {
			sysUser.OpenId = openId
		}

		// 获取用户UserId
		userIdMap := feishu.GetOpenId(req.Mobile, "user_id")
		if userId, ok := userIdMap[sysUser.Mobile]; ok {
			sysUser.UserId = userId
		}

		_, err = l.svcCtx.SysUserModel.Insert(l.ctx, sysUser)
		if err != nil {
			return errs.WithCode(err, errs.ServerErrorCode)
		}

		return nil
	} else {

		return errs.NewCode(errs.AddUserErrorCode)
	}
}
