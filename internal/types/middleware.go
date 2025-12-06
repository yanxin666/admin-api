package types

import "net/http"

// MiddlewareCommonParams 通用参数，为什么如下几个枚举值也需要用string？是因为安卓那边转Map的时候只支持string，这里只能手动兼容设置为string...
type MiddlewareCommonParams struct {
	Product       string `json:"product"`       // 产品线 0.辞源 1.AI搜索 2.白帝辞 3.智普乌托邦
	Source        string `json:"source"`        // 来源 0.APP 1.豆伴匠 2.风的颜色 1001.听力熊 1002.微软
	System        string `json:"system"`        // 操作系统 1.安卓 2.iOS 3.H5 4.鸿蒙
	Screen        string `json:"screen"`        // 设备分辨率
	Imei          string `json:"imei"`          // 设备标识，例：手机设备号
	Dev           string `json:"dev"`           // 具体设备型号，例：iPhone14,5
	SysVersion    string `json:"sysVersion"`    // 系统版本，例：iOS17.5
	AppVersion    string `json:"appVersion"`    // app版本，例：1.0.99
	PackageName   string `json:"packageName"`   // APP包名
	From          string `json:"from"`          // 渠道 暂时预留
	NetEnv        string `json:"netEnv"`        // 网络，例：WiFi
	Platform      string `json:"platform"`      // 设备类别 1.iPhone 2.iPad 3.Android 4.平板 5.学习机 6.演示版
	ChannelDevice string `json:"channelDevice"` // 渠道设备，与 Platform 连用，例：学习机设备号
}

// MiddlewareApiRequestLog 请求日志
type MiddlewareApiRequestLog struct {
	Host     string         `json:"host"`
	ClientIP string         `json:"client_ip"`
	Schema   string         `json:"schema"`
	Header   http.Header    `json:"header"`
	URL      string         `json:"url"`
	Path     string         `json:"path"`
	Method   string         `json:"method"`
	Params   map[string]any `json:"params"`
	UserId   int64          `json:"user_id"`
}
