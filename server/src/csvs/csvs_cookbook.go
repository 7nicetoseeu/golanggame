package csvs

import (
	"golanggame/server/src/utils"
)

func init() {
	// fmt.Println("初始化菜谱配置")
}

type ConfigCookBook struct {
	CookBookId int `json:"CookBookId"`
	Reward     int `json:"Reward"`
}

var ConfigCookBookMap map[int]*ConfigCookBook

func init() {
	ConfigCookBookMap = make(map[int]*ConfigCookBook)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/CookBook", &ConfigCookBookMap)
	return
}
func GetConfigCookBook(itemId int) *ConfigCookBook {
	if ConfigCookBookMap == nil {
		return nil
	}
	ConfigCookBook, ok := ConfigCookBookMap[itemId]
	if !ok {
		return nil
	}
	return ConfigCookBook
}
