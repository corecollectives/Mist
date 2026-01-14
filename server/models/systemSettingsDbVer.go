package models

type SystemSettingsDbVer struct {
	key       string `gorm:"primaryKey"`
	value     string `gorm:"not null"`
	updatedAt int64  `gorm:"autoUpdateTime"`
}
