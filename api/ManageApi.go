package api

import (
	"net/http"
	"strconv"
	"time"

	"../common"

	"../models"
	"github.com/Luxurioust/excelize"
	"github.com/gin-gonic/gin"
)

var logs = common.Loger
var dataObj *models.Data = nil

func LoadData() error {

	dataObj = new(models.Data)
	_, err := dataObj.LoadExistUser()
	_, err = dataObj.LoadDrawedRecord()
	_, err = dataObj.LoadBackMoneyRecord()

	_, err = dataObj.LoadAwards()
	if err != nil {
		logs.Error("load award error:%v", err)
		return err

	}
	_, err = dataObj.LoadAction()
	if err != nil {
		logs.Error("load action error:%v", err)
		return err

	}

	logs.Info("LoadData user size:%d", len(dataObj.Users))

	if dataObj.Awards == nil {
		dataObj.Awards = []models.Award{}
	}

	if dataObj.Actions == nil {
		dataObj.Actions = []models.DrawedAction{}
	}
	if dataObj.Users == nil {
		dataObj.Users = []models.User{}
	}
	//dataObj.DrawedRecords = []models.DrawedRecord{}
	//dataObj.BackMoneyRecords = []models.BackMoneyRecord{}

	dataObj.Count.PoolMoney = 0
	dataObj.Count.AllLuckyCount = 0
	dataObj.Count.AllPeopleCount = 0
	dataObj.Count.LuckyStaffCount = 0
	dataObj.Count.LuckyLeaderCount = 0
	dataObj.Count.NoLuckyLeaderCount = 0
	dataObj.Count.NoLuckyStaffCount = 0
	for _, item := range dataObj.Users {
		dataObj.Count.AllPeopleCount++
		if item.Level == 1 {
			if item.IsDrawed {
				dataObj.Count.LuckyStaffCount++
			} else {
				dataObj.Count.NoLuckyStaffCount++
			}
		} else if item.Level == 2 {
			if item.IsDrawed {
				dataObj.Count.LuckyLeaderCount++
			} else {
				dataObj.Count.NoLuckyLeaderCount++
			}
		}
	}
	logs.Info("data:%v", dataObj)
	return nil

}

func InitData(c *gin.Context) {

	pswd := c.Query("pswd")
	if pswd != "123456" {
		c.JSON(http.StatusCreated, gin.H{"code": "1", "msg": "密码不正确"})
		return
	}

	dataObj = new(models.Data)
	dataObj.LoadUser()
	dataObj.LoadAwards()
	dataObj.LoadAction()

	logs.Info("test,%d", len(dataObj.Users))

	if dataObj.Awards == nil {
		dataObj.Awards = []models.Award{}
	}

	if dataObj.Actions == nil {
		dataObj.Actions = []models.DrawedAction{}
	} else {
		for i := range dataObj.Actions {
			dataObj.Actions[i].Status = "ToDo"

		}
	}
	if dataObj.Users == nil {
		dataObj.Users = []models.User{}
	} else {
		for i := range dataObj.Users {
			dataObj.Users[i].IsDrawed = false
		}
	}
	dataObj.DrawedRecords = []models.DrawedRecord{}
	dataObj.BackMoneyRecords = []models.BackMoneyRecord{}

	dataObj.Count.PoolMoney = 0
	dataObj.Count.AllLuckyCount = 0
	dataObj.Count.AllPeopleCount = 0
	dataObj.Count.LuckyStaffCount = 0
	dataObj.Count.LuckyLeaderCount = 0
	dataObj.Count.NoLuckyLeaderCount = 0
	dataObj.Count.NoLuckyStaffCount = 0
	for _, item := range dataObj.Users {
		dataObj.Count.AllPeopleCount++
		if item.Level == 1 {
			if item.IsDrawed {
				dataObj.Count.LuckyStaffCount++
			} else {
				dataObj.Count.NoLuckyStaffCount++
			}
		} else if item.Level == 2 {
			if item.IsDrawed {
				dataObj.Count.LuckyLeaderCount++
			} else {
				dataObj.Count.NoLuckyLeaderCount++
			}
		}
	}
	logs.Info("data:%v", dataObj)
	setDataErr := dataObj.SetData()
	if setDataErr == nil {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": setDataErr.Error()})
		return
	}
}

func AddPoolMoney(c *gin.Context) {
	money, err := strconv.Atoi(c.Query("money"))
	if money < 1 || err != nil {
		c.JSON(http.StatusCreated, gin.H{"code": "1", "msg": "加奖金额不正确"})
		return
	}
	memo := c.Query("memo")

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取data数据错误", "data": nil})
		return
	}
	dataObj.BackMoneyRecords = append(dataObj.BackMoneyRecords, models.BackMoneyRecord{
		BackTime: time.Now().Format("2006-01-02 15:04:05"),
		Memo:     "临时加奖:" + memo,
		Money:    money,
		UserID:   -1,
		UserName: "",
	})

	setDataErr := dataObj.SetData()
	if setDataErr == nil {
		c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"code": "99", "msg": setDataErr.Error()})
		return
	}
}

func ExportExcel(c *gin.Context) {
	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "-1", "msg": "data为空", "data": nil})
		return
	}
	fileName := "中奖名单.xlsx"
	xlsx := excelize.NewFile()
	xlsx.NewSheet("Sheet1")

	xlsx.SetCellValue("Sheet1", "A1", "中奖人")
	xlsx.SetCellValue("Sheet1", "B1", "奖品名称")
	xlsx.SetCellValue("Sheet1", "C1", "数量")

	for index, record := range dataObj.DrawedRecords {
		luckName := record.LuckyUserFullName
		awardName := record.AwardName
		awardId := record.AwardID
		for _, award := range dataObj.Awards {
			if award.ID == awardId {
				awardName = award.Name
				break
			}

		}
		count := record.AwardCount
		strIndex := strconv.Itoa(index + 2)
		xlsx.SetCellValue("Sheet1", "A"+strIndex, luckName)
		xlsx.SetCellValue("Sheet1", "B"+strIndex, awardName)
		xlsx.SetCellValue("Sheet1", "C"+strIndex, count)

	}
	err := xlsx.SaveAs(fileName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "-1", "msg": "fail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "success"})
}
