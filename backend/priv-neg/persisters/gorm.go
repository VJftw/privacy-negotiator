package persisters

import (
	"fmt"
	"log"
	"os"

	"github.com/VJftw/privacy-negotiator/backend/priv-neg/utils"
	"github.com/jinzhu/gorm"
	// postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// NewGORMDB - Initialises a connection to a GORM storage
func NewGORMDB(logger *log.Logger, models ...interface{}) *gorm.DB {

	if !utils.WaitForService(fmt.Sprintf("%s:%s", os.Getenv("POSTGRES_HOST"), "5432"), logger) {
		panic("Could not find database")
	}

	// db, err := gorm.Open("sqlite3", "test.db")
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
	))

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.AutoMigrate(models...)

	db.LogMode(true)

	return db
}
