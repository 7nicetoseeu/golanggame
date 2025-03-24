package game

import (
	"encoding/json"
	"fmt"
	"golanggame/server/src/csvs"
	"io/ioutil"
	"os"
)

type TaskInfo struct {
	TaskId int
	State  int
}
type ModUniquetask struct {
	MyTaskInfo map[int]*TaskInfo

	player *Player
	path   string
}

func (self *ModUniquetask) SaveData() {
	data, err := json.Marshal(self)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(self.path, data, os.ModePerm)
	if err != nil {
		return
	}
}
func (self *ModUniquetask) InitData() {
	self.MyTaskInfo = make(map[int]*TaskInfo)
}
func (self *ModUniquetask) LoadData(player *Player) {
	self.player = player
	self.path = player.LocalPath + "/uniquetask.json"
	data, err := ioutil.ReadFile(self.path)
	if self.MyTaskInfo == nil {
		self.InitData()
	}
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &self)
	if err != nil {
		return
	}

}

func (self *ModUniquetask) FinishTask(taskId int) bool {
	_, ok := csvs.ConfigUniqueTaskMap[taskId]
	if !ok {
		return false
	}
	self.MyTaskInfo[taskId] = new(TaskInfo)
	self.MyTaskInfo[taskId].State = csvs.TASK_STATE_FINISH
	fmt.Printf("任务Id为%v的任务已完成", taskId)
	return true
}
func (self *ModUniquetask) IsTaskFinish(taskId int) bool {
	_, ok := csvs.ConfigUniqueTaskMap[taskId]
	if !ok {
		return false
	}
	_, ok = self.MyTaskInfo[taskId]
	if !ok {
		fmt.Printf("没有完成id为%d的任务\n", taskId)
		return false
	}
	return true
}
