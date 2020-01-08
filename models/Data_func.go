package models

import (
	"../common"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Luxurioust/excelize"
)

var logs = &common.Loger

func (data *Data) LoadAwards() (*Data, error) {
	fileName := "awards.json"
	dataJsonByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		logs.Error("open award file error:%v", err)
		return data, err
	}
	dataJsonStr := string(dataJsonByte)
	err = json.Unmarshal([]byte(dataJsonStr), &data.Awards)
	if err != nil {
		logs.Error("Unmarshal award file error:%v", err)
		return data, err
	}
	//logs.Info("Awards:%v", data.Awards)
	return data, nil
}

func (data *Data) LoadDrawedRecord() (*Data, error) {
	fileName := "drawed_record.json"
	dataJsonByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		logs.Error("open drawed_record file error:%v", err)
		return data, err
	}
	dataJsonStr := string(dataJsonByte)
	err = json.Unmarshal([]byte(dataJsonStr), &data.DrawedRecords)
	if err != nil {
		logs.Error("Unmarshal drawed_record file error:%v", err)
		return data, err
	}
	return data, nil
}

func (data *Data) LoadBackMoneyRecord() (*Data, error) {
	fileName := "back_money_record.json"
	dataJsonByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		logs.Error("open back_money_record file error:%v", err)
		return data, err
	}
	dataJsonStr := string(dataJsonByte)
	err = json.Unmarshal([]byte(dataJsonStr), &data.BackMoneyRecords)
	if err != nil {
		logs.Error("Unmarshal back_money_record file error:%v", err)
		return data, err
	}
	logs.Info("back_money_record:%v", data.BackMoneyRecords)
	return data, nil
}
func (data *Data) LoadAction() (*Data, error) {
	fileName := "action.json"
	dataJsonByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		logs.Error("open action file error:%v", err)
		return data, err
	}
	dataJsonStr := string(dataJsonByte)
	err = json.Unmarshal([]byte(dataJsonStr), &data.Actions)
	if err != nil {
		logs.Error("Unmarshal action file error:%v", err)
		return data, err
	}
	for index, action := range data.Actions {
		flag := false

		for _, award := range data.Awards {
			if award.ID == action.AwardID {
				flag = true
				data.Actions[index].AwardName = award.Name
				break

			}

		}

		if !flag {
			fmt.Println("action not match award")
			logs.Error("action:%v error,no match award:%v", action.ID, action.AwardID)
			return data, errors.New("action not match award")

		}

	}

	//logs.Info("action:%v", data.Actions)
	return data, nil
}
func (data *Data) LoadExistUser() (*Data, error) {
	fileName := "users.json"
	dataJsonByte, err := ioutil.ReadFile(fileName)
	if err != nil {
		logs.Error("open action file error:%v", err)
		return data, err
	}
	dataJsonStr := string(dataJsonByte)
	err = json.Unmarshal([]byte(dataJsonStr), &data.Users)
	if err != nil {
		logs.Error("Unmarshal action file error:%v", err)
		return data, err
	}
	//logs.Info("action:%v", data.Actions)
	return data, nil
}
func (data *Data) LoadUser() (*Data, error) {

	fileName := "fit_2019_2.xlsx"
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		logs.Error("open file error:%s", err)
		return data, err
	}
	rows, err := xlsx.GetRows("员工明细表")
	if err != nil {
		logs.Error("open file error:%s", err)
		return data, err
	}
	//logs.Info("row:%d", len(rows))
	for _, row := range rows {
		id, err := strconv.Atoi(row[0])
		if err != nil {
			logs.Info(err)
			continue

		}
		name := row[1]
		fullName := row[2]
		level, err := strconv.Atoi(row[4])
		if err != nil {
			logs.Info(err)
			continue

		}
		var user User
		user.ID = id
		user.Name = name
		user.Level = level
		user.IsDrawed = false
		user.FullName = fullName
		data.Users = append(data.Users, user)

	}

	rand.Seed(time.Now().Unix())

	sort.Slice(data.Users, func(i, j int) bool {
		tt := rand.Intn(2)
		//logs.Info("rand:%v", tt)
		if (tt) > 0 {
			return true
		}
		return false
	})
	//logs.Info("users:%v", data.Users)
	return data, nil

}

func (data *Data) SetData() error {

	go data.SetDataCore()

	return nil
}

var (
	mutex sync.Mutex
)

//水电费水电费
//非线程安全胜多负少你好  水电费第三方水电费
func (data *Data) SetDataCore() error {
	mutex.Lock()
	defer mutex.Unlock()

	userJsonByte, err := json.Marshal(data.Users)
	if err != nil {
		logs.Error("error:%v", err)
		return err
	}
	err = ioutil.WriteFile("users.json", userJsonByte, 0644)
	if err != nil {
		logs.Error("error:%v", err)
		return err
	}

	actionJsonByte, err := json.Marshal(data.Actions)
	if err != nil {
		logs.Error("error:%v", err)
		return err
	}
	err = ioutil.WriteFile("action.json", actionJsonByte, 0644)
	if err != nil {
		logs.Error("error:%v", err)
		return err
	}

	drawedRecordJsonByte, err := json.Marshal(data.DrawedRecords)
	err = ioutil.WriteFile("drawed_record.json", drawedRecordJsonByte, 0644)
	if err != nil {
		logs.Error("error:%v", err)
		return err
	}

	backMoneyRecordJsonByte, err := json.Marshal(data.BackMoneyRecords)
	err = ioutil.WriteFile("back_money_record.json", backMoneyRecordJsonByte, 0644)
	if err != nil {
		logs.Error("error:%v", err)
		return err
	}

	logs.Info("save data finish")
	return nil
}

func (data *Data) ReSetData() error {
	data, err := data.CountData()
	if err != nil {
		return err
	}
	dataJsonByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	dataJsonStr := string(dataJsonByte)
	err = ioutil.WriteFile("data.json", []byte(dataJsonStr), 0644)
	if err != nil {
		return err
	}
	return nil
}

var fileLock = false

func (data *Data) CountData() (*Data, error) {
	//PoolMoney
	data.Count.PoolMoney = 0
	for _, item := range data.BackMoneyRecords {
		data.Count.PoolMoney += item.Money
	}
	for _, item := range data.DrawedRecords {
		if item.AwardID == 0 {
			data.Count.PoolMoney -= item.AwardCount
		}
	}
	if data.Count.PoolMoney < 0 {
		return data, errors.New("pool money less than 0")
	}
	//AllPeopleCount
	data.Count.AllPeopleCount = len(data.Users)

	data.Count.AllLuckyCount = 0
	for _, item := range data.Users {
		if item.IsDrawed {
			data.Count.AllLuckyCount += 1
		}
	}
	//LuckyLeaderCount
	data.Count.LuckyLeaderCount = 0
	for _, item := range data.Users {
		if item.IsDrawed && item.Level == 2 {
			data.Count.LuckyLeaderCount += 1
		}
	}

	//NoLuckyLeaderCount
	data.Count.NoLuckyLeaderCount = 0
	for _, item := range data.Users {
		if item.IsDrawed == false && item.Level == 2 {
			data.Count.NoLuckyLeaderCount += 1
		}
	}

	//LuckyStaffCount
	data.Count.LuckyStaffCount = 0
	for _, item := range data.Users {
		if item.IsDrawed == true && item.Level == 1 {
			data.Count.LuckyStaffCount += 1
		}
	}

	//NoLuckyStaffCount
	data.Count.NoLuckyStaffCount = 0
	for _, item := range data.Users {
		if item.IsDrawed == false && item.Level == 1 {
			data.Count.NoLuckyStaffCount += 1
		}
	}
	return data, nil
}
