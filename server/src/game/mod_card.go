package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
)

type CardInfo struct {
	CardId    int
	CardName  string
	FriendNum int
}
type ModCard struct {
	ModCardMap map[int]*CardInfo

	player *Player
	path   string
}

func (self *ModCard) SaveData() {
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, data, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModCard) InitData() {
	self.ModCardMap = make(map[int]*CardInfo)
}
func (self *ModCard) LoadData(player *Player) {
	self.player = player
	self.path = player.LocalPath + "/card.json"
	data, err := ioutil.ReadFile(self.path)
	if self.ModCardMap == nil {
		self.InitData()
	}
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &self)
	if err != nil {
		return
	}
}
func (self *ModCard) CardIsHas(cardId int) bool {
	_, ok := self.ModCardMap[cardId]
	if !ok {
		return false
	}
	return true
}
func (self *ModCard) CheckGetCard(cardId int, friendNum int) {
	cardConfig := csvs.GetConfigCard(cardId)
	itemConfig := csvs.GetItemConfig(cardId)
	if self.CardIsHas(cardId) {
		fmt.Println("已拥有名片", cardConfig.CardId)
		return
	}
	if !itemIsExist(cardId) {
		fmt.Println("当前名片不存在")
		return
	}
	if friendNum < cardConfig.Friendliness {
		fmt.Println("好感度不够")
		return
	}
	cardInfo := &CardInfo{
		CardId:   cardId,
		CardName: itemConfig.ItemName,
	}
	self.ModCardMap[cardId] = cardInfo
	fmt.Println("添加名片", cardInfo.CardName)
}
func (self *ModCard) GetCardNow() string {
	cardId := self.player.GetMod(MOD_PLAYER).(*ModPlayer).Card
	cardConfig := csvs.GetItemConfig(cardId)
	return cardConfig.ItemName
}
func (self *ModCard) AddCard(roleId int) {
	cardCongig := csvs.GetConfigCardByRoleId(roleId)
	if cardCongig == nil {
		return
	}
	self.CheckGetCard(cardCongig.CardId, 10)
}
