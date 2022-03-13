package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/fabrv/watchman-server/models"
	"github.com/fabrv/watchman-server/utils"
)

var (
	DBConn *gorm.DB
)

func InitDatabase() {
	var err error
	DBConn, err = gorm.Open("postgres", utils.GetEnv("DATABASE_URL", "host=localhost port=5432 user=postgres dbname=watchman password=password sslmode=disable"))
	DBConn.LogMode(true)

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")

	DBConn.AutoMigrate(&models.LogType{}, &models.Project{}, &models.Team{}, &models.Role{}, &models.User{})
	fmt.Println("Database migrated")
}
