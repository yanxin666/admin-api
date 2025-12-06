package define

var (
	// LLMChannel 大模型渠道
	LLMChannel = struct {
		ChannelDiFy string
	}{
		ChannelDiFy: "DiFy", // DiFy
	}

	// LLMStreamStatus 大模型流式返回消息内容
	LLMStreamStatus = struct {
		Stop     string // 正常结束
		Error    string // 异常结束
		Continue string // 跳过
	}{
		Stop:     "<Stop>",
		Error:    "<Error>",
		Continue: "<Continue>",
	}

	// LLMRole LLM角色
	LLMRole = struct {
		System    string
		Assistant string
		User      string
	}{
		System:    "system",
		Assistant: "assistant",
		User:      "user",
	}

	// LLMResourceStatus LLM大模型资源节点状态
	LLMResourceStatus = struct {
		Usable   int64
		Disabled int64
	}{
		Usable:   1, // 可用
		Disabled: 2, // 不可用
	}

	// LLMResourceType LLM大模型资源节点类型
	LLMResourceType = struct {
		GPT     int64
		Azure   int64
		FastGPT int64
		DiFy    int64
	}{
		GPT:     1,
		Azure:   2,
		FastGPT: 3,
		DiFy:    4,
	}

	// LLMResourceUseType LLM大模型资源节点用途
	LLMResourceUseType = struct {
		Answer int64
	}{
		Answer: 1, // 答题
	}

	// FormalStepSLice 公式步骤
	FormalStepSLice = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十"}
)
