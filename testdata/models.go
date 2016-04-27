package testdata

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Profile struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Profile      Profile `gorm:"ForeignKey:ProfileRefer"` // use ProfileRefer as foreign key
	ProfileRefer int
}

func CreateTables() {
	db, err := gorm.Open("mysql", "root:@/gorm_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	//db.Model(&User{}).AddForeignKey("profile_refer", "profile(id)", "RESTRICT", "RESTRICT")

	db.AutoMigrate(&User{}, &Profile{})

}
