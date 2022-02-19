package conf

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

type Configuration struct {
	WebPort                     int    `json:"webPort"`                     //端口
	DbPath                      string `json:"dbPath"`                      //数据库目录
	TempPath                    string `json:"tempPath"`                    //上传文件的目录
	IotImagePath                string `json:"iotImagePath"`                //上传文件的目录
	IotAlarmMaxCount            int    `json:"iotAlarmMaxCount"`            //要求100-20000之间， 不得超过20000
	IotEventMaxCount            int    `json:"iotEventMaxCount"`            //要求100-20000之间， 不得超过20000
	IotPropsMaxCount            int    `json:"iotPropsMaxCount"`            //要求100-20000之间， 不得超过20000
	IotEventPushHttptUrl        string `json:"iotEventPushHttptUrl"`        //设备事件推送
	IotAlarmPushHttptUrl        string `json:"iotAlarmPushHttptUrl"`        //盒子IOT告警推送
	IotPropsPushHttptUrl        string `json:"iotPropsPushHttptUrl"`        //设备属性推送
	IotDeviceStatusPushHttptUrl string `json:"iotDeviceStatusPushHttptUrl"` //设备心跳推送
	DefaultSN                   string `json:"defaultSN"`                   //sn
	LogLevel                    string `json:"logLevel"`                    //日志级别
}

var confBean *Configuration = nil

func GetConf() *Configuration {
	return confBean
}

func InitConf(fname string) (*Configuration, error) {
	if confBean == nil {
		exists, err := PathExists(fname)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		if !exists {
			return nil, errors.New("Path is not exists.")
		}

		file, err := os.Open(fname)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		reader := json.NewDecoder(file)
		confBean = new(Configuration)
		err = reader.Decode(confBean)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}
	return confBean, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
