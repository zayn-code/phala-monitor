package db

import (
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"pha/global"

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

func InitSqlite() {
	var err error
	global.DB, err = gorm.Open(sqlite.Open("data/workers.db"), &gorm.Config{})
	if err != nil {
		panic("init db error:" + err.Error())
	}
	err = global.DB.AutoMigrate(&Worker{})
	if err != nil {
		panic("init worker table error:" + err.Error())
	}
	err = global.DB.AutoMigrate(&Config{})
	if err != nil {
		panic("init config table error:" + err.Error())
	}
}
