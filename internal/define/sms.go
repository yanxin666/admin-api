package define

const (
	// SmsLoginTemplateKey 登录验证短信模板ID，模板内容：验证码：{1}，有效期{2}分钟。如非本人操作，请忽略。
	SmsLoginTemplateKey = "zmexing_tpl1"
	SmsLoginSign        = "zmexing"

	// SmsLoginTemplateId 尊敬的用户，您的{1}（权益名称，如快解题一年）使用权限已生效，请用本号码登录辞源App，有问题请联系客服，感谢您的使用。
	SmsLoginTemplateId = "2239202"

	// SmsLoginValidityTime 登陆验证码的有效期
	SmsLoginValidityTime = "5"
)

// WhiteList 短信白名单
var WhiteList = []string{
	"15235656647", // 王鹏飞
}
