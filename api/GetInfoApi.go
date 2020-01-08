package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAwards(c *gin.Context) {

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取award数据错误", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": dataObj.Awards})
	return
}
func GetUsers(c *gin.Context) {

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取user数据错误", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": dataObj.Users})
	return
}

func GetPoolmoney(c *gin.Context) {

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取poolmoney数据错误", "data": nil})
		return
	}

	dataObj.Count.PoolMoney = 0
	for _, item := range dataObj.BackMoneyRecords {
		dataObj.Count.PoolMoney += item.Money
	}
	for _, item := range dataObj.DrawedRecords {
		if item.AwardID == 0 {
			dataObj.Count.PoolMoney -= item.AwardCount
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "poolmoney": dataObj.Count.PoolMoney})
	return
}
func GetResult(c *gin.Context) {

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取result数据错误", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": returnData})
	return
}
func GetCount(c *gin.Context) {

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据错误", "data": nil})
		return
	}
	dataObj.CountData()
	c.JSON(http.StatusCreated, gin.H{"code": "0", "count": dataObj.Count})
}

func GetNextAction(c *gin.Context) {

	if dataObj == nil {
		c.JSON(http.StatusOK, gin.H{"code": "2", "msg": "获取txt数据错误", "data": nil})
		return
	}
	for _, item := range dataObj.Actions {
		if item.Status == "ToDo" {
			c.JSON(http.StatusOK, gin.H{"code": "0", "msg": "ok", "data": item})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": "1", "msg": "奖品已抽完"})
	return
}
