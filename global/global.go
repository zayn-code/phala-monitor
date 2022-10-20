package global

import "gorm.io/gorm"

var (
	DB *gorm.DB
)

const (
	GoDate     = "2006-01-02"
	GoDateTime = "2006-01-02 15:04:05"

	IgnoreWorkers = "ignore_workers"
)
