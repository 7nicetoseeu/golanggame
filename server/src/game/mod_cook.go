package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
)

type CookInfo struct {
	ItemId   int
	ItemName string
}
type ModCook struct {
	CookInfoMap map[int]*CookInfo

	player *Player
	path   string
}

func (self *ModCook) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModCook) InitData() {
	self.CookInfoMap = make(map[int]*CookInfo)
}
func (self *ModCook) LoadData(player *Player) {
	self.player = player
	self.path = self.player.LocalPath + "/cook.json"
	data, err := ioutil.ReadFile(self.path)
	if self.CookInfoMap == nil {
		self.CookInfoMap = make(map[int]*CookInfo)
	}
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
func (self *ModCook) CookIsHas(itemId int) bool {
	_, ok := self.CookInfoMap[itemId]
	if !ok {
		return false
	}
	return true
}
func (self *ModCook) AddCook(itemId int) {
	itemConfig := csvs.GetConfigCook(itemId)
	if itemConfig == nil {
		fmt.Println("当前菜谱技能不存在")
		return
	}
	if self.CookIsHas(itemConfig.CookId) {
		fmt.Println("已习得", csvs.GetItemName(itemConfig.CookId))
		return
	}
	CookInfo := &CookInfo{
		ItemId:   itemConfig.CookId,
		ItemName: csvs.GetItemName(itemConfig.CookId),
	}
	self.CookInfoMap[itemConfig.CookId] = CookInfo
	fmt.Println("学习菜谱技能", csvs.GetItemName(itemConfig.CookId))
}
