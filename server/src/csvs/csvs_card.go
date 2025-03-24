package csvs

import (
	"golanggame/server/src/utils"
)

func init() {
}

type ConfigCard struct {
	CardId       int `json:"CardId"`
	Friendliness int `json:"Friendliness"`
	Check        int `json:"Check"`
}

var ConfigCardMap map[int]*ConfigCard
var ConfigCardByRoleIdMap map[int]*ConfigCard

func init() {
	ConfigCardMap = make(map[int]*ConfigCard)
	ConfigCardByRoleIdMap = make(map[int]*ConfigCard)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Card", &ConfigCardMap)
	for _, v := range ConfigCardMap {
		ConfigCardByRoleIdMap[v.Check] = v
	}
	return
}
func GetConfigCard(itemId int) *ConfigCard {
	if ConfigCardMap == nil {
		return nil
	}
	configCard, ok := ConfigCardMap[itemId]
	if !ok {
		return nil
	}
	return configCard
}
func GetConfigCardByRoleId(RoleId int) *ConfigCard {
	if ConfigCardMap == nil {
		return nil
	}
	configCard, ok := ConfigCardByRoleIdMap[RoleId]
	if !ok {
		return nil
	}
	return configCard
}
