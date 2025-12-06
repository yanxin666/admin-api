package define

import "fmt"

type ModuleConfig struct {
	Key    string
	Secret string
}

type AppConfig struct {
	Modules map[string]ModuleConfig
}

var AppConfigData = AppConfig{
	Modules: map[string]ModuleConfig{
		"passport": {Key: "muse-admin", Secret: "3097ab1b-8700-4419-9dc8-7784624da742"},
		// 可以添加更多的模块配置
	},
}

func GetModuleConfig(moduleName string) (ModuleConfig, error) {
	config, ok := AppConfigData.Modules[moduleName]
	if !ok {
		return ModuleConfig{}, fmt.Errorf("module '%s' not found", moduleName)
	}
	return config, nil
}
