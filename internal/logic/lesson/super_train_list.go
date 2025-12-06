package lesson

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/model/kingclub/supertrain"
	"muse-admin/pkg/errs"

	"muse-admin/internal/svc"
	"muse-admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SuperTrainList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSuperTrainList(ctx context.Context, svcCtx *svc.ServiceContext) *SuperTrainList {
	return &SuperTrainList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SuperTrainList) SuperTrainList(req *types.SuperTrainListReq) (resp *types.SuperTrainListResp, err error) {
	var courseIds []any
	resp = &types.SuperTrainListResp{
		List:       make([]types.SuperTrainInfo, 0),
		Pagination: types.Pagination{},
	}

	conditionChapter := make(map[string]interface{})
	if req.ChapterId != 0 {
		conditionChapter["id"] = req.ChapterId
	}
	if req.ChapterNo != "" {
		conditionChapter["chapter_no"] = req.ChapterNo
	}
	if req.ChapterName != "" {
		conditionChapter["name"] = req.ChapterName
	}
	if len(conditionChapter) != 0 {
		chapterInfos, err := l.svcCtx.TrainCourseChapterModel.FindAllByCondition(l.ctx, conditionChapter)
		if err != nil {
			return nil, errs.NewCode(errs.ServerErrorCode)
		}
		courseIds, _ = util.ArrayColumn(chapterInfos, "CourseId")
		courseIds = util.ArrayUniqueValue(courseIds)
		if len(courseIds) == 0 {
			return resp, nil
		}
	}

	// 构建查询条件
	condition := make(map[string]interface{})
	if req.CourseId != 0 {
		condition["id"] = req.CourseId
	}
	if req.Type != 0 {
		condition["course_type"] = req.Type
	}
	if req.Status != 0 {
		condition["status"] = req.Status
	}
	conditionLike := map[string]string{
		"course_no": req.CourseNo,
		"name":      req.Name,
	}
	list, total, err := l.svcCtx.TrainCourseModel.FindPageByCondition(l.ctx, req.Page, req.Limit, condition, conditionLike, courseIds)
	if err != nil {
		logc.Errorf(l.ctx, "批量查询课程失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	var resourceInfo map[int64][]supertrain.CourseChapter
	if len(list) != 0 {
		courseIds, _ := util.ArrayColumn(list, "Id")
		resourceInfo, err = l.svcCtx.TrainCourseChapterModel.BatchByCourseIds(l.ctx, courseIds)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
	}

	for _, v := range list {
		var subChapterArr []types.SuperTrainSubInfo // 课程下面的章节列表

		for _, u := range resourceInfo[v.Id] {
			subChapter := types.SuperTrainSubInfo{
				ChapterId:  u.Id,
				ChapterNo:  u.ChapterNo,
				Type:       u.Type,
				Name:       u.Name,
				Image:      u.Image,
				Intro:      u.Intro,
				Status:     u.Status,
				IsCanLearn: u.CanLearn,
				IsNew:      u.IsNew,
				Extra:      u.Extra.String,
				CreatedAt:  u.CreatedAt.Unix(),
				UpdatedAt:  u.UpdatedAt.Unix(),
			}

			subChapterArr = append(subChapterArr, subChapter)
		}

		// 课程列表
		resp.List = append(resp.List, types.SuperTrainInfo{
			CourseId:   v.Id,
			CourseNo:   v.CourseNo,
			CourseType: v.CourseType,
			Name:       v.Name,
			Image:      v.Image,
			Intro:      v.Intro,
			Status:     v.Status,
			OpenTime:   v.OpenTime.Unix(),
			Extra:      v.Extra.String,
			SubList:    subChapterArr,
			CreatedAt:  v.CreatedAt.Unix(),
			UpdatedAt:  v.UpdatedAt.Unix(),
		})
	}

	// 分页信息
	resp.Pagination = types.Pagination{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	return
}
