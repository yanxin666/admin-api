package lesson

import (
	"context"
	"e.coding.net/zmexing/nenglitanzhen/biz-lib/util"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"muse-admin/internal/model/knowledge"
	"muse-admin/internal/svc"
	"muse-admin/internal/types"
	"muse-admin/pkg/errs"

	"github.com/zeromicro/go-zero/core/logx"
)

type LessonList struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLessonList(ctx context.Context, svcCtx *svc.ServiceContext) *LessonList {
	return &LessonList{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LessonList) LessonList(req *types.LessonListReq) (resp *types.LessonListResp, err error) {
	var (
		resourceInfo       map[int64][]knowledge.LessonResourceLeftJoinQuestion
		keywordQuestionAsk string
		lessonIdArr        []int64
	)

	resp = &types.LessonListResp{
		List:       make([]types.LessonInfo, 0),
		Pagination: types.Pagination{},
	}

	// 题目搜索条件
	questionCondition := make(map[string]interface{})
	if req.QuestionId != 0 {
		questionCondition["q.id"] = req.QuestionId
	}
	if req.QuestionNo != "" {
		questionCondition["q.question_no"] = req.QuestionNo
	}
	if req.QuestionUsageType != 0 {
		questionCondition["q.usage_type"] = req.QuestionUsageType
	}
	if req.QuestionNodeType != 0 {
		questionCondition["q.node_type"] = req.QuestionNodeType
	}
	if req.QuestionType != 0 {
		questionCondition["q.type"] = req.QuestionType
	}
	if req.QuestionStatus != 0 {
		questionCondition["q.review_status"] = req.QuestionStatus
	}
	if req.QuestionAsk != "" {
		keywordQuestionAsk = req.QuestionAsk
	}

	// 题目查询无数据就返回
	if keywordQuestionAsk != "" || len(questionCondition) != 0 {
		lessonIdArr, err = l.svcCtx.LessonResourceModel.BatchLessonIdsByCondition(l.ctx, keywordQuestionAsk, questionCondition)
		if err != nil {
			logc.Errorf(l.ctx, "批量查询课程失败，Err:%s", err)
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
		if lessonIdArr == nil {
			return resp, nil
		}
	}

	// 构建查询条件
	condition := make(map[string]interface{})
	condition["l.`parent_id`"] = 0
	if req.LessonId != 0 {
		lessonIdArr = append(lessonIdArr, req.LessonId)
		// condition["l.id"] = req.LessonId
	}
	if req.LessonNo != 0 {
		condition["l.lesson_no"] = req.LessonNo
	}
	if req.Status != 0 {
		condition["l.review_status"] = req.Status
	}
	if req.Grade != 0 {
		condition["c.grade"] = req.Grade
	}

	list, total, err := l.svcCtx.LessonModel.FindPageByCondition(l.ctx, req.Page, req.Limit, lessonIdArr, req.Name, req.Title, req.SubTitle, condition)
	if err != nil {
		logc.Errorf(l.ctx, "批量查询课程失败，Err:%s", err)
		return nil, errs.WithCode(err, errs.ServerErrorCode)
	}

	if len(list) != 0 {
		lessonIds, _ := util.ArrayColumn(list, "Id")
		resourceInfo, err = l.svcCtx.LessonResourceModel.BatchByLessonIds(l.ctx, lessonIds)
		if err != nil {
			return nil, errs.WithCode(err, errs.ServerErrorCode)
		}
	}

	for _, v := range list {
		var lessonArr []types.LessonSubInfo // 课程下面的题库列表

		for _, u := range resourceInfo[v.Id] {
			subLesson := types.LessonSubInfo{
				NodeType:        u.NodeType,
				QuestionId:      u.QuestionId,
				QuestionNo:      u.QuestionNo,
				Type:            u.Type,
				UsageType:       u.UsageType,
				GradePhase:      u.GradePhase,
				ReviewStatus:    u.ReviewStatus,
				Level:           u.Level,
				ExampleId:       u.ExampleId,
				MaterialId:      u.MaterialId,
				MaterialContent: u.MaterialContent.String,
				Ask:             u.Ask,
				Answer:          u.Answer.String,
				Analysis:        u.Analysis.String,
				CreatedAt:       u.CreatedAt.Unix(),
				UpdatedAt:       u.UpdatedAt.Unix(),
			}

			lessonArr = append(lessonArr, subLesson)
		}

		// 多个课程目标
		var learnTargetList []string
		_ = json.Unmarshal([]byte(v.LearnTarget), &learnTargetList)

		// 课程列表
		resp.List = append(resp.List, types.LessonInfo{
			LessonId:      v.Id,
			LessonNo:      v.LessonNo,
			LessonGroupNo: v.LessonGroupNo,
			ParentId:      v.ParentId,
			Level:         v.Level,
			Name:          v.Name,
			LessonType:    v.LessonType,
			ReviewStatus:  v.ReviewStatus,
			Title:         v.Title,
			SubTitle:      v.SubTitle,
			LearnTarget:   learnTargetList,
			Remark:        v.Remark,
			Grade:         v.Grade,
			SubList:       lessonArr,
			CreatedAt:     v.CreatedAt.Unix(),
			UpdatedAt:     v.UpdatedAt.Unix(),
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
