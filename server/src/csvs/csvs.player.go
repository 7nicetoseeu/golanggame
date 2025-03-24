package csvs

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type ConfigPlayerLevel struct {
	PlayerLevel int `json:"PlayerLevel"`
	PlayerExp   int `json:"PlayerExp"`
	WorldLevel  int `json:"WorldLevel"`
	ChapterId   int `json:"ChapterId"`
}

var ConfigPlayerLevelSlice []*ConfigPlayerLevel

func init() {
	filepath := "../../csv/PlayerLevel.csv"
	fs, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	r := csv.NewReader(fs)
	//针对大文件，一行一行的读取文件
	for {
		rows, err := r.Read()
		if err != nil && err != io.EOF {
			log.Fatalf("can not read, err is %+v", err)
		}
		if err == io.EOF {
			break
		}
		// fmt.Println(rows)
		configPlayerLevel := new(ConfigPlayerLevel)
		for index, v := range rows {
			num, _ := strconv.Atoi(v)
			if index == 0 {
				configPlayerLevel.PlayerLevel = num
			}
			if index == 1 {
				configPlayerLevel.PlayerExp = num
			}
			if index == 2 {
				configPlayerLevel.WorldLevel = num
			}
			if index == 3 {
				configPlayerLevel.ChapterId = num
			}
		}
		ConfigPlayerLevelSlice = append(ConfigPlayerLevelSlice, configPlayerLevel)
	}
	ConfigPlayerLevelSlice[0].PlayerLevel = 1
	// for _, v := range ConfigPlayerLevelSlice {
	// 	fmt.Printf("v: %v\n", v)
	// }
}

func GetNowLevelConfig(level int) *ConfigPlayerLevel {
	if level < 0 || level > len(ConfigPlayerLevelSlice) {
		return nil
	}
	return ConfigPlayerLevelSlice[level]
}
