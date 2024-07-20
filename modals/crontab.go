package modals

import "time"

type CronJob struct {
	JobId    uint      `gorm:"primaryKey"`
	Command  string    `gorm:"type:text;not null"`
	Schedule string    `gorm:"type:text;not null"`
	NextRun  time.Time `gorm:"not null"`
}
