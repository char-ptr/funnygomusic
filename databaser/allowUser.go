package databaser

import "gorm.io/gorm"

type AllowedUser struct {
	gorm.Model
	UserId uint64
}
