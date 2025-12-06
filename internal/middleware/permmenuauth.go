package middleware

import (
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"muse-admin/internal/config"
	"muse-admin/internal/define"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"muse-admin/pkg/response"
	"net/http"
	"strings"
)

type PermMenuAuthMiddleware struct {
	c     config.Config
	Redis *redis.Redis
}

func NewPermMenuAuthMiddleware(c config.Config, r *redis.Redis) *PermMenuAuthMiddleware {
	return &PermMenuAuthMiddleware{
		c:     c,
		Redis: r,
	}
}

func (m *PermMenuAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		ctx := r.Context()
		if len(token) < 0 {
			response.ErrorCtx(ctx, w, errs.NewCode(errs.RequestIllegal), nil)
			return
		}

		// 获取用户信息
		userId := cast.ToString(GetUserId(r.Context()))
		val, _ := m.Redis.HgetallCtx(ctx, define.SysOnlineUserInfoPrefix+userId)
		ctx = context.WithValue(ctx, define.SysUserInfoCtx, types.UserInfoCtx{
			Account:  val["account"],
			Username: val["username"],
			Mobile:   val["mobile"],
		})

		online, err := m.Redis.Get(define.SysOnlineUserCachePrefix + cast.ToString(GetUserId(r.Context())))
		if err != nil || online == "" {
			response.ErrorCtx(ctx, w, errs.NewCode(errs.AuthErrorCode), nil)
			return
		}

		// 仅线上环境会判断权限
		if m.c.Mode == "pro" {
			// 判断请求头中的Token是否与Redis存储的等效，防止安全漏洞1
			// admin除外，因为admin可能多人会登录导致token互冲，所以不做校验
			if GetUserId(r.Context()) != 1 && token != online {
				response.ErrorCtx(ctx, w, errs.NewCode(errs.OauthTokenFail), nil)
				return
			}

			uri := strings.Split(r.RequestURI, "?")
			is, err := m.Redis.Sismember(define.SysPermMenuCachePrefix+cast.ToString(GetUserId(r.Context())), uri[0])
			if err != nil || !is {
				response.ErrorCtx(ctx, w, errs.NewCode(errs.NotPermMenuErrorCode), nil)
				return
			}
		}

		// 设置日志自定义内容
		logFields := []logx.LogField{
			logx.Field(define.LogFieldsType.UserID, userId),
		}
		ctx = logx.ContextWithFields(ctx, logFields...)

		// 若要给上下文赋值，这行不可缺失
		r = r.WithContext(ctx)

		// 放行
		next(w, r)
	}
}

// GetUserId 上下文中获取用户ID
func GetUserId(ctx context.Context) int64 {
	var uid int64
	if jsonUid, ok := ctx.Value(define.SysJwtUserId).(json.Number); ok {
		if int64Uid, err := jsonUid.Int64(); err == nil {
			uid = int64Uid
		} else {
			logx.WithContext(ctx).Errorf("GetUidFromCtx err : %+v", err)
		}
	}

	return uid
}
