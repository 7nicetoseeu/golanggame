package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
)

type IconInfo struct {
	ItemId   int
	ItemName string
}
type ModIcon struct {
	IconInfoMap map[int]*IconInfo

	player *Player
	path   string
}

func (self *ModIcon) LoadData(player *Player) {
	self.player = player
	self.path = player.LocalPath + "/icon.json"
	config, err := ioutil.ReadFile(self.path)
	if self.IconInfoMap == nil {
		self.InitData()
	}
	if err != nil {
		self.InitData()
		return
	}
	err = json.Unmarshal(config, &self)
	if err != nil {
		return
	}
}
func (self *ModIcon) SaveData() {
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, data, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModIcon) InitData() {
	self.IconInfoMap = make(map[int]*IconInfo)
}

func (self *ModIcon) IconIsHas(itemId int) bool {
	_, ok := self.IconInfoMap[itemId]
	if !ok {
		return false
	}
	return true
}
func (self *ModIcon) AddIcon(itemId int) {
	itemConfig := csvs.GetItemConfig(itemId)
	if self.IconIsHas(itemId) {
		fmt.Println("已拥有头像", itemConfig.ItemName)
		return
	}
	if !itemIsExist(itemId) {
		fmt.Println("当前名片不存在")
		return
	}
	iconInfo := &IconInfo{
		ItemId:   itemId,
		ItemName: itemConfig.ItemName,
	}
	self.IconInfoMap[itemId] = iconInfo
	fmt.Println("添加头像", itemConfig.ItemName)
}
func (self *ModIcon) GetIconNow() string {
	iconId := self.player.GetMod(MOD_PLAYER).(*ModPlayer).Icon
	iconConfig := csvs.GetItemConfig(iconId)
	return iconConfig.ItemName
}

func (self *ModIcon) CheckGetIcon(roleId int) {
	iconConfig := csvs.GetConfigIcon(roleId)
	if iconConfig == nil {
		return
	}
	self.AddIcon(iconConfig.IconId)
}
