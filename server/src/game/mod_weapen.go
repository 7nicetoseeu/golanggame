package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
)

type WeaponInfo struct {
	WeaponId    int
	WeaponName  string
	KeyId       int
	Exp         int
	Level       int
	Star        int
	RefineLevel int
	RoleId      int
}

type ModWeapon struct {
	ModWeaponMap map[int]*WeaponInfo
	ModKeyId     int

	player *Player
	path   string
}

func (self *ModWeapon) SaveData() {
	content, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, content, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModWeapon) InitData() {
	self.ModWeaponMap = make(map[int]*WeaponInfo)
}
func (self *ModWeapon) LoadData(player *Player) {
	self.player = player
	self.path = self.player.LocalPath + "/weapon.json"
	configFile, err := ioutil.ReadFile(self.path)
	if self.ModWeaponMap == nil {
		self.ModWeaponMap = make(map[int]*WeaponInfo)
	}
	if err != nil {
		self.InitData()
		return
	}
	err = json.Unmarshal(configFile, &self)
	if err != nil {
		self.InitData()
		return
	}

}
func (self *ModWeapon) AddWeapon(weaponId int, num int) {
	weaponConfig := csvs.GetConfigWeapon(weaponId)
	itemConfig := csvs.GetItemConfig(weaponId)
	if len(self.ModWeaponMap)+num >= csvs.WEAPON_MAX_COUNT {
		fmt.Println("武器超过最大数量")
		return
	}
	for i := 0; i < num; i++ {
		self.ModKeyId++
		weaponInfo := &WeaponInfo{
			WeaponId:   weaponConfig.WeaponId,
			WeaponName: itemConfig.ItemName,
			KeyId:      self.ModKeyId,
			Level:      1,
		}
		fmt.Println("获得武器", weaponInfo.WeaponName, "武器编号", self.ModKeyId, "武器等级", weaponInfo.Level)
		self.ModWeaponMap[self.ModKeyId] = weaponInfo
	}
}
func (self *ModWeapon) ShowWeaponBag() {
	if len(self.ModWeaponMap) == 0 {
		fmt.Println("当前背包为空")
		return
	}
	for _, v := range self.ModWeaponMap {
		fmt.Printf("武器名称%s，当前武器编号%v\n", v.WeaponName, v.KeyId)
	}
}

func (self *ModWeapon) UpWeaponLevel(weaponId int) {
	weaponInfo := self.ModWeaponMap[weaponId]
	if weaponInfo == nil {
		return
	}
	fmt.Println("获得10000经验")
	weaponInfo.Exp += 10000
	WeaponConfig := csvs.ConfigWeaponMap[weaponInfo.WeaponId]
	if WeaponConfig == nil {
		return
	}
	fmt.Println("————升级前————")
	weaponInfo.ShowInfo()
	for {
		WeaponLevelConfig := csvs.WeaponLevelConfigGroup[WeaponConfig.Star][weaponInfo.Level+1]
		if WeaponLevelConfig == nil {
			fmt.Println("已满级")
			weaponInfo.Exp = 0
			return
		}
		if weaponInfo.Exp < WeaponLevelConfig.NeedExp {
			break
		}
		if weaponInfo.Star < WeaponLevelConfig.NeedStarLevel {
			fmt.Println("武器未进行突破")
			weaponInfo.Exp = 0
			return
		}
		weaponInfo.Exp -= WeaponLevelConfig.NeedExp
		weaponInfo.Level++
	}
	fmt.Println("————升级后————")
	weaponInfo.ShowInfo()
}
func (self *ModWeapon) UpWeaponStar(weaponId int) {
	weaponInfo := self.ModWeaponMap[weaponId]
	if weaponInfo == nil {
		return
	}
	weaponConfig := csvs.ConfigWeaponMap[weaponInfo.WeaponId]
	if weaponConfig == nil {
		return
	}
	weaponStar := csvs.WeaponStarConfigGroup[weaponConfig.Star][weaponInfo.Star+1]
	if weaponStar == nil {
		return
	}
	if weaponStar.Level > weaponInfo.Level {
		fmt.Println("武器没有到达突破等级，下一个突破等级为", weaponStar.Level)
		return
	}
	fmt.Println("检测消耗品")
	weaponInfo.Star++
	fmt.Println("武器突破成功，当前星级：", weaponInfo.Star)
}

func (self *ModWeapon) RefineWeapon(weaponId int, targetId int) {
	if weaponId == targetId {
		fmt.Println("不可以对武器本身操作")
		return
	}
	weaponInfo := self.ModWeaponMap[weaponId]
	if weaponInfo == nil {
		return
	}
	targetInfo := self.ModWeaponMap[targetId]
	if targetInfo == nil {
		return
	}
	if weaponInfo.WeaponId != targetInfo.WeaponId {
		fmt.Println("两把武器不一致")
		return
	}
	if weaponInfo.RefineLevel >= csvs.WEAPON_MAX_REFINE {
		fmt.Println("超过最大精炼次数")
		return
	}
	weaponInfo.RefineLevel++
	self.DeleteWeapon(targetId)
	fmt.Println("精炼成功，消耗武器编号", targetId, "武器精炼等级：", weaponInfo.RefineLevel)
}
func (self *ModWeapon) DeleteWeapon(weaponId int) {
	weaponInfo := self.ModWeaponMap[weaponId]
	if weaponInfo == nil {
		return
	}
	self.ModWeaponMap[weaponId] = nil
}

func (self *WeaponInfo) ShowInfo() {
	fmt.Println("武器", self.WeaponName, "武器编号", self.KeyId, "武器等级", self.Level)
}
