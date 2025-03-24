package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
)

type PoolInfo struct {
	PoolId        int
	FiveStarTimes int
	FourStarTimes int
}
type ModPool struct {
	UpPoolInfo *PoolInfo

	player *Player
	path   string
}

func (self *ModPool) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModPool) InitData() {
	// self.CookInfoMap = make(map[int]*CookInfo)
}
func (self *ModPool) LoadData(player *Player) {
	self.player = player
	self.path = self.player.LocalPath + "/pool.json"
	data, err := ioutil.ReadFile(self.path)
	if err != nil {
		self.InitData()
		return
	}
	err = json.Unmarshal(data, &self)
	if err != nil {
		self.InitData()
		return
	}

}
func (self *ModPool) addTimes() {
	self.UpPoolInfo.FiveStarTimes++
	self.UpPoolInfo.FourStarTimes++
}
func resetDropGroup() *csvs.DropGrop {
	dropGrop := new(csvs.DropGrop)
	dropGrop = csvs.DropGropConfigMap[1000]
	return dropGrop
}
func (self *ModPool) DoDropOld() {

	self.UpPoolInfo = new(PoolInfo)
	result := make(map[int]int)
	fiveTimes := make(map[int]int)
	dropGrop := new(csvs.DropGrop)
	if dropGrop == nil {
		return
	}
	for i := 0; i < 1000000; i++ {
		fmt.Printf("第i次抽奖--- %v\n", i)
		//次数加一
		self.addTimes()
		//生成dropgrop,1.超过73次，2.未超过73次默认
		if self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGrop)
			newDropGroup.DropId = dropGrop.DropId
			newDropGroup.WeightAll = dropGrop.WeightAll
			addFiveWight := csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			addFourWight := csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if !(self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT) {
				addFiveWight = 0
			}
			if !(self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT) {
				addFourWight = 0
			}
			for _, v := range dropGrop.ConfigDrop {
				newDropConfig := v
				if v.Result == 10001 {
					//加
					newDropConfig.Weight = v.Weight + addFiveWight
				} else if v.Result == 10003 {
					//减
					newDropConfig.Weight = v.Weight - addFiveWight - addFourWight
				} else {
					newDropConfig.Weight = v.Weight + addFourWight
				}
				newDropGroup.ConfigDrop = append(newDropGroup.ConfigDrop, newDropConfig)
			}
			dropGrop = newDropGroup
		} else {
			dropGrop = resetDropGroup()
		}
		dropConfig := csvs.GetDropGropNew(dropGrop)
		if dropConfig == nil {
			return
		}
		if dropConfig.DropId == csvs.FIVE_STAT_ID1 || dropConfig.DropId == csvs.FIVE_STAT_ID2 {
			if self.UpPoolInfo.FiveStarTimes > 80 {
				fiveTimes[self.UpPoolInfo.FiveStarTimes]++
			}
			self.UpPoolInfo.FiveStarTimes = 0
			self.UpPoolInfo.FourStarTimes = 0
			dropGrop = resetDropGroup()
		}
		if dropConfig.DropId == csvs.FOUR_STAT_ID1 || dropConfig.DropId == csvs.FOUR_STAT_ID2 || dropConfig.DropId == csvs.FOUR_STAT_ID3 {
			self.UpPoolInfo.FourStarTimes = 0
			dropGrop = resetDropGroup()
		}
		result[dropConfig.Result]++
	}
	for index, v := range fiveTimes {
		fmt.Printf("第%v次抽中五星次数：v: %v\n", index, v)
	}
	// for key, v := range result {
	// 	fmt.Println(csvs.GetItemName(key), ":", v, "个")
	// }
}
func (self *ModPool) DoDrop() {
	dropGroup := resetDropGroup()
	self.UpPoolInfo = new(PoolInfo)
	result := make(map[string]int)
	sumCount := make(map[int]int)
	for i := 0; i < 10000; i++ {
		dropGroupConfig := csvs.GetDropGropNew(dropGroup)
		var name string
		if csvs.GetRoleConfig(dropGroupConfig.Result) == nil {
			name = csvs.GetItemName(dropGroupConfig.Result)
		} else {
			name = csvs.GetRoleConfig(dropGroupConfig.Result).ItemName
		}
		result[name]++
		if dropGroupConfig == nil {
			fmt.Println("数据异常")
			return
		}
		if dropGroupConfig.DropId != csvs.FOUR_STAT_ID1 && dropGroupConfig.DropId != csvs.FOUR_STAT_ID2 && dropGroupConfig.DropId != csvs.FOUR_STAT_ID3 {
			self.UpPoolInfo.FourStarTimes++
		}
		if dropGroupConfig.DropId != csvs.FIVE_STAT_ID1 && dropGroupConfig.DropId != csvs.FIVE_STAT_ID2 {
			self.UpPoolInfo.FiveStarTimes++
		}
		if self.UpPoolInfo.FourStarTimes == 10 {
			dropGroupConfig = csvs.GetDropGropNew(csvs.DropGropConfigMap[10002])
		}
		if dropGroupConfig.DropId == csvs.FIVE_STAT_ID1 || dropGroupConfig.DropId == csvs.FIVE_STAT_ID2 {
			sumCount[5]++
			self.UpPoolInfo.FiveStarTimes = 0
		}
		if dropGroupConfig.DropId == csvs.FOUR_STAT_ID1 || dropGroupConfig.DropId == csvs.FOUR_STAT_ID2 || dropGroupConfig.DropId == csvs.FOUR_STAT_ID3 {
			sumCount[4]++
			self.UpPoolInfo.FourStarTimes = 0
		}
		dropGroup = resetDropGroup()
		if self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGrop)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll

			addWight5 := (self.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			addWight4 := (self.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			// addWight4 = 0
			if addWight4 <= 0 {
				addWight4 = 0
			}
			if addWight5 <= 0 {
				addWight5 = 0
			}
			// fmt.Printf("addWight: %v\n", addWight5)
			for _, v := range dropGroup.ConfigDrop {
				newDropGropConfig := v
				if v.Result == 10003 {
					newDropGropConfig.Weight = v.Weight - addWight5 - addWight4
					// fmt.Println("非五，四星权重：", newDropGropConfig.Weight)
				} else if v.Result == 10001 {
					newDropGropConfig.Weight = v.Weight + addWight5
					// fmt.Println("五星权重：", newDropGropConfig.Weight)
				} else if v.Result == 10002 {
					newDropGropConfig.Weight = v.Weight + addWight4
					if newDropGropConfig.Weight > 10000 {
						newDropGropConfig.Weight = 10000
					}
					// fmt.Println("四星权重：", newDropGropConfig.Weight)
				}
				newDropGroup.ConfigDrop = append(newDropGroup.ConfigDrop, newDropGropConfig)
			}
			dropGroup = newDropGroup
		}
	}
	for index, v := range sumCount {
		fmt.Println(index, ":", v)
	}

}
func (self *ModPool) DoDropbyTimes(times int, player *Player) {
	dropGroup := resetDropGroup()
	self.UpPoolInfo = new(PoolInfo)
	result := make(map[string]int, 0)
	sumCount := make(map[int]int)
	for i := 0; i < times; i++ {
		dropGroupConfig := csvs.GetDropGropNew(dropGroup)
		if dropGroupConfig == nil {
			fmt.Println("数据异常")
			return
		}
		result[csvs.GetItemName(dropGroupConfig.Result)]++
		if dropGroupConfig.DropId != csvs.FOUR_STAT_ID1 && dropGroupConfig.DropId != csvs.FOUR_STAT_ID2 && dropGroupConfig.DropId != csvs.FOUR_STAT_ID3 {
			self.UpPoolInfo.FourStarTimes++
		}
		if dropGroupConfig.DropId != csvs.FIVE_STAT_ID1 && dropGroupConfig.DropId != csvs.FIVE_STAT_ID2 {
			self.UpPoolInfo.FiveStarTimes++
		}
		if self.UpPoolInfo.FourStarTimes == 10 {
			dropGroupConfig = csvs.GetDropGropNew(csvs.DropGropConfigMap[10002])
		}
		// player.GetMod(MOD_BAG).(*ModBag).AddItem(dropGroupConfig.Result, 1)
		if dropGroupConfig.DropId == csvs.FIVE_STAT_ID1 || dropGroupConfig.DropId == csvs.FIVE_STAT_ID2 {
			sumCount[5]++
			self.UpPoolInfo.FiveStarTimes = 0
		}
		if dropGroupConfig.DropId == csvs.FOUR_STAT_ID1 || dropGroupConfig.DropId == csvs.FOUR_STAT_ID2 || dropGroupConfig.DropId == csvs.FOUR_STAT_ID3 {
			sumCount[4]++
			self.UpPoolInfo.FourStarTimes = 0
		}
		dropGroup = resetDropGroup()
		// if self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
		// 	newDropGroup := new(csvs.DropGrop)
		// 	newDropGroup.DropId = dropGroup.DropId
		// 	newDropGroup.WeightAll = dropGroup.WeightAll

		// 	addWight5 := (self.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
		// 	addWight4 := (self.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
		// 	if addWight4 <= 0 {
		// 		addWight4 = 0
		// 	}
		// 	if addWight5 <= 0 {
		// 		addWight5 = 0
		// 	}
		// 	for _, v := range dropGroup.ConfigDrop {
		// 		newDropGropConfig := v
		// 		if v.Result == 10003 {
		// 			newDropGropConfig.Weight = v.Weight - addWight5 - addWight4
		// 		} else if v.Result == 10001 {
		// 			newDropGropConfig.Weight = v.Weight + addWight5
		// 		} else if v.Result == 10002 {
		// 			newDropGropConfig.Weight = v.Weight + addWight4
		// 			if newDropGropConfig.Weight > 10000 {
		// 				newDropGropConfig.Weight = 10000
		// 			}
		// 		}
		// 		newDropGroup.ConfigDrop = append(newDropGroup.ConfigDrop, newDropGropConfig)
		// 	}
		// 	dropGroup = newDropGroup
		// }
	}
	// for index, v := range sumCount {
	// 	fmt.Println(index, ":", v)
	// }
	for name, v := range result {
		fmt.Println(name, ":", v)
	}

}
func (self *ModPool) DoDropbyTimesCheck(times int, player *Player) {
	dropGroup := resetDropGroup()
	self.UpPoolInfo = new(PoolInfo)
	sumCount := make(map[int]int)
	result := make(map[string]int, 0)
	for i := 0; i < times; i++ {
		m, m2 := player.GetMod(MOD_ROLE).(*ModRole).GetRoleInfoForPoolCheck()
		dropGroupConfig := csvs.GetDropGropNew1(dropGroup, m, m2)
		if dropGroupConfig == nil {
			fmt.Println("数据异常")
			return
		}
		result[csvs.GetItemName(dropGroupConfig.Result)]++
		if dropGroupConfig.DropId != csvs.FOUR_STAT_ID1 && dropGroupConfig.DropId != csvs.FOUR_STAT_ID2 && dropGroupConfig.DropId != csvs.FOUR_STAT_ID3 {
			self.UpPoolInfo.FourStarTimes++
		}
		if dropGroupConfig.DropId != csvs.FIVE_STAT_ID1 && dropGroupConfig.DropId != csvs.FIVE_STAT_ID2 {
			self.UpPoolInfo.FiveStarTimes++
		}
		player.GetMod(MOD_BAG).(*ModBag).AddItem(dropGroupConfig.Result, 1)
		if dropGroupConfig.DropId == csvs.FIVE_STAT_ID1 || dropGroupConfig.DropId == csvs.FIVE_STAT_ID2 {
			sumCount[5]++
			self.UpPoolInfo.FiveStarTimes = 0
		}
		if dropGroupConfig.DropId == csvs.FOUR_STAT_ID1 || dropGroupConfig.DropId == csvs.FOUR_STAT_ID2 || dropGroupConfig.DropId == csvs.FOUR_STAT_ID3 {
			sumCount[4]++
			self.UpPoolInfo.FourStarTimes = 0
		}
		dropGroup = resetDropGroup()
		if self.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || self.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGrop)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll

			addWight5 := (self.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			addWight4 := (self.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			// addWight4 = 0
			if addWight4 <= 0 {
				addWight4 = 0
			}
			if addWight5 <= 0 {
				addWight5 = 0
			}
			// fmt.Printf("addWight: %v\n", addWight5)
			for _, v := range dropGroup.ConfigDrop {
				newDropGropConfig := v
				if v.Result == 10003 {
					newDropGropConfig.Weight = v.Weight - addWight5 - addWight4
				} else if v.Result == 10001 {
					newDropGropConfig.Weight = v.Weight + addWight5
				} else if v.Result == 10002 {
					newDropGropConfig.Weight = v.Weight + addWight4
					if newDropGropConfig.Weight > 10000 {
						newDropGropConfig.Weight = 10000
					}
				}
				newDropGroup.ConfigDrop = append(newDropGroup.ConfigDrop, newDropGropConfig)
			}
			dropGroup = newDropGroup
		}
	}
	for name, v := range result {
		fmt.Println(name, ":", v)
	}

}
