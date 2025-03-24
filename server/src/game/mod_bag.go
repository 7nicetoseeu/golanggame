package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
)

type itemInfo struct {
	ItemId   int
	ItemNum  int64
	ItemType int64
}

type ModBag struct {
	BagInfoMap map[int]*itemInfo

	player *Player
	path   string
}

func (self *ModBag) SaveData() {
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, data, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModBag) InitData() {
	self.BagInfoMap = make(map[int]*itemInfo)
}
func (self *ModBag) LoadData(player *Player) {
	self.player = player
	self.path = player.LocalPath
	data, err := ioutil.ReadFile(self.path)
	if self.BagInfoMap == nil {
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
func (self *ModBag) AddItem(ItemId int, num int64) {
	itemConfig := csvs.GetItemConfig(ItemId)
	if itemConfig == nil {
		fmt.Println("当前物品不存在", ItemId)
		return
	}
	switch itemConfig.SortType {
	case csvs.ITEMTYPE_NORMAL:
		self.AddItemToBag(ItemId, int64(num))
	case csvs.ITEMTYPE_ROLE:
		self.player.GetMod(MOD_ROLE).(*ModRole).AddRole(ItemId, int(num))
	case csvs.ITEMTYPE_ICON:
		self.player.GetMod(MOD_ICON).(*ModIcon).AddIcon(ItemId)
	case csvs.ITEMTYPE_CARD:
		self.player.GetMod(MOD_CARD).(*ModCard).CheckGetCard(ItemId, int(num))
	case csvs.ITEMTYPE_WEAPON:
		self.player.GetMod(MOD_WEAPON).(*ModWeapon).AddWeapon(ItemId, int(num))
	case csvs.ITEMTYPE_RELICS:
		self.player.GetMod(MOD_RELICS).(*ModRelics).AddRelics(ItemId, int(num))
	case csvs.ITEMTYPE_COOKBOOK:
		self.AddItemToBag(ItemId, 1)
	case csvs.ITEMTYPE_COOK:
		self.player.GetMod(MOD_COOK).(*ModCook).AddCook(ItemId)
	default:
		self.AddItemToBag(ItemId, int64(num))
	}
}
func itemIsExist(ItemId int) bool {
	itemConfig := csvs.GetItemConfig(ItemId)
	if itemConfig == nil {
		fmt.Println("当前物品不存在")
		return false
	}
	return true
}
func (self *ModBag) AddItemToBag(itemId int, num int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println("添加失败，当前物品不存在")
		return
	}
	item := &itemInfo{
		ItemId:   itemId,
		ItemNum:  num,
		ItemType: int64(itemConfig.SortType),
	}
	_, ok := self.BagInfoMap[itemId]
	if ok {
		self.BagInfoMap[itemId].ItemNum += item.ItemNum
	} else {
		self.BagInfoMap[itemId] = item
	}
	fmt.Println(itemConfig.ItemName, "添加", num, "个成功", "该物品数量还有", self.BagInfoMap[itemId].ItemNum)
}

func (self *ModBag) RemoveItem(ItemId int, num int) {
	if num == 0 {
		return
	}
	itemConfig := csvs.GetItemConfig(ItemId)
	if itemConfig == nil {
		fmt.Println("当前物品不存在")
		return
	}
	switch itemConfig.SortType {
	case csvs.ITEMTYPE_NORMAL:
		self.RemoveItemToBag(ItemId, int64(num))
	}
}
func (self *ModBag) RemoveItemToBag(itemId int, num int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println("扣除失败，当前物品不存在")
		return
	}
	item := &itemInfo{
		ItemId:  itemId,
		ItemNum: num,
	}
	_, ok := self.BagInfoMap[itemId]
	if ok {
		self.BagInfoMap[itemId].ItemNum -= item.ItemNum
		if self.BagInfoMap[itemId].ItemNum < 0 {
			self.BagInfoMap[itemId].ItemNum += item.ItemNum
			fmt.Println("扣除失败，当前物品不够，当前物品数量", self.BagInfoMap[itemId].ItemNum)
			return
		}
	} else {
		fmt.Println("背包没有该物品")
		return
	}
	fmt.Println(itemId, "扣除", num, "成功", "该物品数量还有", self.BagInfoMap[itemId].ItemNum)
}
func (self *ModBag) RemoveItemToBagGM(itemId int, num int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println("GM删除失败，当前物品不存在")
		return
	}
	item := &itemInfo{
		ItemId:  itemId,
		ItemNum: num,
	}
	_, ok := self.BagInfoMap[itemId]
	self.BagInfoMap[itemId].ItemNum = 0
	if ok {
		self.BagInfoMap[itemId].ItemNum -= item.ItemNum
	} else {
		fmt.Println("GM背包没有该物品")
		return
	}
	fmt.Println(itemId, "GM删除成功", "该物品数量还有", self.BagInfoMap[itemId].ItemNum)
}

func (self *ModBag) ItemIsBagHas(itemId int, num int) bool {
	if num == 0 {
		return true
	}
	BagitemConfig := self.BagInfoMap[itemId]
	if BagitemConfig == nil {
		fmt.Println("背包中没有", csvs.GetItemName(itemId), itemId)
		return false
	}
	if int(BagitemConfig.ItemNum)-num < 0 {
		fmt.Println("背包中", csvs.GetItemName(itemId), "数量不够", "，还剩下", BagitemConfig.ItemNum, "需要", num)
		return false
	}
	return true
}

func (self *ModBag) UseItem(ItemId int, num int) {
	if !self.ItemIsBagHas(ItemId, num) {
		return
	}
	switch int(self.BagInfoMap[ItemId].ItemType) {
	case csvs.ITEMTYPE_COOKBOOK:
		self.UseCookBook(ItemId)
	case csvs.ITEMTYPE_COOK:
		self.UseCook(ItemId)
	case csvs.ITEMTYPE_FOOD:
		fmt.Println("使用食物")
	}
}
func (self *ModBag) UseCookBook(ItemId int) {
	cookBookConfig := csvs.GetConfigCookBook(ItemId)
	if cookBookConfig == nil {
		fmt.Println("不存在该菜谱")
		return
	}
	self.RemoveItemToBag(ItemId, 1)
	self.AddItem(cookBookConfig.Reward, 1)
}
func (self *ModBag) UseCook(ItemId int) {

}

func (self *ModBag) ShowBag() {
	if len(self.BagInfoMap) == 0 {
		fmt.Println("当前背包为空")
		return
	}
	for _, v := range self.BagInfoMap {
		fmt.Println(csvs.GetItemName(v.ItemId), ":", v.ItemNum, "个")
	}
}
func (self *ModBag) GetItemNum(ItemId int) int {
	_, ok := self.BagInfoMap[ItemId]
	if !ok {
		return 0
	}
	return int(self.BagInfoMap[ItemId].ItemNum)
}
