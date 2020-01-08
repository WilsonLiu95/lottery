package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var config *Config = nil

func GetConfig() *Config {
	if config == nil {
		LoadConfig()
	}
	return config
}

type Config struct {
	Server struct {
		Listen struct {
			Ip       string `json:"ip"`
			HttpPort int    `json:"http_port"`
			Port     int    `json:"port"`
		} `json:"listen"`

		NumCPU int `json:"num_cpu"`
	} `json:"server"`
	Log struct {
		FileName string `json:"log_file_name"`
		Log4Go   struct {
			Output  string `json:"output"`
			Level   string `json:"level"`
			MaxSize int    `json:"maxsize"`
		} `json:"log4go"`
	} `json:"log"`

	NoticeUrl string `json:"notice_url"`
}

func LoadConfig() error {

	cfgFileName := "./conf/lottery.conf"

	data, err := ioutil.ReadFile(cfgFileName)
	if err != nil {
		fmt.Println("loadconfig err", err)
		return err
	}
	config = new(Config)
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("loadconfig Unmarshal err", err)
		return err
	}

	return nil
}

// func SaveGoldConfig() error {
// 	body, err := json.Marshal(gold_config)
// 	if err != nil {
// 		glog.Error("saveconfig marshal err", err)
// 		return err
// 	}

// 	err = ioutil.WriteFile(goldFileName, body, 0644)
// 	if err != nil {
// 		glog.Error("writeconfig err", err)
// 		return err
// 	}

// 	glog.Info("savegoldconfig success special media:%s", string(body))

// 	return nil
// }
