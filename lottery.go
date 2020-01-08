package main

import (
	"fmt"
	"net/http"
	"os"

	"strconv"

	"./api"
	"./common"

	"./sockects"

	"github.com/gin-gonic/gin"
)

var log = common.Loger

func main() {
	initCfg()
	router := gin.Default()

	//api路由
	router.Group("/api/count").GET("/", api.GetCount)
	router.Group("/api/getAwards").GET("/", api.GetAwards)
	router.Group("/api/initData").GET("/", api.InitData)
	router.Group("/api/getNextAction").GET("/", api.GetNextAction)
	router.Group("/api/ndraw").GET("/", api.NDraw)
	router.Group("/api/exdraw").GET("/", api.ExDraw)
	router.Group("/api/pooldraw").GET("/", api.PoolDraw)
	router.Group("/api/addMoney").GET("/", api.AddPoolMoney)
	router.Group("/api/initSystem").GET("/", api.InitData)
	router.Group("/api/getusers").GET("/", api.GetUsers)
	router.Group("/api/getresult").GET("/", api.GetResult)
	router.Group("/api/getpoolmoney").GET("/", api.GetPoolmoney)
	router.Group("/api/export_excel").GET("/", api.ExportExcel)
	//web socket 路由
	router.GET("/ws", func(c *gin.Context) { sockects.WebSocketHandler(c.Writer, c.Request) })

	//html页面路由
	router.LoadHTMLGlob("views/*")
	router.Group("/view/").GET("/:name", func(c *gin.Context) { c.HTML(http.StatusOK, c.Param("name")+".html", gin.H{}) })

	////静态文件路由
	router.Static("/wwwroot", "./wwwroot")
	port := common.GetConfig().Server.Listen.HttpPort
	log.Info("server run...")
	router.Run(":" + strconv.Itoa(port))

}
func initCfg() {

	err := common.LoadConfig()

	if err != nil {
		fmt.Println("load config file err", err)
		os.Exit(0)
	}

	err = api.LoadData()

	if err != nil {
		fmt.Println("LoadData err", err)
		os.Exit(0)
	}
}
