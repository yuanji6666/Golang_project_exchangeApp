package global

import(
	"gorm.io/gorm"
)

var (
	Db *gorm.DB //define global database
)