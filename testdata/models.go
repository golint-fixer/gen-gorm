package testdata

import (
	// need this
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type profile struct {
	gorm.Model
	Name string
}

type user struct {
	gorm.Model
	Profile      profile `gorm:"ForeignKey:ProfileRefer"` // use ProfileRefer as foreign key
	ProfileRefer int
}

func createTables() {
	db, err := gorm.Open("mysql", "root:@/gorm_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	//db.Model(&User{}).AddForeignKey("profile_refer", "profile(id)", "RESTRICT", "RESTRICT")

	db.AutoMigrate(&user{}, &profile{})

}
