package cron

import (
	"context"
	"muse-admin/internal/svc"
)

// AutoIncrUserLesson 自动增加用户课程
// 此服务可用于定时任务或者其他场景下的自动增加用户课程
// 例如：每天凌晨自动增加用户课程

// todo
//  1、课时调整时，需要注意用户正在训练的进度的完结状态。
//  2、排课时，默认只初始化5节或10节课课程；当小于等2个待学课时时，触发新增排布；
//  3、触发时，需要根据当前最后一个课时所在的位置，填补其后紧挨着的课时；

type AutoIncrLessonCorn struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	//uClassHourModel appointmentModel.UserClassHourModel
	//lessonModel     knowledge.LessonModel
	incrNum map[int64]int64
}

// NewAutoIncrLessonCorn 初始化自动增加用户课时服务
func NewAutoIncrLessonCorn(ctx context.Context, svcCtx *svc.ServiceContext) *AutoIncrLessonCorn {
	return &AutoIncrLessonCorn{
		ctx:    ctx,
		svcCtx: svcCtx,
		//uClassHourModel: appointmentModel.NewUserClassHourModel(svcCtx.MysqlConn),
		//lessonModel:     knowledge.NewLessonModel(svcCtx.MysqlConn),
	}
}

func (c *AutoIncrLessonCorn) Access() bool {
	return true
}

func (c *AutoIncrLessonCorn) SetTitle() string {
	return "定时任务-更新用户排课脚本"
}

func (c *AutoIncrLessonCorn) SetCron() string {
	// 每天凌晨12点执行一次
	return "0 0 0 * * *"
}

func (c *AutoIncrLessonCorn) Run() {
	//_ = c.runCorn()
	// 这里是定时任务执行的具体逻辑
}

//
//// getUserBalance 获取用户课时余额
//func (c *AutoIncrLessonCorn) getUserBalance() (map[int64]int64, error) {
//	result, err := c.uClassHourModel.GetBalance(c.ctx)
//	if err != nil {
//		logc.Errorf(c.ctx, "获取用户课时余额失败，err:%v", err)
//		return nil, err
//	}
//	balance := make(map[int64]int64, len(result))
//	for _, item := range result {
//		balance[item.UserId] = item.Count
//	}
//	return balance, nil
//}
//
//// getHasArrangeNum 获取已经排课的数量
//func (c *AutoIncrLessonCorn) getHasArrangeNum() (map[int64]int64, error) {
//	result, err := c.svcCtx.AssessmentUserLessonModel.GetHasArrangeNum(c.ctx)
//	if err != nil {
//		logc.Errorf(c.ctx, "获取已经排课的数量失败，err:%v", err)
//		return nil, err
//	}
//	lesson := make(map[int64]int64, len(result))
//	for _, item := range result {
//		lesson[item.UserId] = item.Count
//	}
//	return lesson, nil
//}
//
//// getIncrNum 获取增加课时数量
//func (c *AutoIncrLessonCorn) getIncrNum() error {
//	userBalance, err := c.getUserBalance()
//	if err != nil {
//		return err
//	}
//	arrangeNum, err := c.getHasArrangeNum()
//	if err != nil {
//		return err
//	}
//	canUseNum := make(map[int64]int64, len(userBalance))
//	for userId, arrange := range arrangeNum {
//		if arrange >= 5 {
//			// 如果可用数量>=5，则不再增加
//			continue
//		}
//		if balance, ok := userBalance[userId]; ok {
//			if balance == 0 {
//				// 如果余额为0，则不再增加
//				continue
//			}
//			// 如果余额+已排课数量>=5，则可用数量为5-已排课数量
//			if balance+arrange >= 5 {
//				canUseNum[userId] = 5 - arrange
//			} else {
//				// 如果余额+已排课数量<5，则可用数量为余额
//				canUseNum[userId] = balance
//			}
//		}
//	}
//	c.incrNum = canUseNum
//	return nil
//}
//
//// runCorn 执行
//func (c *AutoIncrLessonCorn) runCorn() error {
//	tx, err := sqld.BeginTx(c.ctx, c.svcCtx.MysqlConn)
//	if err != nil {
//		logc.Errorf(c.ctx, "开启事务失败，err:%v", err)
//		return err
//	}
//	// 获取插入的课时
//	lessons, err := c.calInsertLesson()
//	if err != nil {
//		_ = tx.Rollback()
//		return err
//	}
//	err = c.svcCtx.AssessmentUserLessonModel.CreatedMultiWithTx(c.ctx, tx, lessons)
//	if err != nil {
//		_ = tx.Rollback
//		return err
//	}
//	// 提交事务
//	err = tx.Commit()
//	if err != nil {
//		_ = tx.Rollback()
//		logc.Errorf(c.ctx, "提交事务失败，err:%v", err)
//		return err
//	}
//	return nil
//}
//
//// getAllLesson 获取所有课时
//func (c *AutoIncrLessonCorn) getAllLesson() ([]int64, error) {
//	// 获取所有课时
//	allLesson, err := c.lessonModel.FindAll(c.ctx)
//	if err != nil {
//		logc.Errorf(c.ctx, "获取所有课时失败，err:%v", err)
//		return nil, err
//	}
//	result := make([]int64, 0, len(allLesson))
//	for _, item := range allLesson {
//		result = append(result, item.Id)
//	}
//	return result, nil
//}
//
//// calInsertLesson 计算插入课时
//// todo 3、触发时，需要根据当前最后一个课时所在的位置，填补其后紧挨着的课时；
//// todo 具体思路
//// todo 根据用户排序最后的一个课时
//func (c *AutoIncrLessonCorn) calInsertLesson() ([]*assessment.UserLesson, error) {
//	// 获取所有课时
//	allLesson, err := c.getAllLesson()
//	if err != nil {
//		return nil, err
//	}
//	// 获取最后一个用户课时
//	lastLesson, err := c.getLastUserLesson()
//	if err != nil {
//		return nil, err
//	}
//
//	result := make([]*assessment.UserLesson, 0)
//	// 遍历用户待插入课时数量
//	for userId, need := range c.incrNum {
//		// 获取用户最后一个课时
//		if last, ok := lastLesson[userId]; ok {
//			// 获取课时
//			for id, _ := range allLesson {
//				if int64(id) == last.LessonId {
//					// 待用的课时
//					waitList := allLesson[id+1 : id+int(need)+1]
//					for i := 0; i < len(waitList); i++ {
//						c.forkLesson(result, userId, waitList[i], last.LessonIndex+int64(i)+1)
//					}
//				}
//			}
//		} else {
//			waitList := allLesson[:need]
//			for i := 0; i < len(waitList); i++ {
//				c.forkLesson(result, userId, waitList[i], int64(i)+1)
//			}
//		}
//	}
//	return result, nil
//}
//
//// getLastUserLesson 获取最后一个用户课时(未学习)
//func (c *AutoIncrLessonCorn) getLastUserLesson() (map[int64]*incrLastLesson, error) {
//	allLesson, err := c.svcCtx.AssessmentUserLessonModel.GetAllUserLesson(c.ctx, define.UserLessonStatus.UnStudied)
//	if err != nil {
//		return nil, err
//	}
//	result := make(map[int64]*incrLastLesson)
//	for _, item := range allLesson {
//		if _, ok := result[item.UserId]; !ok {
//			result[item.UserId] = &incrLastLesson{
//				LessonId:    item.LessonId,
//				LessonIndex: item.LessonIndex,
//			}
//		}
//	}
//	return result, nil
//}
//
//type incrLastLesson struct {
//	LessonId    int64
//	LessonIndex int64
//}
//
//func (c *AutoIncrLessonCorn) forkLesson(lessons []*assessment.UserLesson, userId, lessonId, index int64) {
//	lessons = append(lessons, &assessment.UserLesson{
//		UserId:      userId,
//		LessonId:    lessonId,
//		Status:      define.UserLessonStatus.UnStudied,
//		LessonIndex: index,
//		CreatedBy:   0,
//		UpdatedBy:   0,
//	})
//}
