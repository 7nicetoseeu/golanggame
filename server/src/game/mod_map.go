package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type MapInfo struct {
	MapId     int
	EventInfo map[int]*Event
}
type Event struct {
	EventId         int
	State           int
	NextRefreshTime int64
}
type StatueInfo struct {
	StatueId int
	Level    int
	ItemInfo map[int]itemInfo
}
type ModMap struct {
	MapInfo    map[int]*MapInfo
	StatueInfo map[int]*StatueInfo

	player *Player
	path   string
}

//改变地点
//触发事件：打怪，物品互动
func (self *ModMap) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModMap) InitData() {
	self.MapInitData()
}
func (self *ModMap) LoadData(player *Player) {
	self.player = player
	self.path = self.player.LocalPath + "/cook.json"
	data, err := ioutil.ReadFile(self.path)
	if self.MapInfo == nil {
		self.InitData()
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
func (self *ModMap) MapInitData() {
	self.MapInfo = make(map[int]*MapInfo)
	self.StatueInfo = make(map[int]*StatueInfo)
	configMapMap := csvs.ConfigMapMap
	configEventMap := csvs.ConfigEventMap

	for _, v := range configMapMap {
		_, ok := self.MapInfo[v.MapId]
		if !ok {
			self.MapInfo[v.MapId] = NewMapInfo(v.MapId)
		}
	}
	for _, event := range configEventMap {
		_, ok := self.MapInfo[event.MapId]
		if !ok {
			continue
		}
		_, ok = self.MapInfo[event.MapId].EventInfo[event.EventId]
		if !ok {
			self.MapInfo[event.MapId].EventInfo[event.EventId] = new(Event)
			self.MapInfo[event.MapId].EventInfo[event.EventId].State = csvs.MAP_ITEM_STATE_START
			self.MapInfo[event.MapId].EventInfo[event.EventId].EventId = event.EventId
		}
	}
}
func NewMapInfo(mapId int) *MapInfo {
	mapInfo := new(MapInfo)
	mapInfo.MapId = mapId
	mapInfo.EventInfo = make(map[int]*Event)
	return mapInfo
}
func (self *ModMap) GetMapEvent(mapId int) {
	mapInfo, ok := self.MapInfo[mapId]
	if !ok {
		fmt.Println("地图不存在")
		return
	}
	fmt.Println("当前处于", csvs.GetConfigMap(mapId).MapName)
	fmt.Println("当前时间：", time.Now().Unix())
	for _, v := range mapInfo.EventInfo {
		self.CheckRefresh(v)
		lastTime := v.NextRefreshTime - time.Now().Unix()
		if lastTime <= 0 {
			lastTime = 0
		}
		fmt.Println("事件ID", v.EventId, "事件名称:", csvs.GetConfigEvent(v.EventId).Name, "事件状态", v.State, "刷新时间：", v.NextRefreshTime, "还剩下", lastTime, "秒")
	}
}

func (self *ModMap) CheckRefresh(event *Event) {
	if event.NextRefreshTime >= time.Now().Unix() {
		return
	}
	eventConfig := csvs.GetConfigEvent(event.EventId)
	if eventConfig == nil {
		return
	}
	switch eventConfig.RefreshType {
	case csvs.MAP_ITEM_REFRESH_DAY:
		count := time.Now().Unix()/csvs.MAP_ITEM_REFRESH_DAY_TIME + 1
		event.NextRefreshTime = count * csvs.MAP_ITEM_REFRESH_DAY_TIME
	case csvs.MAP_ITEM_REFRESH_WEEK:
		count := time.Now().Unix()/csvs.MAP_ITEM_REFRESH_WEEK_TIME + 1
		event.NextRefreshTime = count * csvs.MAP_ITEM_REFRESH_WEEK_TIME
	case csvs.MAP_ITEM_REFRESH_JION:
		return
	}
	event.State = csvs.MAP_ITEM_STATE_START
}
func (self *ModMap) RefreshByPlayer(mapId int) {
	if csvs.GetConfigMap(mapId).MapType != csvs.MAP_TYPE_PLAYER {
		return
	}
	for _, eventConfig := range self.MapInfo[mapId].EventInfo {
		eventConfig.State = 0
	}
}
func (self *ModMap) SetMapEvent(mapId int, eventId int, state int, player *Player) {
	_, ok := self.MapInfo[mapId]
	if !ok {
		fmt.Println("地图不存在")
		return
	}
	_, ok = self.MapInfo[mapId].EventInfo[eventId]
	if !ok {
		fmt.Println("事件不存在")
		return
	}
	if self.MapInfo[mapId].EventInfo[eventId].State >= state {
		fmt.Println("参数异常")
		return
	}
	eventConfig := csvs.GetConfigEvent(eventId)
	if eventConfig == nil {
		return
	}
	mapConfig := csvs.GetConfigMap(mapId)
	if mapConfig == nil {
		return
	}

	if !player.GetMod(MOD_BAG).(*ModBag).ItemIsBagHas(eventConfig.CostItem, eventConfig.CostNum) {
		return
	}
	player.GetMod(MOD_BAG).(*ModBag).RemoveItem(eventConfig.CostItem, eventConfig.CostNum)
	if mapConfig.MapType == csvs.MAP_TYPE_PLAYER && eventConfig.EventType == csvs.EVENT_ITEM_TYPE_REWARD {
		for _, v := range self.MapInfo[mapId].EventInfo {
			if v.EventId == eventId {
				continue
			}
			if csvs.GetConfigEvent(v.EventId).EventType == csvs.EVENT_ITEM_TYPE_REWARD {
				continue
			}
			if v.State != csvs.MAP_ITEM_STATE_END {
				fmt.Println("有事件尚未完成", v.EventId)
				return
			}
		}
	}
	self.MapInfo[mapId].EventInfo[eventId].State = state
	if self.MapInfo[mapId].EventInfo[eventId].State == csvs.MAP_ITEM_STATE_FINISH {
		fmt.Println(csvs.GetConfigEvent(eventId).Name, "事件完成")
	}
	if self.MapInfo[mapId].EventInfo[eventId].State == csvs.MAP_ITEM_STATE_END {
		for i := 0; i < eventConfig.EventDropTimes; i++ {
			fmt.Println("第", i+1, "次")
			dropGropItem := csvs.GetDropGropItemNew(eventConfig.EventDrop)
			if dropGropItem == nil {
				break
			}
			for _, configDropItem := range dropGropItem {
				randNum := rand.Intn(csvs.DORP_ITEM_ALL_PERCENT)
				if randNum < configDropItem.Weight {
					randAll := configDropItem.ItemNumMax - configDropItem.ItemNumMin + 1
					randNum1 := rand.Intn(randAll) + configDropItem.ItemNumMin
					worldLevel := player.GetMod(MOD_PLAYER).(*ModPlayer).GetPlayerWorldLevel()
					if worldLevel > 0 {
						randNum1 = randNum1 * (csvs.DORP_ITEM_ALL_PERCENT + worldLevel*configDropItem.WorldAdd) / csvs.DORP_ITEM_ALL_PERCENT
					}
					player.GetMod(MOD_BAG).(*ModBag).AddItem(configDropItem.ItemId, int64(randNum1))
				}
			}
		}
		fmt.Println(csvs.GetConfigEvent(eventId).Name, "事件领取")
	}
	if state > 0 {
		eventConfig := csvs.GetConfigEvent(eventId)
		if eventConfig == nil {
			return
		}
		switch eventConfig.RefreshType {
		case csvs.MAP_ITEM_REFRESH_SELF:
			self.MapInfo[mapId].EventInfo[eventId].NextRefreshTime = time.Now().Unix() + csvs.MAP_ITEM_REFRESH_SELF_TIME
		}
		return
	}
}

func (self *ModMap) RefreshDay() {
	for _, v1 := range self.MapInfo {
		for _, v := range v1.EventInfo {
			if csvs.GetConfigEvent(v.EventId) == nil {
				continue
			}
			if csvs.GetConfigEvent(v.EventId).RefreshType == csvs.MAP_ITEM_REFRESH_DAY {
				v.State = csvs.MAP_ITEM_STATE_START
				v.NextRefreshTime = time.Now().Unix() + csvs.MAP_ITEM_REFRESH_DAY_TIME
				// fmt.Println("日刷新")
			}
		}
	}
}
func (self *ModMap) RefreshWeek() {
	for _, v1 := range self.MapInfo {
		for _, v := range v1.EventInfo {
			if csvs.GetConfigEvent(v.EventId) == nil {
				continue
			}
			if csvs.GetConfigEvent(v.EventId).RefreshType == csvs.MAP_ITEM_REFRESH_WEEK {
				self.MapInfo[v1.MapId].EventInfo[v.EventId].State = csvs.MAP_ITEM_STATE_START
			}
		}
	}
}
func (self *ModMap) Refreshself() {
	for _, v1 := range self.MapInfo {
		for _, v := range v1.EventInfo {
			if csvs.GetConfigEvent(v.EventId) == nil {
				continue
			}
			if csvs.GetConfigEvent(v.EventId).RefreshType == csvs.MAP_ITEM_REFRESH_SELF {
				if v.NextRefreshTime <= time.Now().Unix() {
					v.State = csvs.MAP_ITEM_STATE_START
				}
			}
		}
	}
}
func NewStatuesInfo(StatueId int) *StatueInfo {
	statueInfo := new(StatueInfo)
	statueInfo.StatueId = StatueId
	statueInfo.ItemInfo = make(map[int]itemInfo)
	return statueInfo
}
func (self *ModMap) UpStatue(StatueId int, player *Player) {
	_, ok := self.StatueInfo[StatueId]
	if !ok {
		self.StatueInfo[StatueId] = NewStatuesInfo(StatueId)
	}
	statueInfo := self.StatueInfo[StatueId]
	levelNext := statueInfo.Level + 1
	ConfigStatue, ok := csvs.ConfigStatueMap[StatueId][levelNext]
	if !ok {
		return
	}

	costNum := ConfigStatue.CostNum
	itemId := ConfigStatue.CostItem
	itemName := csvs.GetItemName(itemId)
	if !player.GetMod(MOD_BAG).(*ModBag).ItemIsBagHas(itemId, costNum) {
		//不够升级
		HasNum := player.GetMod(MOD_BAG).(*ModBag).GetItemNum(itemId)
		player.GetMod(MOD_BAG).(*ModBag).RemoveItem(itemId, HasNum)
		ConfigStatue.CostNum -= HasNum
		fmt.Println("使用数量:", HasNum, "，神像当前等级:", levelNext, ",还差", ConfigStatue.CostNum, itemName)
	} else {
		//够升级
		self.StatueInfo[StatueId].Level++
		player.GetMod(MOD_BAG).(*ModBag).RemoveItem(itemId, costNum)
		fmt.Println("升级成功，神像当前等级：", levelNext)
	}
}
