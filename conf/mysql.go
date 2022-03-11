package conf

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

func ConnectMysql() (err error) {
	dsn := ":root:123456@tcp(stevenwang.top:3306)/Todo?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}