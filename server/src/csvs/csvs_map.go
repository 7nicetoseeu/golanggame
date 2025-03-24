package csvs

import (
	"golanggame/server/src/utils"
)

type ConfigMap struct {
	MapId   int    `json:"MapId"`
	MapName string `json:"MapName"`
	MapType int    `json:"MapType"`
}

var ConfigMapMap map[int]*ConfigMap

func init() {
	ConfigMapMap = make(map[int]*ConfigMap)
	utils.GetCsvUtilMgr().LoadCsv("../../csv/Map", &ConfigMapMap)
	return
}
func GetConfigMap(mapId int) *ConfigMap {
	if ConfigMapMap == nil {
		return nil
	}
	configMap, ok := ConfigMapMap[mapId]
	if !ok {
		return nil
	}
	return configMap
}
