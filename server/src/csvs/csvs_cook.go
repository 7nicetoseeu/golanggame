package csvs

import (
	"golanggame/server/src/utils"
)

func init() {
	// fmt.Println("初始化菜谱技能配置")
}

type ConfigCook struct {
	CookId int `json:"CookId"`
	Star   int `json:"Star"`
}

var ConfigCookMap map[int]*ConfigCook

func init() {
	ConfigCookMap = make(map[int]*ConfigCook)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Cook", &ConfigCookMap)
	return
}
func GetConfigCook(itemId int) *ConfigCook {
	if ConfigCookMap == nil {
		return nil
	}
	ConfigCook, ok := ConfigCookMap[itemId]
	if !ok {
		return nil
	}
	return ConfigCook
}
