package csvs

import (
	"golanggame/server/src/utils"
)

type ConfigDrop struct {
	DropId int `json:"DropId"`
	Weight int `json:"Weight"`
	Result int `json:"Result"`
	IsEnd  int `json:"IsEnd"`
}
type ConfigDropItem struct {
	DropId     int `json:"DropId"`
	DropType   int `json:"DropType"`
	Weight     int `json:"Weight"`
	ItemId     int `json:"ItemId"`
	ItemNumMin int `json:"ItemNumMin"`
	ItemNumMax int `json:"ItemNumMax"`
	WorldAdd   int `json:"WorldAdd"`
}
type ConfigStatue struct {
	StatueId    int `json:"StatueId"`
	Level       int `json:"Level"`
	CostItem    int `json:"CostItem"`
	CostNum     int `json:"CostNum"`
	RewardItem1 int `json:"RewardItem1"`
	RewardNum1  int `json:"RewardNum1"`
	RewardItem2 int `json:"RewardItem2"`
	RewardNum2  int `json:"RewardNum2"`
}

var ConfigDropSlice []*ConfigDrop
var ConfigDropItemSlice []*ConfigDropItem
var ConfigStatueSlice []*ConfigStatue

func init() {
	ConfigDropSlice = make([]*ConfigDrop, 0)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Drop", &ConfigDropSlice)
	ConfigDropItemSlice = make([]*ConfigDropItem, 0)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/DropItem", &ConfigDropItemSlice)
	ConfigStatueSlice = make([]*ConfigStatue, 0)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Statue", &ConfigStatueSlice)
	return
}
