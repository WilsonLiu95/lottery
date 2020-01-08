package api

import (
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"../models"

	"github.com/gin-gonic/gin"
)

type UserAward struct {
	Name      string
	AwardName string
	AwardId   int
	Money     int
}

var returnData []models.DrawedRecord
var NoticeChan chan UserAward
var recordChan chan models.DrawedRecord

func init() {
	NoticeChan = make(chan UserAward, 10000)
	recordChan = make(chan models.DrawedRecord, 10000)
	/*go func() {
		logs.Error("save record start")
		for {
			record := <-recordChan
			user := record.LuckyUserFullName
			awardName := record.AwardName

		}

	}()*/
}

func NDraw(c *gin.Context) {
	drawer := c.Query("drawer")
	memo := c.Query("memo")
	leaderCount, err := strconv.Atoi(c.Query("leaderCount"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "领导人数人数错误", "data": nil})
		return
	}

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据失败", "data": nil})
		return
	}
	actionID := c.Query("actionID")

	var awardID int
	var awardName string
	var allPeopleCount int
	var backMoney int
	for i, item := range dataObj.Actions {
		if item.Status == "ToDo" {
			if strconv.Itoa(item.ID) == actionID {
				dataObj.Actions[i].Status = "Done"
				awardID = item.AwardID
				awardName = item.AwardName
				allPeopleCount = item.PeopleCount
				backMoney = item.BackMoney
				break
			} else {
				logs.Error("action:%v", item)
				c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "当前奖品状态不正确", "data": nil})
				return
			}
		} else if i+1 == len(dataObj.Actions) {
			c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "奖品已经抽完", "data": nil})
			return
		}
	}
	money := 0
	for i, item := range dataObj.Awards {
		if item.ID == awardID {
			money = item.Money
			break
		} else if i == len(dataObj.Awards) {
			c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "奖品类型不在字典内", "data": nil})
			return
		}
	}

	staffCount := 0
	if leaderCount >= 0 {

		staffCount = allPeopleCount - leaderCount
		allPeopleCount = 0

	} else {
		allPeopleCount = allPeopleCount
	}
	logs.Info("draw leaderCount:%d,staffCount:%d,actionID:%v", leaderCount, staffCount, actionID)
	if allPeopleCount < 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工人数与领导人数之和 超过设置总人数", "data": nil})
		return
	}
	var temp []models.DrawedRecord
	returnData = temp

	for i := 0; i < leaderCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 2)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}
		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardName:         awardName,
			AwardCount:        1,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}

	for i := 0; i < staffCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}
		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardCount:        1,
			AwardName:         awardName,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}

	for i := 0; i < allPeopleCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 0)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}

		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardCount:        1,
			AwardName:         awardName,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}

	for _, item := range returnData {

		var userAward UserAward
		userAward.Name = item.LuckyUserName
		userAward.AwardName = awardName
		userAward.AwardId = awardID
		userAward.Money = money

		NoticeChan <- userAward

		if item.LuckyUserLevel == 2 {
			backMoneyRecord := models.BackMoneyRecord{
				BackTime: time.Now().Format("2006-01-02 15:04:05"),
				Memo:     "奖品返奖",
				Money:    backMoney,
				UserID:   item.LuckyUserID,
				UserName: item.LuckyUserName,
			}
			dataObj.BackMoneyRecords = append(dataObj.BackMoneyRecords, backMoneyRecord)

		}
	}
	sort.Slice(returnData, func(i, j int) bool {
		if rand.Intn(2) > 0 {
			return true
		}
		return false
	})

	dataObj.SetData()
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": returnData})
	return

}

func ExDraw(c *gin.Context) {
	drawer := c.Query("drawer")
	memo := c.Query("memo")

	awardID, err := strconv.Atoi(c.Query("awardID"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "奖品ID错误", "data": nil})
		return
	}
	awardCount, err := strconv.Atoi(c.Query("awardCount"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "奖品数量错误", "data": nil})
		return
	}
	backMoney, err := strconv.Atoi(c.Query("backMoney"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "返奖金额错误", "data": nil})
		return
	}

	mixPeopleCount, err := strconv.Atoi(c.Query("mixPeopleCount"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "混合人数错误", "data": nil})
		return
	}
	leaderCount, err := strconv.Atoi(c.Query("leaderCount"))
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "2", "msg": "领导人数人数错误", "data": nil})
		return
	}
	staffCount, err := strconv.Atoi(c.Query("staffCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工人数错误", "data": nil})
		return
	}

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据错误", "data": nil})
		return
	}
	awardName := ""
	money := 0
	for i, item := range dataObj.Awards {
		if item.ID == awardID {
			awardName = item.Name
			if awardName == "特别奖" {
				awardName = drawer + awardName + memo

			}
			money = item.Money
			break
		} else if i == len(dataObj.Awards) {
			c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "奖品类型不在字典内", "data": nil})
			return
		}
	}

	var returnData []models.DrawedRecord

	for i := 0; i < mixPeopleCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 0)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}

		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardCount:        awardCount,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}
	for i := 0; i < leaderCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 2)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}
		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardCount:        awardCount,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}
	for i := 0; i < staffCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}
		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardCount:        awardCount,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}

	for _, item := range returnData {
		var userAward UserAward
		userAward.Name = item.LuckyUserName
		userAward.AwardName = awardName
		userAward.AwardId = awardID
		userAward.Money = money
		NoticeChan <- userAward
		if item.LuckyUserLevel == 2 {
			backMoneyRecord := models.BackMoneyRecord{
				BackTime: time.Now().Format("2006-01-02 15:04:05"),
				Memo:     "奖品返奖",
				Money:    backMoney,
				UserID:   item.LuckyUserID,
				UserName: item.LuckyUserName,
			}
			dataObj.BackMoneyRecords = append(dataObj.BackMoneyRecords, backMoneyRecord)
		}
	}
	sort.Slice(returnData, func(i, j int) bool {
		if rand.Intn(2) > 0 {
			return true
		}
		return false
	})
	dataObj.SetData()
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": returnData})
	return

}

func PoolDraw(c *gin.Context) {
	awardID := 888
	drawer := c.Query("drawer")
	memo := c.Query("memo")
	backMoney, err := strconv.Atoi(c.Query("backMoney"))
	if err != nil || backMoney < 0 {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "返奖金额错误", "data": nil})
		return
	}
	awardCount, err := strconv.Atoi(c.Query("awardCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "奖品数量错误", "data": nil})
		return
	}

	mixPeopleCount, err := strconv.Atoi(c.Query("mixPeopleCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "混合人数错误", "data": nil})
		return
	}
	leaderCount, err := strconv.Atoi(c.Query("leaderCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "领导人数人数错误", "data": nil})
		return
	}
	staffCount, err := strconv.Atoi(c.Query("staffCount"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "员工人数错误", "data": nil})
		return
	}
	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "dataobj err", "data": nil})
		return
	}
	//preRecords := dataObj.DrawedRecords[:len(dataObj.DrawedRecords)]

	var returnData []models.DrawedRecord

	for i := 0; i < mixPeopleCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 0)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}

		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardCount:        awardCount,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}
	for i := 0; i < leaderCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 2)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}
		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardCount:        awardCount,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}
	for i := 0; i < staffCount; i++ {
		index, err := GetLuckyUserID(dataObj.Users, 1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": returnData})
			return
		}
		record := models.DrawedRecord{
			AwardID:           awardID,
			AwardCount:        awardCount,
			Drawer:            drawer,
			DrawTime:          time.Now().Format("2006-01-02 15:04:05"),
			LuckyUserID:       dataObj.Users[index].ID,
			LuckyUserLevel:    dataObj.Users[index].Level,
			LuckyUserName:     dataObj.Users[index].Name,
			LuckyUserFullName: dataObj.Users[index].FullName,
			Memo:              memo,
		}
		dataObj.DrawedRecords = append(dataObj.DrawedRecords, record)
		returnData = append(returnData, record)
		dataObj.Users[index].IsDrawed = true
	}
	for _, item := range returnData {
		var userAward UserAward
		userAward.Name = item.LuckyUserName
		userAward.AwardName = "现金奖"
		userAward.AwardId = awardID
		userAward.Money = awardCount
		NoticeChan <- userAward
		if item.LuckyUserLevel == 2 {
			backMoneyRecord := models.BackMoneyRecord{
				BackTime: time.Now().Format("2006-01-02 15:04:05"),
				Memo:     "奖品返奖",
				Money:    backMoney,
				UserID:   item.LuckyUserID,
				UserName: item.LuckyUserName,
			}
			dataObj.BackMoneyRecords = append(dataObj.BackMoneyRecords, backMoneyRecord)
		}
	}
	sort.Slice(returnData, func(i, j int) bool {
		if rand.Intn(2) > 0 {
			return true
		}
		return false
	})

	//dataObj, err = dataObj.PoolDraw(common.GetUUID(), staffCount, leaderCount, mixPeopleCount, awardCount, backRatio, drawer, memo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": err.Error(), "data": ""})
		return
	}
	//currentRecords := dataObj.DrawedRecords[len(preRecords):]
	dataObj.SetData()
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": returnData})
	return

}
