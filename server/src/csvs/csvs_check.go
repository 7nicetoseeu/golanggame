package csvs

import (
	"fmt"
	"math/rand"
)

type DropGrop struct {
	DropId     int
	WeightAll  int
	ConfigDrop []*ConfigDrop
}
type DropGropItem struct {
	ItemId         int
	ConfigDropItem []*ConfigDropItem
}
type RelicsEntry struct {
	RlicsId           int
	ConfigRelicsEntry []*ConfigRelicsEntry
}
type RelicsSuit struct {
	SuitId           int
	ConfigRelicsSuit []*ConfigRelicsSuit
}

var DropGropConfigMap map[int]*DropGrop
var DropGropItemConfigMap map[int]*DropGropItem
var RlicsEntryConfigMap map[int]*RelicsEntry
var ConfigStatueMap map[int]map[int]*ConfigStatue
var RelicsLevelConfigMap map[int]map[int]*ConfigRelicsLevel
var RelicsSuitConfigMap map[int]map[int]*ConfigRelicsSuit
var WeaponLevelConfigGroup map[int]map[int]*ConfigWeaponLevel
var WeaponStarConfigGroup map[int]map[int]*ConfigWeaponStar

func MakeStatueConfig() {
	ConfigStatueMap = make(map[int]map[int]*ConfigStatue, 0)
	for _, configStatue := range ConfigStatueSlice {
		statueMap, ok := ConfigStatueMap[configStatue.StatueId]
		if !ok {
			statueMap = make(map[int]*ConfigStatue)
			ConfigStatueMap[configStatue.StatueId] = statueMap
		}
		statueMap[configStatue.Level] = configStatue
	}
	return
}

func MakeDropMapConfig() {
	DropGropConfigMap = make(map[int]*DropGrop)
	for _, configDrop := range ConfigDropSlice {
		dropGrop, ok := DropGropConfigMap[configDrop.DropId]
		if !ok {
			DropGropConfigMap[configDrop.DropId] = new(DropGrop)
			DropGropConfigMap[configDrop.DropId].DropId = configDrop.DropId
			DropGropConfigMap[configDrop.DropId].WeightAll = configDrop.Weight
			DropGropConfigMap[configDrop.DropId].ConfigDrop = append(DropGropConfigMap[configDrop.DropId].ConfigDrop, configDrop)
		} else {
			dropGrop.WeightAll += configDrop.Weight
			dropGrop.ConfigDrop = append(dropGrop.ConfigDrop, configDrop)
		}
	}
}

func MakeDropItemMapConfig() {
	DropGropItemConfigMap = make(map[int]*DropGropItem)
	for _, ConfigDropItem := range ConfigDropItemSlice {
		dropGropItem, ok := DropGropItemConfigMap[ConfigDropItem.DropId]
		if !ok {
			DropGropItemConfigMap[ConfigDropItem.DropId] = new(DropGropItem)
			DropGropItemConfigMap[ConfigDropItem.DropId].ItemId = ConfigDropItem.ItemId
			DropGropItemConfigMap[ConfigDropItem.DropId].ConfigDropItem = append(DropGropItemConfigMap[ConfigDropItem.DropId].ConfigDropItem, ConfigDropItem)
		} else {
			dropGropItem.ConfigDropItem = append(dropGropItem.ConfigDropItem, ConfigDropItem)
		}
	}
}

func RandDropItmeTest() *ConfigDrop {
	dropGropItem := DropGropItemConfigMap[1]
	if dropGropItem == nil {
		fmt.Println("掉落空")
		return nil
	}
	randNum := rand.Intn(DORP_ITEM_ALL_PERCENT)
	for _, configDropItem := range dropGropItem.ConfigDropItem {
		if randNum < configDropItem.Weight {
			fmt.Println("掉落", configDropItem.ItemId)
		}
	}
	return nil
}
func GetDropGropItem(DropId int) *DropGropItem {
	return DropGropItemConfigMap[DropId]
}

//递归思想
func GetDropGropItemNew(dropId int) []*ConfigDropItem {
	rel := make([]*ConfigDropItem, 0)
	DropGropItem := DropGropItemConfigMap[dropId]
	if DropGropItem == nil {
		return nil
	}
	allItem := make([]*ConfigDropItem, 0)
	for _, v := range DropGropItem.ConfigDropItem {
		if v.DropType == DORP_ITEM_TYPE_GROUP {
			randNum := rand.Intn(DORP_ITEM_ALL_PERCENT)
			if randNum > v.Weight {
				continue
			}
			dropItems := GetDropGropItemNew(v.ItemId)
			rel = append(rel, dropItems...)
		} else if v.DropType == DORP_ITEM_TYPE_ITEM {
			rel = append(rel, v)
		} else if v.DropType == DORP_ITEM_TYPE_WEIGHT {
			allItem = append(allItem, v)
		}
	}
	allRate := 0
	if len(allItem) != 0 {
		for _, v := range allItem {
			allRate += v.Weight
		}
		randNum := rand.Intn(allRate)
		randNow := 0
		for _, v := range allItem {
			randNow += v.Weight
			if randNow > randNum {
				configItem := new(ConfigDropItem)
				*configItem = *v
				configItem.Weight = allRate
				rel = append(rel, configItem)
				return rel
			}
		}
	}
	return rel
}
func RandDropTest() *ConfigDrop {
	dropGrop := DropGropConfigMap[1000]
	if dropGrop == nil {
		return nil
	}
	for {
		dropGropConfig := GetDropGrop(dropGrop)
		if dropGropConfig.IsEnd == LOGIC_TRUE {
			return dropGropConfig
		}
		dropGrop = DropGropConfigMap[dropGropConfig.Result]
		if dropGrop == nil {
			break
		}
	}
	return nil
}
func GetDropGrop(dropGrop *DropGrop) *ConfigDrop {
	// time.Sleep(1 * time.Second)
	// rand.Seed(time.Now().Unix())
	randNum := rand.Intn(dropGrop.WeightAll)
	randNow := 0
	for _, grop := range dropGrop.ConfigDrop {
		randNow += grop.Weight
		if randNow > randNum {
			return grop
		}
	}
	return nil
}

//递归思想
func GetDropGropNew(dropGrop *DropGrop) *ConfigDrop {
	randNum := rand.Intn(dropGrop.WeightAll)
	randNow := 0
	for _, grop := range dropGrop.ConfigDrop {
		randNow += grop.Weight
		if randNow > randNum {
			if grop.IsEnd == LOGIC_TRUE {
				return grop
			}
			dropGrop := DropGropConfigMap[grop.Result]
			if dropGrop == nil {
				return nil
			}
			return GetDropGropNew(dropGrop)
		}
	}
	return nil
}
func GetDropGropNew1(dropGrop *DropGrop, fiveInfo map[int]int, fourInfo map[int]int) *ConfigDrop {

	for _, v := range dropGrop.ConfigDrop {
		_, ok := fiveInfo[v.Result]
		if ok {
			maxGetTimes := 0
			index := new(ConfigDrop)
			for _, v1 := range dropGrop.ConfigDrop {
				times, newOk := fiveInfo[v1.Result]
				if !newOk {
					continue
				}
				if maxGetTimes < times {
					maxGetTimes = times
					index = v1
				}

			}
			return index
		}
		_, ok = fourInfo[v.Result]
		if ok {
			maxGetTimes := 0
			index := new(ConfigDrop)
			for _, v1 := range dropGrop.ConfigDrop {
				times, newOk := fourInfo[v1.Result]
				if !newOk {
					continue
				}
				if maxGetTimes < times {
					maxGetTimes = times
					index = v1
				}

			}
			return index
		}
	}

	randNum := rand.Intn(dropGrop.WeightAll)
	randNow := 0
	for _, grop := range dropGrop.ConfigDrop {
		randNow += grop.Weight
		if randNow > randNum {
			if grop.IsEnd == LOGIC_TRUE {
				return grop
			}
			dropGrop := DropGropConfigMap[grop.Result]
			if dropGrop == nil {
				return nil
			}
			return GetDropGropNew(dropGrop)
		}
	}
	return nil
}
func MakeRelicsEntryMapConfig() {
	RlicsEntryConfigMap = make(map[int]*RelicsEntry)
	for _, v := range ConfigRelicsEntrySlice {
		_, ok := RlicsEntryConfigMap[v.Group]
		if !ok {
			RlicsEntryConfigMap[v.Group] = new(RelicsEntry)
			RlicsEntryConfigMap[v.Group].RlicsId = v.Group
			RlicsEntryConfigMap[v.Group].ConfigRelicsEntry = append(RlicsEntryConfigMap[v.Group].ConfigRelicsEntry, v)
		} else {
			RlicsEntryConfigMap[v.Group].ConfigRelicsEntry = append(RlicsEntryConfigMap[v.Group].ConfigRelicsEntry, v)
		}
	}
}
func MakeRelicsLevelMapConfig() {
	RelicsLevelConfigMap = make(map[int]map[int]*ConfigRelicsLevel)
	for _, v := range ConfigRelicsLevelSlice {
		_, ok := RelicsLevelConfigMap[v.EntryId]
		if !ok {
			RelicsLevelConfigMap[v.EntryId] = make(map[int]*ConfigRelicsLevel)
			RelicsLevelConfigMap[v.EntryId][v.Level] = v
		} else {
			RelicsLevelConfigMap[v.EntryId][v.Level] = v
		}
	}
}
func MakeRelicsSuitMapConfig() {
	RelicsSuitConfigMap = make(map[int]map[int]*ConfigRelicsSuit)
	for _, v := range ConfigRelicsSuitSlice {
		_, ok := RelicsSuitConfigMap[v.Type]
		if !ok {
			RelicsSuitConfigMap[v.Type] = make(map[int]*ConfigRelicsSuit)
			RelicsSuitConfigMap[v.Type][v.Num] = v
		} else {
			RelicsSuitConfigMap[v.Type][v.Num] = v
		}
	}
}
func MakeWeaponLevelGourpConfig() {
	WeaponLevelConfigGroup = make(map[int]map[int]*ConfigWeaponLevel)
	for _, v := range ConfigWeaponLevelSlice {
		_, ok := WeaponLevelConfigGroup[v.WeaponStar]
		if !ok {
			WeaponLevelConfigGroup[v.WeaponStar] = make(map[int]*ConfigWeaponLevel)
			WeaponLevelConfigGroup[v.WeaponStar][v.Level] = v
		} else {
			WeaponLevelConfigGroup[v.WeaponStar][v.Level] = v
		}
	}
}
func MakeWeaponStarGourpConfig() {
	WeaponStarConfigGroup = make(map[int]map[int]*ConfigWeaponStar)
	for _, v := range ConfigWeaponStarSlice {
		_, ok := WeaponStarConfigGroup[v.WeaponStar]
		if !ok {
			WeaponStarConfigGroup[v.WeaponStar] = make(map[int]*ConfigWeaponStar)
			WeaponStarConfigGroup[v.WeaponStar][v.StarLevel] = v
		} else {
			WeaponStarConfigGroup[v.WeaponStar][v.StarLevel] = v
		}
	}
}

//主词条配置
func GetRelicsLevelConfig(mainEntry int, level int) *ConfigRelicsLevel {
	if ConfigRelicsEntryMap[mainEntry] == nil {
		return nil
	}
	AttrTypedId := ConfigRelicsEntryMap[mainEntry].AttrType
	m1 := RelicsLevelConfigMap[AttrTypedId]
	if m1 == nil {
		return nil
	}
	relicsLevelConfig := RelicsLevelConfigMap[AttrTypedId][level]
	if relicsLevelConfig == nil {
		return nil
	}
	// relicsLevelConfig.AttrValue += ConfigRelicsEntryMap[mainEntry].AttrValue
	return relicsLevelConfig
}

//副词条配置
func GetRelicsOtherConfig(otherEntry int, level int) *ConfigRelicsEntry {
	configRelicsEntry := ConfigRelicsEntryMap[otherEntry]
	m1 := RelicsLevelConfigMap[configRelicsEntry.AttrType]
	if m1 == nil {
		return nil
	}
	return configRelicsEntry
}
func CheckCsvs() {
	MakeDropMapConfig()
	MakeDropItemMapConfig()
	MakeStatueConfig()
	MakeRelicsEntryMapConfig()
	MakeRelicsLevelMapConfig()
	MakeRelicsSuitMapConfig()
	MakeWeaponLevelGourpConfig()
	MakeWeaponStarGourpConfig()
	fmt.Println("初始化配置完成")
}
