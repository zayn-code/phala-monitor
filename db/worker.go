package db

import "time"

type Worker struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Uuid        string    `json:"uuid" gorm:"uuid"`
	Pid         uint      `json:"pid" gorm:"pid"`
	Name        string    `json:"name" gorm:"name"`
	Status      string    `json:"status" gorm:"status"`
	State       string    `json:"state" gorm:"state"`
	LastMessage string    `json:"lastMessage" gorm:"lastMessage"`
	Endpoint    string    `json:"endpoint" gorm:"endpoint"`
	Ve          float64   `json:"ve" gorm:"ve"`
	V           float64   `json:"v" gorm:"v"`
	TotalReward float64   `json:"totalReward" gorm:"totalReward"`
	Date        string    `json:"date" gorm:"date"`
	CreatedAt   time.Time `json:"createdAt" gorm:"createdAt"`
}

func (Worker) TableName() string {
	return "worker"
}
