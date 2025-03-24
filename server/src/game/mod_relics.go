package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"math/rand"
	"os"
)

type RelicsInfo struct {
	RelicsId   int
	RelicsName string
	KeyId      int
	MainEntry  int
	OtherEntry []int
	Level      int
	Exp        int
	RoleId     int
}

type ModRelics struct {
	ModRelicsMap map[int]*RelicsInfo
	ModKeyId     int

	player *Player
	path   string
}

func (self *ModRelics) SaveData() {
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, data, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModRelics) InitData() {
	if self.ModRelicsMap == nil {
		self.ModRelicsMap = make(map[int]*RelicsInfo)
	}
}
func (self *ModRelics) LoadData(player *Player) {
	self.player = player
	self.path = player.LocalPath + "/relics.json"
	self.InitData()
	data, err := ioutil.ReadFile(self.path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &self)
	if err != nil {
		return
	}
}
func (self *ModRelics) AddRelics(relicsId int, num int) {
	if len(self.ModRelicsMap)+num >= csvs.RELICS_MAX_COUNT {
		fmt.Println("圣遗物超过最大数量")
		return
	}
	for i := 0; i < num; i++ {
		RelicsInfo := self.NewRelics(relicsId)
		RelicsInfo.ShowInfo()
		self.ModRelicsMap[self.ModKeyId] = RelicsInfo
	}
}
func (self *ModRelics) ShowRelicsBag() {
	if len(self.ModRelicsMap) == 0 {
		fmt.Println("当前背包为空")
		return
	}
	for _, v := range self.ModRelicsMap {
		v.ShowInfo()
	}
}
func (self *ModRelics) NewRelics(relicsId int) *RelicsInfo {
	RelicsConfig := csvs.GetConfigRelics(relicsId)
	itemConfig := csvs.GetItemConfig(relicsId)
	self.ModKeyId++
	RelicsInfo := &RelicsInfo{
		RelicsId:   RelicsConfig.RelicsId,
		RelicsName: itemConfig.ItemName,
		KeyId:      self.ModKeyId,
		// Level:      1,
		MainEntry: self.MakeMainEntry(RelicsConfig.MainGroup),
	}
	for i := 0; i < RelicsConfig.OtherGroupNum; i++ {
		if i == RelicsConfig.OtherGroupNum-1 {
			//只有百分之20的几率会生成第四个词条
			randNum := rand.Intn(csvs.DORP_ITEM_ALL_PERCENT)
			if randNum < csvs.OTHER_FOUR_ENTRY_YES {
				otherEntry := self.MakeOtherEntry(RelicsConfig.OtherGroup, RelicsInfo)
				RelicsInfo.OtherEntry = append(RelicsInfo.OtherEntry, otherEntry)
			}
		} else {
			otherEntry := self.MakeOtherEntry(RelicsConfig.OtherGroup, RelicsInfo)
			RelicsInfo.OtherEntry = append(RelicsInfo.OtherEntry, otherEntry)
		}
	}
	return RelicsInfo
}
func (self *ModRelics) MakeOtherEntry(otherEntry int, relicesInfo *RelicsInfo) int {
	allEntry := make(map[int]int, 0)
	allEntry[csvs.ConfigRelicsEntryMap[relicesInfo.MainEntry].AttrType] = csvs.LOGIC_TRUE
	for _, v := range relicesInfo.OtherEntry {
		allEntry[csvs.ConfigRelicsEntryMap[v].AttrType] = csvs.LOGIC_TRUE
	}

	configEntry := csvs.RlicsEntryConfigMap[otherEntry]
	if configEntry == nil {
		return 0
	}
	allRate := 0
	for _, v := range configEntry.ConfigRelicsEntry {
		allRate += v.Weight
	}
	configRelics := csvs.ConfigRelicsMap[relicesInfo.RelicsId]
	if len(relicesInfo.OtherEntry) >= configRelics.OtherGroupNum {
		for {
			randNum := rand.Intn(allRate)
			randNow := 0
			for _, v := range configEntry.ConfigRelicsEntry {
				randNow += v.Weight
				if randNum < randNow {
					if allEntry[csvs.ConfigRelicsEntryMap[v.Id].AttrType] == 1 {
						return v.Id
					} else {
						break
					}
				}
			}
		}
	} else {
		for {
			randNum := rand.Intn(allRate)
			randNow := 0
			for _, v := range configEntry.ConfigRelicsEntry {
				randNow += v.Weight
				if randNum < randNow {
					if allEntry[csvs.ConfigRelicsEntryMap[v.Id].AttrType] == 0 {
						return v.Id
					} else {
						break
					}
				}
			}
		}
	}
}
func (self *ModRelics) TestMakeOtherEntry() {
	testCount := 27372
	entryTypeMap := make(map[int]int, 0)

	for i := 0; i < testCount; i++ {
		relicsInfo := self.NewRelics(7000005)
		for _, v := range relicsInfo.OtherEntry {
			entryTypeMap[v]++
		}
	}
	for index, v := range entryTypeMap {
		fmt.Println(index, "--", v/4)
	}
}
func (self *ModRelics) MakeMainEntry(mainGroup int) int {
	configRelics := csvs.RlicsEntryConfigMap[mainGroup]
	if configRelics == nil {
		return 0
	}
	allRate := 0
	for _, v := range configRelics.ConfigRelicsEntry {
		allRate += v.Weight
	}
	randNum := rand.Intn(allRate)
	randNow := 0
	for _, v := range configRelics.ConfigRelicsEntry {
		randNow += v.Weight
		if randNum < randNow {
			return v.Id
		}
	}
	return 0
}
func (self *ModRelics) TestMakeMainEntry() {
	testCount := 10000
	entryTypeMap := make(map[int]int, 0)
	for i := 0; i < testCount; i++ {
		entryId := self.MakeMainEntry(1)
		entryType := csvs.ConfigRelicsEntryMap[entryId].AttrType
		entryTypeMap[entryType]++
	}
	for index, v := range entryTypeMap {
		fmt.Println(index, "--", v)
	}
}
func (self *RelicsInfo) ShowInfo() {
	fmt.Println("圣遗物", self.RelicsName, "圣遗物编号---", self.KeyId)
	relicesEntry := csvs.GetRelicsLevelConfig(self.MainEntry, self.Level)
	if relicesEntry != nil {
		fmt.Println("主词条：", relicesEntry.AttrName, "值：", relicesEntry.AttrValue)
		for index, v := range self.OtherEntry {
			relicesOtherEntry := csvs.GetRelicsOtherConfig(v, self.Level)
			if relicesOtherEntry == nil {
				fmt.Printf("relicesOtherEntry: %v\n", relicesOtherEntry)
			}
			if index >= 4 {
				fmt.Println("附加副词条：", relicesOtherEntry.AttrName, "值：", relicesOtherEntry.AttrValue)
			} else {
				fmt.Println("副词条：", relicesOtherEntry.AttrName, "值：", relicesOtherEntry.AttrValue)
			}
		}
	}
}
func (self *ModRelics) UpRelicsLevel(player *Player, relicsKeyId int) {
	relices := self.ModRelicsMap[relicsKeyId]
	if relices == nil {
		return
	}
	mainEntry := relices.MainEntry
	configRelicsEntry := csvs.ConfigRelicsEntryMap[mainEntry]
	if configRelicsEntry == nil {
		return
	}
	for {
		relicesLevelConfig := csvs.RelicsLevelConfigMap[configRelicsEntry.AttrType][relices.Level+1]
		if relicesLevelConfig == nil {
			fmt.Println("已满级")
			break
		}
		if relices.Exp < relicesLevelConfig.NeedExp {
			fmt.Println("经验不够升级，当前经验", relices.Exp)
			break
		}
		relices.Exp -= relicesLevelConfig.NeedExp
		relices.Level++
		if relices.Level%4 == 0 {
			configRelics := csvs.ConfigRelicsMap[relices.RelicsId]
			otherEntry := self.MakeOtherEntry(configRelics.OtherGroup, relices)
			relices.OtherEntry = append(relices.OtherEntry, otherEntry)
		}
		fmt.Println("升级成功，当前圣遗物等级", relices.Level)
	}
	self.ModRelicsMap[relicsKeyId].ShowInfo()
}

//不需要经验直接升级
func (self *ModRelics) UpRelicsLevelTest(player *Player, relicsKeyId int) {
	relices := self.ModRelicsMap[relicsKeyId]
	if relices == nil {
		return
	}
	mainEntry := relices.MainEntry
	configRelicsEntry := csvs.ConfigRelicsEntryMap[mainEntry]
	if configRelicsEntry == nil {
		return
	}
	for {
		relicesLevelConfig := csvs.RelicsLevelConfigMap[configRelicsEntry.AttrType][relices.Level+1]
		if relicesLevelConfig == nil {
			fmt.Println("已满级")
			break
		}
		relices.Level++
		if relices.Level%4 == 0 {
			configRelics := csvs.ConfigRelicsMap[relices.RelicsId]
			otherEntry := self.MakeOtherEntry(configRelics.OtherGroup, relices)
			relices.OtherEntry = append(relices.OtherEntry, otherEntry)
		}
		fmt.Println("升级成功，当前圣遗物等级", relices.Level)
	}
	self.ModRelicsMap[relicsKeyId].ShowInfo()
}
func (self *ModRelics) AddRelicsExp(player *Player, relicsKeyId int) {
	self.ModRelicsMap[relicsKeyId].Exp += 10000
}
func (self *ModRelics) TestBestRelics() {
	dropCount := 100000
	bestHeadRelics := make(map[int]*RelicsInfo, 0)
	for i := 0; i < dropCount; i++ {
		fmt.Println("第", i, "次")
		relicsInfo := self.NewRelics(7000005)
		mainEntry := relicsInfo.MainEntry
		mainType := csvs.ConfigRelicsEntryMap[mainEntry].AttrType
		if mainType != 4 && mainType != 5 {
			continue
		}
		configRelicsEntry := csvs.ConfigRelicsEntryMap[mainEntry]
		if configRelicsEntry == nil {
			return
		}
		isadd := 1
		for {
			relicesLevelConfig := csvs.RelicsLevelConfigMap[configRelicsEntry.AttrType][relicsInfo.Level+1]
			if relicesLevelConfig == nil {
				break
			}
			relicsInfo.Level++
			if relicsInfo.Level%4 == 0 {
				configRelics := csvs.ConfigRelicsMap[relicsInfo.RelicsId]
				otherEntry := self.MakeOtherEntry(configRelics.OtherGroup, relicsInfo)
				otherType := csvs.ConfigRelicsEntryMap[otherEntry].AttrType
				if otherType != 4 && otherType != 5 {
					isadd = 0
					break
				}
				relicsInfo.OtherEntry = append(relicsInfo.OtherEntry, otherEntry)
			}
		}
		if isadd == 1 {
			bestHeadRelics[i] = relicsInfo
		}
	}
	for _, v := range bestHeadRelics {
		v.ShowInfo()
	}
	fmt.Println("测试", dropCount, "次，获得", len(bestHeadRelics), "次极品头")
}
