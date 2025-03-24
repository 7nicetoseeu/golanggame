package csvs

import (
	"golanggame/server/src/utils"
)

func init() {
}

type ConfigEvent struct {
	EventId        int    `json:"EventId"`
	EventType      int    `json:"EventType"`
	RefreshType    int    `json:"RefreshType"`
	Name           string `json:"Name"`
	EventDrop      int    `json:"EventDrop"`
	EventDropTimes int    `json:"EventDropTimes"`
	MapId          int    `json:"MapId"`
	CostItem       int    `json:"CostItem"`
	CostNum        int    `json:"CostNum"`
}

var ConfigEventMap map[int]*ConfigEvent

func init() {
	ConfigEventMap = make(map[int]*ConfigEvent)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/MapEvent", &ConfigEventMap)
	return
}
func GetConfigEvent(EventId int) *ConfigEvent {
	if ConfigEventMap == nil {
		return nil
	}
	ConfigEvent, ok := ConfigEventMap[EventId]
	if !ok {
		return nil
	}
	return ConfigEvent
}
