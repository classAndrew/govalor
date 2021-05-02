package models

import (
	"log"
	"os"

	// for mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DB .
var DB *gorm.DB

// ConnectDatabase Connects to db
func ConnectDatabase() {
	dbpass := os.Getenv("DBPASS")
	dbname := os.Getenv("DBNAME")
	dbuser := os.Getenv("DBUSER")
	dbhost := os.Getenv("DBHOST")
	database, err := gorm.Open("mysql", dbuser+":"+dbpass+"@tcp("+dbhost+")/"+dbname)
	if err != nil {
		log.Println(err)
	}
	database.AutoMigrate(&UserSliceXP{})
	database.AutoMigrate(&UserTotalXP{})
	database.AutoMigrate(&ActivityMember{})
	database.AutoMigrate(&MemberRecordXP{})
	database.AutoMigrate(&GuildMember{})
	DB = database
}
