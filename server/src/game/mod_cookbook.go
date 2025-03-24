package game

type CookBookInfo struct {
	ItemId   int
	ItemName string
}
type ModCookBook struct {
	CookBookInfoMap map[int]*CookBookInfo
}

func (self *ModCookBook) CookBookIsHas(itemId int) bool {
	_, ok := self.CookBookInfoMap[itemId]
	if !ok {
		return false
	}
	return true
}
func (self *ModCookBook) AddCookBook(itemId int) {

}
