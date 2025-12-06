package _import

//import (
//	"muse-admin/internal/model/ability/question_group"
//	taskModel "muse-admin/internal/model/excel"
//	"muse-admin/internal/svc"
//)
//
//type QuestionGroup struct {
//	svcCtx                    *svc.ServiceContext
//	taskInfo                  *taskModel.SyncTask
//	backendManageTaskLogModel taskModel.SyncTaskLogModel
//	filename                  string
//	sheetData                 []map[string]string // 所选中的表对象数据
//	groupIds                  []string            // 表对象数据中的组卷ID,数据处理完后需要对其进行总数时使用
//
//	materialModel      question_group.MaterialModel
//	exampleModel       question_group.ExampleModel
//	questionModel      question_group.QuestionModel
//	optionModel        question_group.QuestionOptionModel
//	groupModel         question_group.BizGroupModel
//	groupQuestionModel question_group.BizGroupQuestionModel
//}

//
//func NewGroupBuilder(svcCtx *svc.ServiceContext, taskInfo *taskModel.SyncTask, filename string) IBuilder {
//	return &QuestionGroup{
//		svcCtx:                    svcCtx,
//		taskInfo:                  taskInfo,
//		backendManageTaskLogModel: taskModel.NewSyncTaskLogModel(svcCtx.MysqlConn),
//		filename:                  filename,
//
//		materialModel:      question_group.NewMaterialModel(svcCtx.MysqlLocal),
//		exampleModel:       question_group.NewExampleModel(svcCtx.MysqlLocal),
//		questionModel:      question_group.NewQuestionModel(svcCtx.MysqlLocal),
//		optionModel:        question_group.NewQuestionOptionModel(svcCtx.MysqlLocal),
//		groupModel:         question_group.NewBizGroupModel(svcCtx.MysqlLocal),
//		groupQuestionModel: question_group.NewBizGroupQuestionModel(svcCtx.MysqlLocal),
//	}
//}
//
//func (q *QuestionGroup) GetName() string {
//	return "导入题库且进行组卷"
//}
//
//func (q *QuestionGroup) Access() bool {
//	return q.taskInfo.Type == 30
//}
//
//func (q *QuestionGroup) PreBuild(ctx context.Context) (length int, err error) {
//	// 构建所需要的信息
//	header := map[string]string{
//		"关卡ID":   "game_level_id",
//		"组卷ID":   "group_id",
//		"业务类型":   "biz_type",
//		"题目类型":   "question_type",
//		"使用类型":   "usage_type",
//		"素材标题":   "material_title",
//		"素材作者":   "material_author",
//		"素材来源":   "material_source",
//		"素材内容":   "material_content",
//		"题干":     "ask",
//		"答案":     "answer",
//		"选项一":    "opt1",
//		"选项二":    "opt2",
//		"选项三":    "opt3",
//		"选项四":    "opt4",
//		"解析":     "analysis",
//		"题目引导":   "start_tts",
//		"题干来源":   "source",
//		"例题内容":   "example_content",
//		"例题纯文本":  "example_plain",
//		"小黑板大标题": "example_title",
//		"小黑板小标题": "example_sub_title",
//	}
//	o := excel.NewExCelImport(
//		excel.WithSheetNum(cast.ToInt(q.taskInfo.FileSheet)),
//		excel.WithHeader(header),
//	)
//	q.sheetData, err = o.GetExcelDataOneSheet(ctx, q.filename)
//
//	return len(q.sheetData), err
//}
//
//func (q *QuestionGroup) BuildBefore(ctx context.Context) (*ExecuteResult, error) {
//	var (
//		failedIndex []int // 数据转换失败时需要过滤的下标
//	)
//	resp := &ExecuteResult{
//		TotalCnt:    0,
//		SuccessList: make([]BatchContent, 0),
//		FailList:    make([]BatchContent, 0),
//	}
//
//	for k, v := range q.sheetData {
//		// 数据快照
//		data, _ := json.Marshal(v)
//		_, _ = q.backendManageTaskLogModel.Insert(ctx, &taskModel.SyncTaskLog{
//			TaskId: q.taskInfo.Id,
//			Data:   string(data),
//		})
//
//		// 业务类型转换 1.评测 2.爬天梯 3.原理讲解 4.精析拔高 5.冲刺高考
//		bizType, ok := bizTypeMap[v["biz_type"]]
//		if !ok {
//			failedIndex = append(failedIndex, k)
//			resp.FailList = append(resp.FailList, BatchContent{
//				Index: k,
//				Err:   errors.New(fmt.Sprintf("业务类型:%v,枚举转换失败,请修正导入源后重试", v["biz_type"])),
//			})
//			continue
//		}
//		v["biz_type"] = bizType
//
//		// 题目类型转换 1.单选 2.多选 3.填空 4.判断 5.简答 6.阅读 7.作文
//		questionType, ok := questionTypeMap[v["question_type"]]
//		if !ok {
//			failedIndex = append(failedIndex, k)
//			resp.FailList = append(resp.FailList, BatchContent{
//				Index: k,
//				Err:   errors.New(fmt.Sprintf("题目类型:%v,枚举转换失败,请修正导入源后重试", v["question_type"])),
//			})
//			continue
//		}
//		v["question_type"] = questionType
//
//		// 使用类型转换 1.例题 2.练习题 3.候补题
//		usageType, ok := usageTypeMap[v["usage_type"]]
//		if !ok {
//			failedIndex = append(failedIndex, k)
//			resp.FailList = append(resp.FailList, BatchContent{
//				Index: k,
//				Err:   errors.New(fmt.Sprintf("使用类型:%v,枚举转换失败,请修正导入源后重试", v["usage_type"])),
//			})
//			continue
//		}
//		v["usage_type"] = usageType
//	}
//
//	// 移除不符合的元素
//	q.sheetData = util.ArrayRemoveElements(q.sheetData, failedIndex)
//
//	return resp, nil
//}
//
//func (q *QuestionGroup) BuildData(ctx context.Context, resp *ExecuteResult) {
//	for k, v := range q.sheetData {
//		// 开启事务：todo 如果考虑到事务比较耗时的话，可以考虑批量开启事务，每插入200条提交，减少事务次数，后续考虑优化
//		// 1. 插入素材表：kn_material 获取素材ID
//		// 2. 插入例题表：kn_example 获取例题ID
//		// 3. 插入题库表：kn_question 获取题目ID
//		// 4. 插入题库选项表：kn_question_option
//		// 5. 最后进行组卷
//		// 6. 先获取 kn_biz_group 表的组卷ID，若不存在的话就新增一个组卷ID
//		// 7. 拿到组卷ID后，把题目ID归到组卷题目表里 kn_biz_group_question
//		err := q.svcCtx.MysqlLocal.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
//			var (
//				materialId int64
//				exampleId  int64
//				err        error
//			)
//
//			m := q.convertMaterial(v) // 转换为可直接入库的数据
//			if m.Title != "" || m.Author != "" || m.Source != "" || m.Content.String != "" || m.Background.String != "" || m.AuthorIntro.String != "" {
//				// 查询素材ID是否存在
//				materialId, err = q.materialModel.FindMaterialId(ctx, m)
//				if err != nil {
//					logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,查询素材kn_material表失败,原因:err:%v", q.taskInfo.FileSheetName, k, err)
//					return err
//				}
//				// 若查询的素材ID不存在，新增素材表
//				if materialId == 0 {
//					materialId, err = q.materialModel.InsertSession(ctx, session, m)
//					if err != nil {
//						logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,新增素材kn_material表失败,原因:err:%v", q.taskInfo.FileSheetName, k, err)
//						return err
//					}
//				}
//			}
//
//			e := q.convertExample(v) // 转换为可直接入库的数据
//			if e.Content != "" || e.PlainText != "" || e.Title != "" || e.SubTitle != "" {
//				// 查询例题ID是否存在
//				exampleId, err = q.exampleModel.FindExampleId(ctx, e)
//				if err != nil {
//					logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,查询例题kn_example表失败,原因:err:%v", q.taskInfo.FileSheetName, k, err)
//					return err
//				}
//				// 若查询的例题ID不存在，新增例题表
//				if exampleId == 0 {
//					exampleId, err = q.exampleModel.InsertSession(ctx, session, e)
//					if err != nil {
//						logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,新增例题kn_example表失败,原因:err:%v", q.taskInfo.FileSheetName, k, err)
//						return err
//					}
//				}
//			}
//
//			// 新增题目表
//			questionId, err := q.questionModel.InsertSession(ctx, session, q.convertQuestion(v, materialId, exampleId))
//			if err != nil {
//				logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,新增题目kn_question表失败,原因:err:%v", q.taskInfo.FileSheetName, k, err)
//				return err
//			}
//
//			// 有题目选项
//			opt := q.convertOption(v, questionId)
//			if len(opt) > 0 {
//				for _, op := range opt {
//					_, err = q.optionModel.InsertSession(ctx, session, &op)
//					if err != nil {
//						logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,题目ID:%d 的选项表新增失败,原因:err:%v", q.taskInfo.FileSheetName, k, questionId, err)
//						return err
//					}
//				}
//			}
//
//			// 查询本次需要组卷归类的groupId
//			groupId, err := q.groupModel.FindGroupId(ctx, cast.ToInt64(v["game_level_id"]), cast.ToInt64(v["biz_type"]), v["group_id"])
//			if err != nil {
//				logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,查询组卷kn_biz_group表失败,原因:err:%v", q.taskInfo.FileSheetName, k, err)
//				return err
//			}
//			// 若查询的组卷ID不存在，新增一个组卷记录
//			if groupId == 0 {
//				groupId, err = q.groupModel.InsertSession(ctx, session, q.convertGroup(v))
//				if err != nil {
//					logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,新增组卷kn_biz_group表失败,原因:err:%v", q.taskInfo.FileSheetName, k, err)
//					return err
//				}
//			}
//
//			// 把题目归到组卷里
//			_, err = q.groupQuestionModel.InsertSession(ctx, session, q.convertGroupQuestion(groupId, questionId, k))
//			if err != nil {
//				logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d,新增组卷kn_biz_group表失败,原因:err:%v", q.taskInfo.FileSheetName, k, err)
//				return err
//			}
//
//			// 记录此条题目的组卷ID
//			if !util.IsExist(cast.ToString(groupId), q.groupIds) {
//				q.groupIds = append(q.groupIds, cast.ToString(groupId))
//			}
//
//			// 执行成功
//			resp.SuccessList = append(resp.SuccessList, BatchContent{
//				Index: k,
//			})
//
//			return nil
//		})
//		// 当前事务出错，记录日志，继续执行下一个
//		if err != nil {
//			logc.Errorf(ctx, "表对象名称为:%s,数据下标编号为:%d的事务处理失败,err:%v", q.taskInfo.FileSheetName, k, err)
//
//			// 添加到failList后，继续下一个事务操作
//			resp.FailList = append(resp.FailList, BatchContent{
//				Index: k,
//				Err:   err,
//			})
//			continue
//		}
//	}
//
//	// 本次执行的任务中处理的数据个数总和
//	resp.TotalCnt = cast.ToInt64(len(resp.SuccessList) + len(resp.FailList))
//}
//
//func (q *QuestionGroup) BuildAfter(ctx context.Context) error {
//	// 读取组卷后每个组下面的题库个数总和
//	countGroup, err := q.groupQuestionModel.FindCountByGameLevelIds(ctx, q.groupIds)
//	if err != nil {
//		return err
//	}
//
//	// 将题库个数总和更新到对应的组卷下
//	for _, v := range countGroup {
//		_, err = q.groupModel.UpdateFillFieldsById(ctx, v.BizGroupId, &question_group.BizGroup{
//			QuestionNums: v.Count,
//		})
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//// 转换为kn_material表可录入的数据
//func (q *QuestionGroup) convertMaterial(data map[string]string) *question_group.Material {
//	return &question_group.Material{
//		Title:  strings.TrimSpace(data["material_title"]),
//		Author: strings.TrimSpace(data["material_author"]),
//		Source: strings.TrimSpace(data["material_source"]),
//		Content: sql.NullString{
//			String: strings.TrimSpace(data["material_content"]),
//			Valid:  true,
//		},
//		Background: sql.NullString{
//			String: strings.TrimSpace(data["material_background"]),
//			Valid:  true,
//		},
//		AuthorIntro: sql.NullString{
//			String: strings.TrimSpace(data["material_author_intro"]),
//			Valid:  true,
//		},
//	}
//}
//
//// 转换为kn_example表可录入的数据
//func (q *QuestionGroup) convertExample(data map[string]string) *question_group.Example {
//	return &question_group.Example{
//		Content:   strings.TrimSpace(data["example_content"]),
//		PlainText: strings.TrimSpace(data["example_plain"]),
//		Title:     strings.TrimSpace(data["example_title"]),
//		SubTitle:  strings.TrimSpace(data["example_sub_title"]),
//	}
//}
//
//// 转换为kn_question表可录入的数据
//func (q *QuestionGroup) convertQuestion(data map[string]string, materialId, exampleId int64) *question_group.Question {
//	return &question_group.Question{
//		BizType:      cast.ToInt64(data["biz_type"]),
//		Type:         cast.ToInt64(data["question_type"]),
//		UsageType:    cast.ToInt64(data["usage_type"]),
//		ReviewStatus: 5,
//		GameLevelId:  cast.ToInt64(data["game_level_id"]),
//		MaterialId:   materialId,
//		ExampleId:    exampleId,
//		Ask:          strings.TrimSpace(data["ask"]),
//		Answer: sql.NullString{
//			String: strings.TrimSpace(data["answer"]),
//			Valid:  true,
//		},
//		Analysis: sql.NullString{
//			String: strings.TrimSpace(data["analysis"]),
//			Valid:  true,
//		},
//		StartTts: sql.NullString{
//			String: strings.TrimSpace(data["start_tts"]),
//			Valid:  true,
//		},
//		Source:        strings.TrimSpace(data["source"]),
//		EffectiveTime: time.Now(),
//		ExpiryTime:    util.GetForeverTime(),
//	}
//}
//
//// 转换为kn_question_option表可录入的数据
//func (q *QuestionGroup) convertOption(data map[string]string, questionId int64) []question_group.QuestionOption {
//	optionData := map[string]string{
//		"opt1":   strings.TrimSpace(data["opt1"]),
//		"opt2":   strings.TrimSpace(data["opt2"]),
//		"opt3":   strings.TrimSpace(data["opt3"]),
//		"opt4":   strings.TrimSpace(data["opt4"]),
//		"answer": strings.TrimSpace(data["answer"]),
//	}
//
//	var opt []question_group.QuestionOption
//	answer := strings.TrimSpace(data["answer"])
//
//	labels := []string{"A", "B", "C", "D"}
//
//	for i, label := range labels {
//		optionKey := util.ConcatString("opt", cast.ToString(i+1))
//		if optionData[optionKey] != "" {
//			isAnswer := int64(2)
//			if answer == label {
//				isAnswer = 1
//			}
//			opt = append(opt, question_group.QuestionOption{
//				QuestionId:  questionId,
//				Sequence:    cast.ToInt64(i + 1),
//				OptionLabel: label,
//				Content:     strings.TrimSpace(data[optionKey]),
//				IsAnswer:    isAnswer,
//			})
//		}
//	}
//	return opt
//}
//
//// 转换为kn_biz_group表可录入的数据
//func (q *QuestionGroup) convertGroup(data map[string]string) *question_group.BizGroup {
//	return &question_group.BizGroup{
//		GameLevelId:  cast.ToInt64(data["game_level_id"]),
//		BizType:      cast.ToInt64(data["biz_type"]),
//		ReviewStatus: 5,
//		QuestionNums: 0, // 新增时置为0，count值需要在后置收尾func中读出来回写：BuildAfter()
//		Remark:       data["group_id"],
//	}
//}
//
//// 转换为kn_biz_group_question表可录入的数据
//func (q *QuestionGroup) convertGroupQuestion(groupId, questionId int64, k int) *question_group.BizGroupQuestion {
//	return &question_group.BizGroupQuestion{
//		BizGroupId: groupId,
//		QuestionId: questionId,
//		Sequence:   cast.ToFloat64(k),
//	}
//}
