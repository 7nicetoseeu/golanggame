package csvs

import (
	"golanggame/server/src/utils"
)

type ConfigWeapon struct {
	WeaponId int `json:"WeaponId"`
	Type     int `json:"Type"`
	Star     int `json:"Star"`
}
type ConfigWeaponLevel struct {
	WeaponStar    int `json:"WeaponStar"`
	Level         int `json:"Level"`
	NeedExp       int `json:"NeedExp"`
	NeedStarLevel int `json:"NeedStarLevel"`
}
type ConfigWeaponStar struct {
	WeaponStar int `json:"WeaponStar"`
	StarLevel  int `json:"StarLevel"`
	Level      int `json:"Level"`
}

var ConfigWeaponMap map[int]*ConfigWeapon
var ConfigWeaponLevelSlice []*ConfigWeaponLevel
var ConfigWeaponStarSlice []*ConfigWeaponStar

func init() {
	ConfigWeaponMap = make(map[int]*ConfigWeapon)
	ConfigWeaponLevelSlice = make([]*ConfigWeaponLevel, 0)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Weapon", &ConfigWeaponMap)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/WeaponLevel", &ConfigWeaponLevelSlice)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/WeaponStar", &ConfigWeaponStarSlice)
	return
}
func GetConfigWeapon(itemId int) *ConfigWeapon {
	if ConfigWeaponMap == nil {
		return nil
	}
	ConfigWeapon, ok := ConfigWeaponMap[itemId]
	if !ok {
		return nil
	}
	return ConfigWeapon
}
