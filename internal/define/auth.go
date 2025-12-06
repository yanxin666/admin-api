package define

// Auth相关枚举请用如下方式来定义
var (
	// AuthReqType 临时授权方式 细分粒度展开
	AuthReqType = struct {
		OSSArticleAuth int
		OSSAppAuth     int
		OSSFormulaAuth int
		OSSH5Auth      int
		ASRAuth        int
		TTSAuth        int
		AzureTTSAuth   int
	}{
		OSSArticleAuth: 1, // OSS上传文章资源临时授权
		OSSAppAuth:     2, // OSS上传APP通用资源临时授权
		OSSFormulaAuth: 3, // OSS上传Formula资源临时授权
		OSSH5Auth:      4, // OSS上传H5资源临时授权
		ASRAuth:        5, // 语音转换临时授权
		TTSAuth:        6, // TTS临时授权
		AzureTTSAuth:   7, // 微软TTS临时授权
	}

	// AppSingKey 应用签名秘钥
	AppSingKey = map[string]string{
		"muse-admin-h5": "0853f451-eed1-46c5-be52-ffa3e8758661", // AI伴学H5端签名秘钥
		"doubanjiang":   "0316753d-e248-47bd-a8c3-a2333418fd9d", // 豆伴匠
	}

	// AuthOssPath OSS上传目录
	AuthOssPath = struct {
		App     string
		Article string
		H5      string
		Formula string
		Audio   string
	}{
		App:     "app",     // 基础通用资源 例如avatar
		Article: "article", // 文章图片
		Formula: "formula", // 其他
		H5:      "h5",      // h5资源
		Audio:   "audio",   // 音频
	}
	// AuthProduct 请求passport产品表标识
	AuthProduct = struct {
		Product string
	}{
		Product: "0", // 能力探真后端项目
	}
)

const (
	ExpireTimeSeconds = 3600 // 令牌过期时间
	AzureTtsToken     = "https://%s.api.cognitive.microsoft.com/sts/v1.0/issueToken"
)
