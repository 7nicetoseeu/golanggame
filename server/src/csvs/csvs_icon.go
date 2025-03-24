package csvs

import (
	"golanggame/server/src/utils"
)

type ConfigIcon struct {
	IconId int `json:"IconId"`
	Check  int `json:"Check"`
}

var ConfigIconMap map[int]*ConfigIcon
var ConfigIconMapByRoleId map[int]*ConfigIcon

func init() {
	ConfigIconMap = make(map[int]*ConfigIcon)
	ConfigIconMapByRoleId = make(map[int]*ConfigIcon)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Icon", &ConfigIconMap)
	for _, v := range ConfigIconMap {
		ConfigIconMapByRoleId[v.Check] = v
	}
	return
}
func GetConfigIcon(itemId int) *ConfigIcon {
	if ConfigIconMap == nil {
		return nil
	}
	ConfigIcon, ok := ConfigIconMapByRoleId[itemId]
	if !ok {
		return nil
	}
	return ConfigIcon
}
