package game

import (
	"fmt"
	"golanggame/server/src/csvs"
	"regexp"
	"time"
)

var managerBanWord *ManagerBanWord

type ManagerBanWord struct {
	BanWordBase  []string //配置生成
	BanWordExtra []string //定时器调用
	MsgChannel   chan int
}

func GetManagerBanWord() *ManagerBanWord {
	if managerBanWord == nil {
		managerBanWord = new(ManagerBanWord)
		managerBanWord.BanWordBase = csvs.GetBanWordBase()
		managerBanWord.BanWordExtra = make([]string, 0)
		managerBanWord.MsgChannel = make(chan int)
	}
	return managerBanWord
}

//false不包含违禁词汇，ture包含违禁词汇
func (self *ManagerBanWord) IsBanWord(txt string) bool {
	for _, v := range self.BanWordBase {
		matche, err := regexp.MatchString(v, txt)
		if err != nil {
			fmt.Println("IsBanWord err:" + err.Error())
			return false
		}
		if matche {
			return matche
		}
	}
	for _, v := range self.BanWordExtra {
		matche, err := regexp.MatchString(v, txt)
		if err != nil {
			fmt.Println("IsBanWord err:" + err.Error())
			return false
		}
		if matche {
			return matche
		}
	}
	return false
}

//十秒一次更新词库，并且输出“更新词库”
func (self *ManagerBanWord) Run() {
	GetServer().GoAdd()
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			GetServer().UpBanWords(self.BanWordBase)
			// fmt.Println("更新词库")
		case _, ok := <-self.MsgChannel:
			if !ok {
				// self.BanWordExtra = append(self.BanWordExtra, "添加的词汇")
				GetServer().GoDone()
				return
			}
		}
	}
}

func (self *ManagerBanWord) Close() {
	close(self.MsgChannel)
}
