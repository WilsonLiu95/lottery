package api

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"../models"
)

func init() {
	rand.Seed(time.Now().Unix())
}

///////
func GetLuckyUserID(users []models.User, level int) (luckyID int, err error) {
	//len := len(users)
	canDrawCount := 0

	for _, item := range users {
		if item.IsDrawed == false && (level == 0 || item.Level == level) {
			canDrawCount++

		}
	}

	if canDrawCount <= 0 {
		return 0, errors.New("对应level:" + strconv.Itoa(level) + "未中奖人数不足")
	}

	step := rand.Intn(canDrawCount) + 1
	//logs.Info("step:%d", step)

	for index, item := range users {
		if item.IsDrawed == false && (level == 0 || item.Level == level) {
			step--
			if step <= 0 {
				return index, nil

			}
		}

	}

	return 0, errors.New("对应level:" + strconv.Itoa(level) + " 抽奖失败")

}
