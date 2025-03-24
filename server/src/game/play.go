package game

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

const (
	MOD_PLAYER     = "player"
	MOD_ICON       = "icon"
	MOD_CARD       = "card"
	MOD_UNIQUETASK = "uniquetask"
	MOD_BAG        = "bag"
	MOD_ROLE       = "role"
	MOD_WEAPON     = "weapon"
	MOD_RELICS     = "relics"
	MOD_COOK       = "cook"
	MOD_POOL       = "pool"
	MOD_MAP        = "map"
)

type ModBase interface {
	LoadData(player *Player)
	SaveData()
	InitData()
}

type Player struct {
	UserId    int64
	modManage map[string]ModBase
	LocalPath string
	Ws        *websocket.Conn
}

func NewTestPlayer(ws *websocket.Conn, userId int64) *Player {
	player := new(Player)

	//************泛型架构*****************
	player.UserId = userId
	player.modManage = map[string]ModBase{
		MOD_PLAYER:     new(ModPlayer),
		MOD_ICON:       new(ModIcon),
		MOD_CARD:       new(ModCard),
		MOD_UNIQUETASK: new(ModUniquetask),
		MOD_BAG:        new(ModBag),
		MOD_ROLE:       new(ModRole),
		MOD_WEAPON:     new(ModWeapon),
		MOD_RELICS:     new(ModRelics),
		MOD_COOK:       new(ModCook),
		MOD_POOL:       new(ModPool),
		MOD_MAP:        new(ModMap),
	}
	player.InitData()
	player.InitMod()
	player.Ws = ws
	return player
}
func (self *Player) GetMod(modName string) ModBase {
	return self.modManage[modName]
}
func (self *Player) InitData() {
	//创建以userid为名字的文件夹
	path := GetServer().ServerConfig.Savepath
	_, err := os.Stat(path)
	if err != nil {
		os.Mkdir(path, os.ModePerm)
	}
	userId := strconv.Itoa(int(self.UserId))
	self.LocalPath = path + "/" + userId
	_, err = os.Stat(self.LocalPath)
	if err != nil {
		err = os.Mkdir(self.LocalPath, os.ModePerm)
		if err != nil {
			return
		}
	}
}
func (self *Player) InitMod() {
	for _, v := range self.modManage {
		v.LoadData(self)
	}
}

func (self *Player) SaveData() {
	for _, v := range self.modManage {
		v.SaveData()
	}
}
func (self *Player) Close() {

}

func (self *Player) RecvSetIcon(IconId int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetIcon(IconId)
	return
}
func (self *Player) RecvSetCard(IconId int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetCard(IconId)
	return
}
func (self *Player) RecvSetName(name string) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetName(name)
	return
}
func (self *Player) RecvSetSign(sign string) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetSign(sign)
	return
}
func (self *Player) ReduceWorldLevel() {
	self.GetMod(MOD_PLAYER).(*ModPlayer).ReduceWorldLevel()
	return
}
func (self *Player) ReturnWorldLevel() {
	self.GetMod(MOD_PLAYER).(*ModPlayer).ReturnWorldLevel()
	return
}
func (self *Player) RecvSetBrith(brith int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetBrith(brith)
	return
}
func (self *Player) RecvSetShowCard(Cards []int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetShowCard(Cards)
	return
}
func (self *Player) RecvSetShowTeam(Teams []int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).SetShowTeam(Teams)
	return
}
func (self *Player) RecvSetHideShowTeam(hide int) {
	self.GetMod(MOD_PLAYER).(*ModPlayer).RecvSetHideShowTeam(hide)
	return
}

func (self *Player) Run() {

	fmt.Println("北华航天工业学院毕业设计——golang游戏服务端")
	fmt.Println("作者：李浩岩")
	if self.GetMod(MOD_PLAYER).(*ModPlayer).IsProhibit == 1 {
		fmt.Println("对不起，你的账号已经被封禁，请联系管理员")
		GetServer().Close()
		return
	}
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓输入你的操作↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	var action string
	for {
		str := `
		0.初始化测试平台（未启用）
		1.基础信息
		2.背包功能
		3.地图
		4.抽奖功能
		5.圣遗物功能
		6.武器功能
		7.关闭服务器(存档)
		`
		fmt.Println(str)
		// self.GetMod(MOD_RELICS).(*ModRelics).TestBestRelics()
		fmt.Scan(&action)
		switch action {
		case "0":
			self.HandlerInitial()
		case "1":
			self.HandleBase()
		case "2":
			self.HandleBag()
		case "3":
			self.HandleMap()
		case "4":
			self.HandleDrop()
		case "5":
			self.HandleRelics()
		case "6":
			self.HandleWeapon()
		case "7":
			self.SaveData()
			GetServer().Close()
		case "esc":
			os.Exit(-1)
		}
	}
}
func (self *Player) HandlerInitial() {
	self.GetMod(MOD_RELICS).(*ModRelics).TestBestRelics()
}
func (self *Player) HandleMap() {
	var action string
	for {
		fmt.Println("地图：")
		str := `
		0.退出
		1.蒙德
		2.璃月
		3.稻妻
		1001.深入风龙废墟
		2001.无妄引咎密宫
		`
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			self.ChangeMap(action)
		case "2":
			self.ChangeMap(action)
		case "3":
			self.ChangeMap(action)
		case "1001":
			self.ChangeMap(action)
		case "2001":
			self.ChangeMap(action)
		case "esc":
			os.Exit(-1)
		}
	}
}
func (self *Player) ChangeMap(mapId string) {
	for {
		Id, _ := strconv.Atoi(mapId)
		self.GetMod(MOD_MAP).(*ModMap).GetMapEvent(Id)
		fmt.Println("\n事件触发")
		var eventId int
		fmt.Println("输入你想触发事件的ID(0返回)")
		fmt.Scan(&eventId)
		if eventId == 0 {
			id, _ := strconv.Atoi(mapId)
			fmt.Printf("id: %v\n", id)
			self.GetMod(MOD_MAP).(*ModMap).RefreshByPlayer(id)
			return
		}
		self.GetMod(MOD_MAP).(*ModMap).SetMapEvent(Id, eventId, 10, self)
	}
}
func (self *Player) HandleDrop() {
	var action string
	for {
		fmt.Println("进入背包")
		str := `
		0.退出
		1.正常抽奖（自定义抽奖次数）
		------概率如下-------------
		5星奖励 60
		4星奖励 510
		武器奖励 9430
		--------------------------
		`
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			self.DropByTimes()
		case "2":
			return
			self.TestDrop()
		case "3":
			return
			self.DropByTimesCheck()
		case "esc":
			os.Exit(-1)
		}
	}
}
func (self *Player) HandleWeapon() {
	var action string
	for {
		fmt.Println("武器界面，测试默认操作一号武器")
		str := `
		0.退出
		1.升级武器
		2.突破武器
		3.穿戴武器
		4.卸下武器
		5.查看装备武器
		6.精炼武器
		`
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			self.GetMod(MOD_WEAPON).(*ModWeapon).UpWeaponLevel(1)
		case "2":
			self.GetMod(MOD_WEAPON).(*ModWeapon).UpWeaponStar(1)
		case "3":
			self.HandleWearWeapon()
		case "4":
			self.HandleTakeOffWeapon()
		case "5":
			self.HandleShowWeapon()
		case "6":
			self.HandleRefineWeapon()
		case "esc":
			os.Exit(-1)
		}
	}
}
func (self *Player) HandleShowWeapon() {
	var roleId int
	fmt.Println("输入目标人物")
	fmt.Scan(&roleId)
	for {
		if roleId == 999 {
			self.GetMod(MOD_ROLE).(*ModRole).ShowAllWeaponInfo(self)
			return
		}
		if self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId] == nil {
			fmt.Println("没有当前人物，请重新输入")
			continue
		}
		roleInfo := self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId]
		roleInfo.ShowWeaponInfo(self)
		break
	}
}
func (self *Player) HandleTakeOffWeapon() {
	var weaponId int
	var roleId int
	roleInfo := new(RoleInfo)
	weaponInfo := new(WeaponInfo)
	fmt.Println("输入目标武器")
	for {
		fmt.Scan(&weaponId)
		if weaponId == 0 {
			return
		}
		if self.GetMod(MOD_WEAPON).(*ModWeapon).ModWeaponMap[weaponId] == nil {
			fmt.Println("背包中没有当前武器，请重新输入")
			continue
		}
		weaponInfo = self.GetMod(MOD_WEAPON).(*ModWeapon).ModWeaponMap[weaponId]
		break
	}
	fmt.Println("输入目标角色")
	for {
		fmt.Scan(&roleId)
		if roleId == 0 {
			return
		}
		if self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId] == nil {
			fmt.Println("没有当前角色，请重新输入")
			continue
		}
		roleInfo = self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId]
		break
	}
	self.GetMod(MOD_ROLE).(*ModRole).TakeOffWeapon(weaponInfo, roleInfo)
}
func (self *Player) HandleWearWeapon() {
	var weaponId int
	var roleId int
	roleInfo := new(RoleInfo)
	weaponInfo := new(WeaponInfo)
	fmt.Println("输入目标武器")
	for {
		fmt.Scan(&weaponId)
		if weaponId == 0 {
			return
		}
		if self.GetMod(MOD_WEAPON).(*ModWeapon).ModWeaponMap[weaponId] == nil {
			fmt.Println("背包中没有当前武器，请重新输入")
			continue
		}
		weaponInfo = self.GetMod(MOD_WEAPON).(*ModWeapon).ModWeaponMap[weaponId]
		break
	}
	fmt.Println("输入目标角色")
	for {
		fmt.Scan(&roleId)
		if roleId == 0 {
			return
		}
		if self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId] == nil {
			fmt.Println("没有当前角色，请重新输入")
			continue
		}
		roleInfo = self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId]
		break
	}
	self.GetMod(MOD_ROLE).(*ModRole).WearWeapon(weaponInfo, roleInfo, self)
}
func (self *Player) HandleRefineWeapon() {
	var weaponId int
	var targetId int
	fmt.Println("输入精炼的武器编号")
	for {
		fmt.Scan(&weaponId)
		if weaponId == 0 {
			return
		}
		if self.GetMod(MOD_WEAPON).(*ModWeapon).ModWeaponMap[weaponId] == nil {
			fmt.Println("无当前武器")
			continue
		}

		break
	}
	fmt.Println("输入消耗的武器编号")
	for {
		fmt.Scan(&targetId)
		if targetId == 0 {
			return
		}
		if self.GetMod(MOD_WEAPON).(*ModWeapon).ModWeaponMap[targetId] == nil {
			fmt.Println("无当前武器")
			continue
		}
		break
	}
	self.GetMod(MOD_WEAPON).(*ModWeapon).RefineWeapon(weaponId, targetId)
}
func (self *Player) HandleRelics() {
	var action string
	for {
		fmt.Println("圣遗物界面，测试默认操作一号圣遗物")
		str := `
		0.退出
		1.升级圣遗物
		2.增加1W经验
		3.模拟一个满级的圣遗物
		4.模拟测试1W次，出现几个极品头
		5.穿戴圣遗物
		6.卸下圣遗物
		7.查看装备圣遗物
		`
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			self.GetMod(MOD_RELICS).(*ModRelics).UpRelicsLevel(self, 1)
		case "2":
			self.HandleAddRelicsExp()
		case "3":
			self.HandleMaxRelics()
		case "4":
			self.HandleTestBestRelics()
		case "5":
			self.HandleWearRelics()
		case "6":
			self.HandleTakeOffRelics()
		case "7":
			self.HandleShowRelics()
		case "8":
			self.RelicsTest()
		case "esc":
			os.Exit(-1)
		}
	}
}
func (self *Player) RelicsTest() {
	three := 0
	four := 0
	for i := 0; i < 10000; i++ {
		i := len(self.GetMod(MOD_RELICS).(*ModRelics).NewRelics(7000001).OtherEntry)
		if i == 3 {
			three++
		} else if i == 4 {
			four++
		}
	}
	fmt.Printf("three: %v\n", three)
	fmt.Printf("four: %v\n", four)
}
func (self *Player) HandleWearRelics() {
	str := "输入穿戴角色编号（0返回）"
	var roleId int
	var roleInfo *RoleInfo
	var relicsInfo *RelicsInfo
	fmt.Println(str)
	for {
		fmt.Scan(&roleId)
		if roleId == 0 {
			return
		}
		roleInfo = self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId]
		if roleInfo != nil {
			break
		}

		fmt.Println("没有当前人物请重新输入（0返回）")
	}
	str = "输入穿戴圣遗物编号（0返回）"
	fmt.Println(str)
	var relicsId int
	for {
		fmt.Scan(&relicsId)
		if relicsId == 0 {
			return
		}
		relicsInfo = self.GetMod(MOD_RELICS).(*ModRelics).ModRelicsMap[relicsId]
		if relicsInfo != nil {
			break
		}
		fmt.Println("没有当前圣遗物请重新输入（0返回）")
	}
	self.GetMod(MOD_ROLE).(*ModRole).WearRelics(relicsInfo, roleInfo, self)
}
func (self *Player) HandleTakeOffRelics() {
	str := "输入目标角色编号（0返回）"
	var roleId int
	var roleInfo *RoleInfo
	var relicsInfo *RelicsInfo
	fmt.Println(str)
	for {
		fmt.Scan(&roleId)
		if roleId == 0 {
			return
		}
		roleInfo = self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId]
		if roleInfo != nil {
			break
		}
		fmt.Println("没有当前人物请重新输入（0返回）")
	}
	str = "输入卸下圣遗物编号（0返回）"
	fmt.Println(str)
	var relicsId int
	for {
		fmt.Scan(&relicsId)
		if relicsId == 0 {
			return
		}
		relicsInfo = self.GetMod(MOD_RELICS).(*ModRelics).ModRelicsMap[relicsId]
		if relicsInfo != nil {
			break
		}
		fmt.Println("没有当前圣遗物请重新输入（0返回）")
	}
	self.GetMod(MOD_ROLE).(*ModRole).TakeOffRelics(relicsInfo, roleInfo, self)
}
func (self *Player) HandleShowRelics() {
	str := "输入查看角色编号（0返回,999查询所有角色圣遗物佩戴情况）"
	var roleId int
	var roleInfo *RoleInfo
	fmt.Println(str)
	for {
		fmt.Scan(&roleId)
		if roleId == 999 {
			self.GetMod(MOD_ROLE).(*ModRole).ShowAllRoleRelics(self)
			return
		}
		if roleId == 0 {
			return
		}
		roleInfo = self.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[roleId]
		if roleInfo != nil {
			break
		}
		fmt.Println("没有当前人物请重新输入（0返回）")
	}
	roleInfo.ShowRelicsInfo(self)
}
func (self *Player) HandleMaxRelics() {
	// self.GetMod(MOD_RELICS).(*ModRelics).AddRelics(7000005, 1)
	self.GetMod(MOD_RELICS).(*ModRelics).UpRelicsLevelTest(self, 1)
}
func (self *Player) HandleTestBestRelics() {
	self.GetMod(MOD_RELICS).(*ModRelics).TestBestRelics()
}
func (self *Player) HandleAddRelicsExp() {
	fmt.Println("获取1w经验")
	self.GetMod(MOD_RELICS).(*ModRelics).AddRelicsExp(self, 1)
}
func (self *Player) TestDrop() {
	fmt.Println("1秒后开始抽奖1W次")
	time.Sleep(1 * time.Second)
	self.GetMod(MOD_POOL).(*ModPool).DoDrop()
}
func (self *Player) DropByTimes() {
	fmt.Println("抽奖功能")
	var times int
	fmt.Println("输入你想抽奖的次数")
	fmt.Scan(&times)
	self.GetMod(MOD_POOL).(*ModPool).DoDropbyTimes(times, self)
}
func (self *Player) DropByTimesCheck() {
	fmt.Println("抽奖功能(带仓检)")
	var times int
	fmt.Println("输入你想抽奖的次数")
	fmt.Scan(&times)
	self.GetMod(MOD_POOL).(*ModPool).DoDropbyTimesCheck(times, self)
}
func (self *Player) HandleBag() {
	var action string
	for {
		fmt.Println("进入背包")
		str := `
		0.退出
		1.增加物品
		2.删除物品
		3.使用物品
		4.查看背包
		5.升级7天神像
		6.神像角色血池测试（1秒恢复10点血量）
		`
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			self.HandleBagAddItem()
		case "2":
			self.HandleBagRemoveItem()
		case "3":
			self.HandleBagUseItem()
		case "4":
			self.ShowBagItem()
		case "5":
			self.HandleUpStatue()
		case "6":
			self.HandleCalHP()
		case "esc":
			os.Exit(-1)
		}
	}
}
func (self *Player) HandleUpStatue() {
	fmt.Println("默认升级风属性七天神像")
	self.GetMod(MOD_MAP).(*ModMap).UpStatue(1, self)
}
func (self *Player) HandleCalHP() {
	fmt.Println("血池测试")
	self.GetMod(MOD_ROLE).(*ModRole).CalHpPool()
}
func (self *Player) ShowBagItem() {
	fmt.Println("普通背包")
	fmt.Println("*********************************")
	self.GetMod(MOD_BAG).(*ModBag).ShowBag()
	fmt.Println("*********************************\n")
	fmt.Println("武器背包")
	fmt.Println("*********************************")
	self.GetMod(MOD_WEAPON).(*ModWeapon).ShowWeaponBag()
	fmt.Println("*********************************\n")
	fmt.Println("圣遗物背包")
	fmt.Println("*********************************")
	self.GetMod(MOD_RELICS).(*ModRelics).ShowRelicsBag()
	fmt.Println("*********************************\n")
	fmt.Println("人物背包")
	fmt.Println("*********************************")
	self.GetMod(MOD_ROLE).(*ModRole).ShowRole()
	fmt.Println("*********************************\n")
}
func (self *Player) HandleBagAddItem() {
	fmt.Println("输入增加物品的Id")
	var itemId int
	fmt.Scan(&itemId)
	if !itemIsExist(itemId) {
		return
	}
	fmt.Println("输入增加物品的数量")
	var itemNum int
	fmt.Scan(&itemNum)
	self.GetMod(MOD_BAG).(*ModBag).AddItem(itemId, int64(itemNum))
}
func (self *Player) HandleBagRemoveItem() {
	fmt.Println("输入删除物品的Id")
	var itemId int
	fmt.Scan(&itemId)
	if !itemIsExist(itemId) {
		return
	}
	fmt.Println("输入删除物品的数量")
	var itemNum int
	fmt.Scan(&itemNum)
	self.GetMod(MOD_BAG).(*ModBag).RemoveItem(itemId, itemNum)
}
func (self *Player) HandleBagUseItem() {
	fmt.Println("输入使用物品的Id")
	var itemId int
	fmt.Scan(&itemId)
	if !itemIsExist(itemId) {
		return
	}
	fmt.Println("输入使用物品的数量")
	var itemNum int
	fmt.Scan(&itemNum)
	self.GetMod(MOD_BAG).(*ModBag).UseItem(itemId, itemNum)
}

func (self *Player) HandleBase() {
	fmt.Println("基础信息界面")
	str := `
	0.返回
	1.查询信息
	2.设置名字
	3.设置签名
	4.头像设置
	5.名片设置
	6.设置生日
	7.人物升级
	`
	var action string
	for {
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			self.SaveData()
			self.showPlayer()
		case "2":
			self.setPlayerName()
		case "3":
			self.setPlayerSign()
		case "4":
			self.PlayerIcon()
		case "5":
			self.PlayerCard()
		case "6":
			self.setPlayerBrith()
		case "7":
			self.PlayerUpLevel()
		case "esc":
			os.Exit(-1)
		}
	}
}

func (self *Player) showPlayer() {
	if self.GetMod(MOD_PLAYER).(*ModPlayer).Name == "" {
		fmt.Println("名字：旅行者")
	} else {
		fmt.Println("名字：", self.GetMod(MOD_PLAYER).(*ModPlayer).Name)
	}
	if self.GetMod(MOD_PLAYER).(*ModPlayer).Sign == "" {
		fmt.Println("签名：未设置")
	} else {
		fmt.Println("签名：", self.GetMod(MOD_PLAYER).(*ModPlayer).Sign)
	}
	if self.GetMod(MOD_PLAYER).(*ModPlayer).Icon == 0 {
		fmt.Println("头像：未设置")
	} else {
		fmt.Println("卡片：", self.GetMod(MOD_ICON).(*ModIcon).GetIconNow())
	}
	if self.GetMod(MOD_PLAYER).(*ModPlayer).Card == 0 {
		fmt.Println("卡片：未设置")
	} else {
		fmt.Println("卡片：", self.GetMod(MOD_CARD).(*ModCard).GetCardNow())
	}
	if self.GetMod(MOD_PLAYER).(*ModPlayer).Birth == 0 {
		fmt.Println("生日：未设置")
	} else {
		month := self.GetMod(MOD_PLAYER).(*ModPlayer).Birth / 100
		day := self.GetMod(MOD_PLAYER).(*ModPlayer).Birth % 100
		fmt.Println("生日：", month, "月", day, "日")
	}
	fmt.Println("大世界等级：", self.GetMod(MOD_PLAYER).(*ModPlayer).WorldLevelNow)
	fmt.Println("人物等级：", self.GetMod(MOD_PLAYER).(*ModPlayer).PlayerLevel)
	fmt.Println("是否为管理员(0否1是)", self.GetMod(MOD_PLAYER).(*ModPlayer).IsGm)
	fmt.Println("是否处于封号状态(0否1是)", self.GetMod(MOD_PLAYER).(*ModPlayer).IsProhibit)
}
func (self *Player) setPlayerName() {
	fmt.Print("输入你的新名字：")
	var name string
	fmt.Scan(&name)
	self.RecvSetName(name)
}
func (self *Player) setPlayerSign() {
	fmt.Print("输入你的新签名：")
	var sign string
	fmt.Scan(&sign)
	self.RecvSetSign(sign)
}
func (self *Player) PlayerIcon() {
	fmt.Println("头像界面")
	str := `
		0.返回
		1.查看头像
		2.修改头像
	`
	var action string
	for {
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			iconInfoMap := self.GetMod(MOD_ICON).(*ModIcon).IconInfoMap
			for _, v := range iconInfoMap {
				fmt.Printf("v: %v\n", v)
			}
		case "2":
			self.PlayerIconAdd()
		}
	}

}
func (self *Player) PlayerIconAdd() {
	fmt.Print("输入你的新头像：")
	var icon int
	fmt.Scan(&icon)
	self.RecvSetIcon(icon)
}
func (self *Player) PlayerCard() {
	fmt.Println("卡片界面")
	str := `
		0.返回
		1.查看名片
		2.修改名片
	`
	var action string
	for {
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			cardInfoMap := self.GetMod(MOD_CARD).(*ModCard).ModCardMap
			for _, v := range cardInfoMap {
				fmt.Printf("v: %v\n", v)
			}
		case "2":
			self.PlayerCardAdd()
		}
	}

}

func (self *Player) PlayerCardAdd() {
	fmt.Print("输入你的新名片：")
	var card int
	fmt.Scan(&card)
	self.RecvSetCard(card)
}
func (self *Player) setPlayerBrith() {
	fmt.Print("输入你的生日：")
	var brith int
	fmt.Scan(&brith)
	self.RecvSetBrith(brith)
}

func (self *Player) PlayerUpLevel() {
	fmt.Println("人物升级")
	str := `
		0.返回
		1.增加5000经验
		2.完成任务(25级)
		3.完成任务(35级)
		4.完成任务(45级)
		5.完成任务(50级)
		6.大世界等级
	`
	var action string
	for {
		fmt.Println(str)
		fmt.Scan(&action)
		switch action {
		case "0":
			return
		case "1":
			self.GetMod(MOD_PLAYER).(*ModPlayer).AddExp(5000)
		case "2":
			self.GetMod(MOD_UNIQUETASK).(*ModUniquetask).FinishTask(10001)
		case "3":
			self.GetMod(MOD_UNIQUETASK).(*ModUniquetask).FinishTask(10002)
		case "4":
			self.GetMod(MOD_UNIQUETASK).(*ModUniquetask).FinishTask(10003)
		case "5":
			self.GetMod(MOD_UNIQUETASK).(*ModUniquetask).FinishTask(10004)
		}
	}

}
