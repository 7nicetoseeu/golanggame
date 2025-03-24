package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
	"time"
)

type RoleInfo struct {
	RoleId     int
	RoleName   string
	GetTimes   int
	Star       int
	WeaponId   int
	RelicsInfo []int //[0,0,1,0,0]
}

type ModRole struct {
	RoleInfoMap map[int]*RoleInfo
	HpPool      int
	HpPoolTime  int64

	player *Player
	path   string
}

func (self *ModRole) SaveData() {
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, data, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModRole) InitData() {
	self.RoleInfoMap = make(map[int]*RoleInfo)
}
func (self *ModRole) LoadData(player *Player) {
	self.player = player
	self.path = player.LocalPath + "/role.json"
	if self.RoleInfoMap == nil {
		self.InitData()
	}
	data, err := ioutil.ReadFile(self.path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &self)
	if err != nil {
		return
	}

}
func (self *ModRole) RoleIsHas(itemId int) bool {
	_, ok := self.RoleInfoMap[itemId]
	if !ok {
		return false
	}
	return true
}
func (self *ModRole) AddRole(itemId int, itemNum int) {
	roleConfig := csvs.GetRoleConfig(itemId)
	if roleConfig == nil {
		fmt.Println("没有当前人物")
		return
	}
	for i := 0; i < itemNum; i++ {
		if self.RoleIsHas(itemId) {
			//拥有人物
			self.RoleInfoMap[itemId].GetTimes++
			if self.RoleInfoMap[itemId].GetTimes >= csvs.ROLE_TIMES_NORMAL_MIN && self.RoleInfoMap[itemId].GetTimes <= csvs.ROLE_TIMES_NORMAL_MAX {
				self.player.GetMod(MOD_BAG).(*ModBag).AddItemToBag(roleConfig.Stuff, roleConfig.StuffNum)
				self.player.GetMod(MOD_BAG).(*ModBag).AddItemToBag(roleConfig.StuffItem, roleConfig.StuffItemNum)
			} else {
				self.player.GetMod(MOD_BAG).(*ModBag).AddItemToBag(roleConfig.MaxStuffItem, roleConfig.MaxStuffItemNum)
			}
		} else {
			//未拥有人物
			fmt.Println("添加人物", roleConfig.ItemName, "人物编号", roleConfig.RoleId)
			self.RoleInfoMap[itemId] = &RoleInfo{
				RoleId:     itemId,
				RoleName:   roleConfig.ItemName,
				GetTimes:   1,
				RelicsInfo: []int{0, 0, 0, 0, 0},
				Star:       roleConfig.Star,
			}
			self.player.GetMod(MOD_ICON).(*ModIcon).CheckGetIcon(itemId)
			self.player.GetMod(MOD_CARD).(*ModCard).AddCard(itemId)
		}
	}
}
func (self *ModRole) ShowRole() {
	if len(self.RoleInfoMap) == 0 {
		fmt.Println("当前背包为空")
		return
	}
	for _, v := range self.RoleInfoMap {
		fmt.Printf("人物名称%s，获得次数%v\n", v.RoleName, v.GetTimes)
	}
}

func (self *ModRole) GetRoleInfoForPoolCheck() (map[int]int, map[int]int) {
	fiveInfo := make(map[int]int, 0)
	fourInfo := make(map[int]int, 0)

	for _, v := range self.RoleInfoMap {
		if v.Star == 5 {
			fiveInfo[v.RoleId] = v.GetTimes
		}
		if v.Star == 4 {
			fourInfo[v.RoleId] = v.GetTimes
		}
	}
	return fiveInfo, fourInfo
}
func (self *ModRole) CalHpPool() {
	if self.HpPoolTime == 0 {
		self.HpPoolTime = time.Now().Unix()
	}
	calTime := time.Now().Unix() - self.HpPoolTime
	self.HpPool += int(calTime) * 10
	self.HpPoolTime = time.Now().Unix()
	fmt.Println("当前血池恢复：", self.HpPool)
}

func (self *ModRole) WearRelics(relicsInfo *RelicsInfo, roleInfo *RoleInfo, player *Player) {
	relicsConfig := csvs.ConfigRelicsMap[relicsInfo.RelicsId]
	if relicsConfig == nil {
		return
	}
	self.CheckRelicsPos(roleInfo, relicsConfig.Pos)
	if relicsInfo.RoleId == roleInfo.RoleId {
		fmt.Println("当前角色已穿戴当前装备，请勿重新穿戴")
		return
	}
	if relicsInfo.RoleId != 0 && roleInfo.RoleId != relicsInfo.RoleId {
		fmt.Println(relicsInfo.RelicsName, "圣遗物已经被其他人物穿戴，开始操作")
		//卸下
		oldRoleInfo := self.RoleInfoMap[relicsInfo.RoleId]
		oldRoleInfo.RelicsInfo[relicsConfig.Pos-1] = csvs.LOGIC_FALSE
		relicsInfo.RoleId = csvs.LOGIC_FALSE
		//装备
		//装备新圣遗物旧人物
		oldRelics := roleInfo.RelicsInfo[relicsConfig.Pos-1]
		oldRoleInfo.RelicsInfo[relicsConfig.Pos-1] = oldRelics
		relicsInfo.RoleId = oldRoleInfo.RoleId
		//装备旧圣遗物新人物
		roleInfo.RelicsInfo[relicsConfig.Pos-1] = relicsInfo.KeyId
		relicsInfo.RoleId = roleInfo.RoleId
		fmt.Println(oldRoleInfo.RoleName, roleInfo.RoleName, "装备已替换")
		// self.WearRelics(relicsInfo, roleInfo, player)
		return
	}
	if roleInfo.RelicsInfo[relicsConfig.Pos-1] != 0 {
		fmt.Println(roleInfo.RoleName, "角色原穿戴圣遗物", relicsInfo.RelicsName, "穿戴位置", relicsConfig.Pos)
		fmt.Println("————————替换————————")
		roleInfo.RelicsInfo[relicsConfig.Pos-1] = csvs.LOGIC_FALSE
		relicsInfo.RoleId = csvs.LOGIC_FALSE
		self.WearRelics(relicsInfo, roleInfo, player)
		return
	}
	roleInfo.RelicsInfo[relicsConfig.Pos-1] = relicsInfo.KeyId
	relicsInfo.RoleId = roleInfo.RoleId
	fmt.Println(roleInfo.RoleName, "角色已穿戴圣遗物", relicsInfo.RelicsName, "穿戴位置", relicsConfig.Pos)
}
func (self *ModRole) TakeOffRelics(relicsInfo *RelicsInfo, roleInfo *RoleInfo, player *Player) {
	relicsConfig := csvs.ConfigRelicsMap[relicsInfo.RelicsId]
	if relicsConfig == nil {
		return
	}
	if roleInfo.RelicsInfo[relicsConfig.Pos-1] == 0 {
		fmt.Println(roleInfo.RoleName, "角色当前没有穿戴圣遗物")
		return
	}
	if roleInfo.RelicsInfo[relicsConfig.Pos-1] != relicsInfo.KeyId {
		fmt.Println("目标角色和目标圣遗物不匹配")
		return
	}
	roleInfo.RelicsInfo[relicsConfig.Pos-1] = csvs.LOGIC_FALSE
	relicsInfo.RoleId = csvs.LOGIC_FALSE
	fmt.Println(roleInfo.RoleName, "角色已卸下圣遗物", relicsInfo.RelicsName, "卸下位置", relicsConfig.Pos)
	// roleInfo.ShowInfo(player)
}
func (self *ModRole) CheckRelicsPos(roleInfo *RoleInfo, pos int) {
	addSize := pos - len(roleInfo.RelicsInfo)
	for i := 0; i < addSize; i++ {
		roleInfo.RelicsInfo = append(roleInfo.RelicsInfo, 0)
	}
}
func (self *RoleInfo) ShowRelicsInfo(player *Player) {
	fmt.Println("当前角色", self.RoleName, "角色ID", self.RoleId)
	suitNum := make(map[int]int, 0)
	for _, relics := range self.RelicsInfo {
		relicsNow := player.GetMod(MOD_RELICS).(*ModRelics).ModRelicsMap[relics]
		if relicsNow == nil {
			fmt.Println("未穿戴")
			continue
		}
		suitId := csvs.ConfigRelicsMap[relicsNow.RelicsId].Type
		suitNum[suitId]++
		fmt.Println("穿戴装备：", relicsNow.RelicsName, relicsNow.KeyId)
	}
	fmt.Println("※※套装效果※※")
	relicsSuit := new(csvs.ConfigRelicsSuit)
	for suitId, num := range suitNum {
		for i := 1; i <= num; i++ {
			if csvs.RelicsSuitConfigMap[suitId][i] == nil {
				continue
			} else {
				relicsSuit = csvs.RelicsSuitConfigMap[suitId][i]
			}
		}
		if relicsSuit != nil {
			if relicsSuit.SkillString == "" {
				relicsSuit.SkillString = "空"
			}
			fmt.Println("套装ID", suitId, "套装技能", relicsSuit.SkillString)
		}
	}
}
func (self *ModRole) ShowAllRoleRelics(player *Player) {
	if len(self.RoleInfoMap) == 0 {
		fmt.Println("当前没有角色")
	}
	for _, RoleInfo := range self.RoleInfoMap {
		RoleInfo.ShowRelicsInfo(player)
	}
}
func (self *ModRole) WearWeapon(weaponInfo *WeaponInfo, roleInfo *RoleInfo, player *Player) {
	if weaponInfo.RoleId == roleInfo.RoleId {
		fmt.Println("请勿重复携带武器")
		return
	}
	//角色已经其他携带武器（卸下）
	if roleInfo.WeaponId != 0 {
		fmt.Println("当前角色已经其他携带武器，武器编号", roleInfo.WeaponId)
		fmt.Println("————开始替换————")
		weaponKeyId := roleInfo.WeaponId
		oldWeaponInfo := player.GetMod(MOD_WEAPON).(*ModWeapon).ModWeaponMap[weaponKeyId]
		if oldWeaponInfo == nil {
			return
		}
		oldWeaponInfo.RoleId = 0
		roleInfo.WeaponId = 0
		self.WearWeapon(weaponInfo, roleInfo, player)
		return
	}
	//武器已经被其他角色携带
	if weaponInfo.RoleId != 0 {
		//卸下武器
		fmt.Println("当前武器已经被其他角色携带，角色编号", weaponInfo.RoleId)
		fmt.Println("————开始替换————")
		oldRoleInfo := player.GetMod(MOD_ROLE).(*ModRole).RoleInfoMap[weaponInfo.RoleId]
		if oldRoleInfo == nil {
			return
		}
		oldRoleInfo.WeaponId = 0
		weaponInfo.RoleId = 0
		self.WearWeapon(weaponInfo, roleInfo, player)
		return
	}
	roleInfo.WeaponId = weaponInfo.KeyId
	weaponInfo.RoleId = roleInfo.RoleId
	fmt.Println(roleInfo.RoleName, "角色穿戴成功，武器", weaponInfo.WeaponName, "编号", weaponInfo.KeyId)

}
func (self *ModRole) TakeOffWeapon(weaponInfo *WeaponInfo, roleInfo *RoleInfo) {
	if roleInfo.WeaponId == 0 {
		fmt.Println("当前人物没有佩戴武器")
		return
	}
	if weaponInfo.KeyId != roleInfo.WeaponId {
		fmt.Println("目标武器和目标人物不匹配")
		return
	}
	weaponInfo.RoleId = csvs.LOGIC_FALSE
	roleInfo.WeaponId = csvs.LOGIC_FALSE
	fmt.Println(roleInfo.RoleName, "已经卸下武器", weaponInfo.WeaponName, "武器编号", weaponInfo.KeyId)
}

func (self *RoleInfo) ShowWeaponInfo(player *Player) {
	fmt.Println("当前角色", self.RoleName, "角色ID", self.RoleId)
	if self.WeaponId == 0 {
		fmt.Println("武器：未穿戴")
		return
	}
	weaponInfo := player.GetMod(MOD_WEAPON).(*ModWeapon).ModWeaponMap[self.WeaponId]
	fmt.Println("武器：", weaponInfo.WeaponName, "编号：", weaponInfo.KeyId, "等级：", weaponInfo.Level)
}

func (self *ModRole) ShowAllWeaponInfo(player *Player) {
	if len(self.RoleInfoMap) == 0 {
		fmt.Println("当前没有人物，请添加人物")
		return
	}
	for _, roleInfo := range self.RoleInfoMap {
		roleInfo.ShowWeaponInfo(player)
	}
}
