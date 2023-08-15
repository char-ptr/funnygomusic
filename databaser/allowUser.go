package databaser

import "gorm.io/gorm"

type AllowedUser struct {
	gorm.Model
	ID uint64 `gorm:"primaryKey"`
}
