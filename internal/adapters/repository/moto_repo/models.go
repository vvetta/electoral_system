package motorepo

import (
	"time"

	"gorm.io/gorm"
)

type GormMoto struct {
	ID uint `gorm:"primaryKey;autoIncrement;unique"`
	Name string `gorm:"type:varchar(100)"`
	Year int `gorm:"not null"`
	Mileage int `gorm:"not null"`
	EngineSize int `gorm:"not null"`
	MotoType string	`gorm:"type:varchar(255)"`
	Location string `gorm:"type:varchar(255)"`
	Price int64 `gorm:"not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

func (GormMoto) TableName() string {
	return "motos"
}
