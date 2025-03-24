package csvs

import (
	"golanggame/server/src/utils"
)

func init() {
}

type ConfigRelics struct {
	RelicsId      int `json:"RelicsId"`
	Type          int `json:"Type"`
	Pos           int `json:"Pos"`
	Star          int `json:"Star"`
	MainGroup     int `json:"MainGroup"`
	OtherGroup    int `json:"OtherGroup"`
	OtherGroupNum int `json:"OtherGroupNum"`
}
type ConfigRelicsEntry struct {
	Id        int    `json:"Id"`
	Group     int    `json:"Group"`
	AttrType  int    `json:"AttrType"`
	AttrName  string `json:"AttrName"`
	AttrValue int    `json:"AttrValue"`
	Weight    int    `json:"Weight"`
}

type ConfigRelicsLevel struct {
	EntryId   int    `json:"EntryId"`
	Level     int    `json:"Level"`
	AttrType  int    `json:"AttrType"`
	AttrName  string `json:"AttrName"`
	AttrValue int    `json:"AttrValue"`
	NeedExp   int    `json:"NeedExp"`
}
type ConfigRelicsSuit struct {
	Type        int    `json:"Type"`
	Num         int    `json:"Num"`
	SuitSkill   int    `json:"SuitSkill"`
	SkillString string `json:"SkillString"`
}

var ConfigRelicsMap map[int]*ConfigRelics
var ConfigRelicsEntryMap map[int]*ConfigRelicsEntry
var ConfigRelicsLevelSlice []*ConfigRelicsLevel
var ConfigRelicsEntrySlice []*ConfigRelicsEntry
var ConfigRelicsSuitSlice []*ConfigRelicsSuit

func init() {
	ConfigRelicsMap = make(map[int]*ConfigRelics)
	ConfigRelicsEntryMap = make(map[int]*ConfigRelicsEntry)
	ConfigRelicsEntrySlice = make([]*ConfigRelicsEntry, 0)
	ConfigRelicsLevelSlice = make([]*ConfigRelicsLevel, 0)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Relics", &ConfigRelicsMap)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/RelicsEntry", &ConfigRelicsEntrySlice)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/RelicsEntry", &ConfigRelicsEntryMap)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/RelicsLevel", &ConfigRelicsLevelSlice)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/RelicsSuit", &ConfigRelicsSuitSlice)
	return
}
func GetConfigRelics(itemId int) *ConfigRelics {
	if ConfigRelicsMap == nil {
		return nil
	}
	ConfigRelics, ok := ConfigRelicsMap[itemId]
	if !ok {
		return nil
	}
	return ConfigRelics
}
