package global

import (
	"gorm.io/gorm"
)

type MailConfigType struct {
	From     string
	To       []string
	Subject  string
	Host     string
	Port     int
	Username string
	Password string
}

type PrbConfigType struct {
	Origin string
	PeerId string
}

var (
	DB         *gorm.DB
	MailConfig MailConfigType

	PrbConfig PrbConfigType
)

const (
	GoDate     = "2006-01-02"
	GoDateTime = "2006-01-02 15:04:05"

	IgnoreWorkers = "ignore_workers"
)
