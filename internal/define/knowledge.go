package define

var (
	// UserLessonStatus 用户排课状态
	UserLessonStatus = struct {
		// 已学习
		Studied int64
		// 未学习
		UnStudied int64
	}{
		Studied:   2, // 已学习
		UnStudied: 1, // 未学习
	}

	ReviewStatus = struct {
		InitialEntry int64 // 初始录入
		Reviewing    int64 // 审核中
		Passed       int64 // 审核通过
		NotPassed    int64 // 审核未通过
		OffShelf     int64 // 下架
	}{
		InitialEntry: 1,
		Reviewing:    2,
		Passed:       3,
		NotPassed:    4,
		OffShelf:     5,
	}

	LessonUsageType = struct {
		Example  int64 // 例题
		Question int64 // 练习题
	}{
		Example:  1,
		Question: 2,
	}

	// LessonType 课节类型
	LessonType = struct {
		Class              int64 // 主线课 包含 导入时的大小语文
		MonthlyExamination int64 // 月考
		SpecialCourse      int64 // 小灶课
	}{
		Class:              1,
		MonthlyExamination: 2,
		SpecialCourse:      3,
	}

	// NodeType 节点类型
	NodeType = struct {
		SpecialCourse int64 // 小灶课[原理讲解]
		SmallLanguage int64 // 小语文[目标课]
		BigLanguage   int64 // 大语文[目标课]
	}{
		SpecialCourse: 1,
		SmallLanguage: 2,
		BigLanguage:   3,
	}

	// LessonTypeMapping 导入时对业务的映射
	LessonTypeMapping = map[int64]int64{
		NodeType.SmallLanguage: LessonType.Class,
		NodeType.BigLanguage:   LessonType.Class,
		NodeType.SpecialCourse: LessonType.SpecialCourse,
	}

	// LessonTypeMap 导入时对业务的映射
	LessonTypeMap = map[int64]string{
		NodeType.SmallLanguage: "小语文[目标课]",
		NodeType.BigLanguage:   "大语文[目标课]",
		NodeType.SpecialCourse: "小灶课[原理讲解]",
	}
)
