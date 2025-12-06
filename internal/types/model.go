package types

// AudioReq 音频参数
type AudioReq struct {
	Model         string         `json:"model"`          // 模型
	Text          string         `json:"text"`           // 文本
	Stream        bool           `json:"stream"`         // 是否流式：true流式/false非流式
	TimberWeights []TimberWeight `json:"timber_weights"` // 音色&权重
	VoiceSetting  VoiceSetting   `json:"voice_setting"`  // 音色设置
	AudioSetting  AudioSetting   `json:"audio_setting"`  // 音频设置
}

// TimberWeight 音色&权重
type TimberWeight struct {
	VoiceId string `json:"voice_id"` // 音色编号
	Weight  int64  `json:"weight"`   // 音色权重
}

// VoiceSetting 音色设置
type VoiceSetting struct {
	VoiceId string  `json:"voice_id"` // 音色编号
	Speed   float32 `json:"speed"`    // 语速
	Vol     float32 `json:"vol"`      // 音量
	Pitch   int64   `json:"pitch"`    // 语调
}

// AudioSetting 音频设置
type AudioSetting struct {
	AudioSampleRate int64  `json:"audio_sample_rate"` // 采样率
	Bitrate         int64  `json:"bitrate"`           // 比特率
	Format          string `json:"format"`            // 音频格式
}
