package db

import (
	"pha/global"
)

type Config struct {
	Key   string `json:"key" gorm:"key"`
	Value string `json:"value" gorm:"value"`
}

func (Config) TableName() string {
	return "config"
}

//获取配置
func GetConfig(k string) string {
	var config Config
	global.DB.Where("key = ?", k).First(&config)
	return config.Value
}

func SaveConfig(k string, v string) error {
	err := global.DB.Table("config").Where("key = ?", k).Save(&Config{
		Key:   k,
		Value: v,
	})
	if err != nil {
		return err.Error
	}
	return nil
}
