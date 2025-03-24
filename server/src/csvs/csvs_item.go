package csvs

import (
	"fmt"
	"golanggame/server/src/utils"
)

type ConfigItem struct {
	ItemId   int    `json:"ItemId"`
	SortType int    `json:"SortType"`
	ItemName string `json:"ItemName"`
}

var ConfigItemMap map[int]*ConfigItem

func init() {
	ConfigItemMap = make(map[int]*ConfigItem)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Item", &ConfigItemMap)
	return
}
func GetItemConfig(itemId int) *ConfigItem {
	if ConfigItemMap == nil {
		return nil
	}
	configitem, ok := ConfigItemMap[itemId]
	if !ok {
		return nil
	}
	return configitem
}

func GetItemName(itemId int) string {
	if ConfigItemMap[itemId].ItemName == "" {
		fmt.Println("当前物品不存在")
		return ""
	}
	return ConfigItemMap[itemId].ItemName
}
