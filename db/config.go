package db

import (
	"errors"
	"gorm.io/gorm"
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
	if err := global.DB.Where("key = ?", k).First(&config).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		_ = SaveConfig(k, "")
	}
	return config.Value
}

func SaveConfig(k string, v string) error {
	err := global.DB.Where("key = ?", k).Save(&Config{
		Key:   k,
		Value: v,
	})
	if err != nil {
		return err.Error
	}
	return nil
}
