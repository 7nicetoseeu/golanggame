package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
	"time"
)

type ShowRole struct {
	RoleId    int
	RoleLevel int
}
type ModPlayer struct {
	UserId         int//玩家id
	Icon           int//头像
	Card           int//名片
	Name           string//名字
	Sign           string//签名
	PlayerLevel    int//玩家等级
	PlayerExp      int//玩家经验
	WorldLevel     int//大世界真实等级
	WorldLevelNow  int//大世界当前等级
	WorldLevelCool int//大世界冷却
	Birth          int//生日
	ShowTeam       []*ShowRole
	HideShowTeam   int
	ShowCard       []int
	//看不见的字段
	IsProhibit int//是否为管理员
	IsGm       int//是否封号

	player *Player
	path   string
}

//Recv开头的都是和客户端打交道的函数

func (self *ModPlayer) SaveData() {
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, data, os.ModePerm)
	if err != nil {
		return
	}
	fmt.Println("存档成功")
}
func (self *ModPlayer) LoadData(player *Player) {
	self.player = player
	self.path = player.LocalPath + "/player.json"
	self.UserId = int(player.UserId)
	config, err := ioutil.ReadFile(self.path)
	if err != nil {
		fmt.Println("--创建一个新玩家--from LoadData")
		self.InitData()
		return
	}
	err = json.Unmarshal(config, &self)
	if err != nil {
		return
	}
}

func (self *ModPlayer) InitData() {
	self.PlayerLevel = 1
	self.WorldLevel = 6
	self.WorldLevelNow = 6
}

func (self *ModPlayer) SetIcon(IconId int) {
	if !self.player.GetMod(MOD_ICON).(*ModIcon).IconIsHas(IconId) {
		//通知客户端操作非法
		fmt.Println("未拥有头像", IconId)
		return
	}
	self.Icon = IconId
	fmt.Printf("头像id替换成功：%v\n", IconId)
}
func (self *ModPlayer) SetCard(CardId int) {
	if !self.player.GetMod(MOD_CARD).(*ModCard).CardIsHas(CardId) {
		fmt.Println("未拥有名片", CardId)
		//通知客户端操作非法
		return
	}
	self.Card = CardId
	fmt.Printf("卡片id替换成功：%v\n", CardId)
}

func (self *ModPlayer) SetName(name string) {
	if GetManagerBanWord().IsBanWord(name) {
		fmt.Printf("名字有违禁词汇：%v\n", name)
		return
	}
	self.Name = name
	fmt.Printf("名字替换成功：%v\n", name)
}

func (self *ModPlayer) SetSign(sign string) {
	if GetManagerBanWord().IsBanWord(sign) {
		fmt.Printf("签名有违禁词汇：%v\n", sign)
		return
	}
	self.Sign = sign
	fmt.Printf("签名替换成功：%v\n", sign)
}
func (self *ModPlayer) AddExp(exp int) {
	self.PlayerExp += exp
	for {
		config := csvs.GetNowLevelConfig(self.PlayerLevel)
		if self.PlayerLevel >= len(csvs.ConfigPlayerLevelSlice) {
			self.PlayerExp = 0
			break
		}
		if config == nil || config.PlayerExp == 0 {
			break
		}
		if self.PlayerExp < config.PlayerExp {
			break
		}
		//判读是否完成任务
		if config.ChapterId != 0 && !self.player.GetMod(MOD_UNIQUETASK).(*ModUniquetask).IsTaskFinish(config.ChapterId) {
			break
		}
		self.PlayerExp -= config.PlayerExp
		self.PlayerLevel = config.PlayerLevel + 1
	}
	fmt.Printf("角色升级已经升级到：%v---当前经验值：%v\n", self.PlayerLevel, self.PlayerExp)
}
func (self *ModPlayer) ReduceWorldLevel() {
	if self.WorldLevel < csvs.REDUCE_WORLD_LEVEL_START {
		fmt.Printf("操作失败，当前世界等级小于最小限度，当前世界等级%v，真实世界等级%v\n", self.WorldLevelNow, self.WorldLevel)
		return
	}
	if time.Now().Unix() < int64(self.WorldLevelCool) {
		fmt.Printf("操作失败，处于冷却时间，当前世界等级%v，真实世界等级%v\n", self.WorldLevelNow, self.WorldLevel)
		return
	}
	if self.WorldLevel-self.WorldLevelNow >= csvs.REDUCE_WORLD_LEVEL_MAX {
		fmt.Printf("操作失败，当前世界等级%v，真实世界等级%v\n", self.WorldLevelNow, self.WorldLevel)
		return
	}

	self.WorldLevelNow--
	self.WorldLevelCool = int(time.Now().Unix()) + csvs.REDUCE_WORLD_LEVEL_TIME
	fmt.Printf("操作成功，当前世界等级%v，真实世界等级%v\n", self.WorldLevelNow, self.WorldLevel)
}

func (self *ModPlayer) ReturnWorldLevel() {
	if self.WorldLevel == self.WorldLevelNow {
		fmt.Printf("操作失败，当前世界等级等于真实世界等级，当前世界等级%v，真实世界等级%v\n", self.WorldLevelNow, self.WorldLevel)
		return
	}
	self.WorldLevelNow = self.WorldLevel
	fmt.Printf("返回世界等级操作成功，当前世界等级%v，真实世界等级%v\n", self.WorldLevelNow, self.WorldLevel)
}

func (self *ModPlayer) SetBrith(brith int) {
	if self.IsGm == 0 {
		if self.Birth != 0 {
			fmt.Println("生日只设置一次")
			return
		}
	}
	month := brith / 100
	day := brith % 100
	if month == 0 {
		month = brith / 10
		day = brith % 10
		brith = month*100 + day
	}
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		if day < 0 || day > 31 {
			fmt.Println(month, "月没有", day, "号")
			return
		}
	case 4, 6, 9, 11:
		if day < 0 || day > 30 {
			fmt.Println(month, "月没有", day, "号")
			return
		}
	case 2:
		if day < 0 || day > 28 {
			fmt.Println(month, "月没有", day, "号")
			return
		}
	default:
		fmt.Println("没有", month, "月")
		return
	}
	fmt.Println("生日设置成功，你的生日是", month, "月", day, "日")
	self.Birth = brith
	if IsBrithDay(month, day) {
		fmt.Println("今天是你的生日，祝你生日快乐")
	} else {
		fmt.Println("期待你生日的到来")
	}
}
func IsBrithDay(month, day int) bool {
	nowDay := time.Now().Day()
	nowMonth := int(time.Now().Month())
	if nowMonth == month && nowDay == day {
		return true
	}
	return false
}

func (self *ModPlayer) SetShowCard(Cards []int) {
	if len(Cards) > csvs.SHOWCARD_MAX {
		return
	}
	existCard := make(map[int]int)
	showCard := make([]int, 0)
	for _, cardId := range Cards {
		_, ok := existCard[cardId]
		if ok {
			continue
		}
		if !self.player.GetMod(MOD_CARD).(*ModCard).CardIsHas(cardId) {
			//通知客户端操作非法
			continue
		}
		existCard[cardId] = 1
		showCard = append(showCard, cardId)
	}
	self.ShowCard = showCard
	fmt.Println("设置展示卡片成功：", showCard)
}

func (self *ModPlayer) SetShowTeam(RoleIds []int) {
	if len(RoleIds) > csvs.SHOWTEAM_MAX {
		return
	}
	existRole := make(map[int]int)
	showRoles := make([]*ShowRole, 0)
	for _, RoldId := range RoleIds {
		_, ok := existRole[RoldId]
		if ok {
			continue
		}
		if !self.player.GetMod(MOD_CARD).(*ModCard).CardIsHas(RoldId) {
			//通知客户端操作非法
			continue
		}
		existRole[RoldId] = 1
		showRole := new(ShowRole)
		showRole.RoleId = RoldId
		showRole.RoleLevel = 80
		showRoles = append(showRoles, showRole)
	}
	self.ShowTeam = showRoles
	fmt.Println("设置展示阵容成功：")
	for _, v := range showRoles {
		fmt.Printf("v: %v\n", v)
	}
}
func (self *ModPlayer) RecvSetHideShowTeam(hide int) {
	if hide == csvs.SHOWTEAM_HIDE || hide == csvs.SHOWTEAM_NOHIDE {
		return
	}
	self.HideShowTeam = hide
}
func (self *ModPlayer) GetPlayerWorldLevel() int {
	return self.WorldLevelNow
}
func (self *ModPlayer) RecvUpRelicsLevel(relicsKeyId int) {
	self.player.GetMod(MOD_RELICS).(*ModRelics).UpRelicsLevel(self.player, relicsKeyId)
}
