package db

import (
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"log"
	"pha/common"
	"pha/global"

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"gorm.io/gorm"
)

func InitSqlite() {
	var err error
	global.DB, err = gorm.Open(sqlite.Open("workers.db"), &gorm.Config{})
	if err != nil {
		log.Println("init db error:", err)
		common.ErrorExit()
	}
	err = global.DB.AutoMigrate(&Worker{})
	if err != nil {
		log.Println("init worker table error:", err)
		common.ErrorExit()
	}
	err = global.DB.AutoMigrate(&Config{})
	if err != nil {
		log.Println("init config table error:", err)
		common.ErrorExit()
	}
}
