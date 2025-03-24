package csvs

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type ConfigUniqueTask struct {
	TaskId    int
	SortType  int
	OpenLevel int
	TaskType  int
	Condition int
}

var ConfigUniqueTaskMap map[int]*ConfigUniqueTask

func init() {
	ConfigUniqueTaskMap = make(map[int]*ConfigUniqueTask)
	filepath := "../../csv/UniqueTask.csv"
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
		configUniqueTask := new(ConfigUniqueTask)
		for index, v := range rows {
			num, _ := strconv.Atoi(v)
			if index == 0 {
				configUniqueTask.TaskId = num
			}
			if index == 1 {
				configUniqueTask.SortType = num
			}
			if index == 2 {
				configUniqueTask.OpenLevel = num
			}
			if index == 3 {
				configUniqueTask.TaskType = num
			}
			if index == 4 {
				configUniqueTask.Condition = num
			}
		}
		ConfigUniqueTaskMap[configUniqueTask.TaskId] = configUniqueTask
	}
	ConfigUniqueTaskMap[0].TaskId = 10001
}
