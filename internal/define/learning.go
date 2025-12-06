package define

var (
	// UserLearningStatus 用户学习进度状态
	UserLearningStatus = struct {
		Proceed  int64
		Finish   int64
		Abnormal int64
	}{
		Proceed:  1, // 学习中
		Finish:   2, // 学习完成
		Abnormal: 3, // 异常结束
	}

	// LearningRecordStatus 学习记录状态
	LearningRecordStatus = struct {
		Proceed int64
		Finish  int64
		Break   int64
		Init    int64
	}{
		Proceed: 1, // 进行中
		Finish:  2, // 结束
		Break:   3, // 中断
		Init:    4, // 预生成
	}

	// BizType 学习业务类型
	BizType = struct {
		Question int64
		Analysis int64
	}{
		Question: 1, // 题目
		Analysis: 2, // 讲解
	}

	// OriginType 问题起源
	OriginType = struct {
		Original int64
		Append   int64
		Dynamic  int64
	}{
		Original: 1, // 原始问题
		Append:   2, // 追问问题
		Dynamic:  3, // 动态生成问题
	}

	// LearningDetailResult 学习详情结果
	LearningDetailResult = struct {
		Proceed int64
		Finish  int64
		Skip    int64
		Init    int64
	}{
		Proceed: 1, // 进行
		Finish:  2, // 完成
		Skip:    3, // 跳过
		Init:    4, // 预生成
	}
)
