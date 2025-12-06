package define

var (
	// DifficultyUpgrade 能力提升难度
	DifficultyUpgrade = struct {
		Easy       int64
		Middle     int64
		Difficulty int64
	}{
		Easy:       1, // 进行
		Middle:     2, // 完成
		Difficulty: 3, // 中断
	}

	// DifficultyArray 能力提升难度数组
	DifficultyArray = []int64{DifficultyUpgrade.Easy, DifficultyUpgrade.Middle, DifficultyUpgrade.Difficulty}

	// AbilityLevel 能力等级
	AbilityLevel = struct {
		Doubt     int64
		Excellent int64
		Qualified int64
		Lack      int64
	}{
		Doubt:     0, // 存疑
		Excellent: 1, // 优秀
		Qualified: 2, // 合格
		Lack:      3, // 欠缺
	}

	// AbilityLevelMap 能力等级名称映射
	AbilityLevelMap = map[int64]string{
		AbilityLevel.Doubt:     "存疑",
		AbilityLevel.Excellent: "优秀",
		AbilityLevel.Qualified: "合格",
		AbilityLevel.Lack:      "欠缺",
	}

	// AbilityLevelArray 能力等级数组
	AbilityLevelArray = []int64{AbilityLevel.Doubt, AbilityLevel.Excellent, AbilityLevel.Qualified, AbilityLevel.Lack}

	// AbilityBelongTo 能力所属
	AbilityBelongTo = struct {
		T string
		B string
		G string
		D string
	}{
		T: "T", // 技巧
		B: "B", // 知识
		G: "G", // 共情
		D: "D", // 表达
	}

	// AbilityBelongToArray 能力所属数组
	AbilityBelongToArray = []string{AbilityBelongTo.T, AbilityBelongTo.B, AbilityBelongTo.G, AbilityBelongTo.D}

	// WordsLevel 词汇掌握度
	WordsLevel = struct {
		UnMastered int64
		Mastered   int64
	}{
		UnMastered: 0, // 未掌握
		Mastered:   1, // 已掌握
	}

	// WordsRegion B1词汇图谱 4大区
	WordsRegion = []string{
		"日常高频",
		"课内知识表",
		"文学高频",
		"文言文",
	}

	WordsTitle = "B1 语法词汇图谱"
)
