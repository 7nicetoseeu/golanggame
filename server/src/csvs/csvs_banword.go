package csvs

import (
	"fmt"
)

type ConfigBanWord struct {
	Id  int
	Txt string
}

var ConfigBanWordSlice []*ConfigBanWord

func init() {
	ConfigBanWordSlice = append(ConfigBanWordSlice,
		&ConfigBanWord{
			Id:  1,
			Txt: "外挂",
		},
		&ConfigBanWord{
			Id:  2,
			Txt: "陪玩",
		},
		&ConfigBanWord{
			Id:  3,
			Txt: "原神",
		},
		&ConfigBanWord{
			Id:  4,
			Txt: "钱",
		},
	)
	fmt.Println("初始化违禁词汇配置")
}

func GetBanWordBase() []string {
	banWord := make([]string, 0)
	for _, v := range ConfigBanWordSlice {
		banWord = append(banWord, v.Txt)
	}
	return banWord
}
