package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(mysql.Open("root:@tcp(localhost:3306)/abskrywn1_db"))
	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&Karyawan{}, &User{}, &Absensi{})

	DB = database
}
