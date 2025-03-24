package csvs

import "golanggame/server/src/utils"

type ConfigRole struct {
	RoleId          int    `json:"RoleId"`
	ItemName        string `json:"ItemName"`
	Star            int    `json:"Star"`
	Stuff           int    `json:"Stuff"`
	StuffNum        int64  `json:"StuffNum"`
	StuffItem       int    `json:"StuffItem"`
	StuffItemNum    int64  `json:"StuffItemNum"`
	MaxStuffItem    int    `json:"MaxStuffItem"`
	MaxStuffItemNum int64  `json:"MaxStuffItemNum"`
}

var ConfigRoleMap map[int]*ConfigRole

func init() {
	ConfigRoleMap = make(map[int]*ConfigRole)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Role", &ConfigRoleMap)
	return
}
func GetRoleConfig(itemId int) *ConfigRole {
	if ConfigRoleMap == nil {
		return nil
	}
	ConfigRole, ok := ConfigRoleMap[itemId]
	if !ok {
		return nil
	}
	return ConfigRole
}
