package member

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"fmt"
	userModel "muse-admin/internal/model/user"
	"muse-admin/internal/svc"
	ctxt "muse-admin/internal/tools"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type WhiteAdd struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWhiteAdd(ctx context.Context, svcCtx *svc.ServiceContext) *WhiteAdd {
	return &WhiteAdd{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WhiteAdd) WhiteAdd(req *types.WhiteEditReq) (resp *types.Result, err error) {
	// 时间格式检查
	startTime, endTime, err := checkDate(req.Date[0], req.Date[1])
	if err != nil {
		return nil, err
	}

	// 获取全部手机号
	phones, zhCnMap, err := getPhones(req.Phone)
	if err != nil {
		return nil, err
	}

	// 检查手机号是否已存在白名单
	existMap, err := l.svcCtx.UserWhiteModel.BatchPhonesByUnique(l.ctx, phones, req.Product, req.Source)
	if err != nil {
		return nil, err
	}

	var (
		existIds    []any                  // 已存在的白名单ID
		unExistData []*userModel.WhiteList // 不存在的手机号
	)

	remark := "【" + ctxt.GetUserInfoByCtx(l.ctx).Username + "】"
	if req.Remark != "" {
		remark = remark + req.Remark
	}

	for _, v := range phones {
		// 默认备注
		r := remark
		// 筛出存在的手机号
		if data, ok := existMap[v]; ok {
			existIds = append(existIds, data.Id)
		} else {
			// 如果手机号有对应的姓名，则取出
			if zhCnMap != nil && zhCnMap[v] != "" {
				r = remark + "-" + zhCnMap[v]
			}
			unExistData = append(unExistData, &userModel.WhiteList{
				Phone:     v,
				StartTime: startTime,
				EndTime:   endTime,
				Product:   req.Product,
				Source:    req.Source,
				Status:    req.Status,
				Remark:    r,
			})
		}
	}

	// 修改已存在的白名单数据
	if len(existIds) > 0 {
		_, err = l.svcCtx.UserWhiteModel.UpdateColumn(l.ctx, existIds, &userModel.WhiteList{
			StartTime: startTime,
			EndTime:   endTime,
			Status:    req.Status,
			Remark:    req.Remark,
		})
		if err != nil {
			return nil, err
		}
	}

	// 添加未存在的白名单数据
	if len(unExistData) > 0 {
		err = l.svcCtx.UserWhiteModel.BatchInsert(l.ctx, unExistData)
		if err != nil {
			return nil, err
		}
	}

	return
}

// 获取全部手机号，如果有姓名的话，还需要获取手机号针对姓名的映射关系
func getPhones(phone string) ([]string, map[string]string, error) {
	var arr []string

	// 过滤所有的中文逗号、双换行符、所有空格
	phone = strings.ReplaceAll(phone, "，", ",")
	phone = strings.ReplaceAll(phone, "\n\n", " ")
	phone = strings.ReplaceAll(phone, "\u00A0", " ")
	phone = strings.ReplaceAll(phone, " ", "")

	// 检查是否包含逗号或换行符
	hasComma := strings.Contains(phone, ",")
	hasNewline := strings.Contains(phone, "\n")

	switch {
	case hasComma && hasNewline:
		// 同时包含逗号和换行符，先统一替换为逗号再分割
		normalized := strings.ReplaceAll(phone, "\n", ",")
		arr = strings.Split(normalized, ",")
	case hasComma:
		// 只包含逗号
		arr = strings.Split(phone, ",")
	case hasNewline:
		// 只包含换行符
		arr = strings.Split(phone, "\n")
	default:
		arr = []string{phone}
	}

	var (
		pure, errPhones []string
		zhCnMap         = make(map[string]string)
	)
	// 过滤掉空字符串
	for _, s := range arr {
		if s == "" {
			continue
		}

		// 提取纯连续数字
		re := regexp.MustCompile(`(\d+)`)
		p := re.FindString(s)

		// 效验格式是否为11位手机号
		if !util.CheckMobile(p) {
			errPhones = append(errPhones, p)
		}

		// 手机号去重
		if slices.Contains(pure, p) {
			continue
		}
		pure = append(pure, p)

		// 手机号以及名称映射
		zhCn := util.ExtractChinese(s)
		if zhCn != "" {
			zhCnMap[p] = zhCn
		}
	}

	// 手机号有误就返回
	if len(errPhones) > 0 {
		return nil, nil, errs.NewMsg(errs.ErrCodeParamsAbnormal, fmt.Sprintf("手机号格式错误: %s", strings.Join(errPhones, " "))).ShowMsg()
	}

	return pure, zhCnMap, nil
}

func checkDate(startDate, endDate string) (time.Time, time.Time, error) {
	var (
		startTime, endTime time.Time
		err                error
	)

	startTime, err = util.GetStandardDatetime(startDate)
	if err != nil {
		return startTime, endTime, errs.NewMsg(errs.ErrCodeParamsAbnormal, "开始时间日期格式错误").ShowMsg()
	}

	endTime, err = util.GetStandardDatetime(endDate)
	if err != nil {
		return startTime, endTime, errs.NewMsg(errs.ErrCodeParamsAbnormal, "结束时间日期格式错误").ShowMsg()
	}

	return startTime, endTime, nil
}
