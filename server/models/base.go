package models

import (
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

func init() {
	dbUri := os.Getenv("DB_URI")

	conn, err := gorm.Open("postgres", dbUri)

	if err != nil {
		retries := 0

		for retries < 10 && err != nil {
			time.Sleep(200 * time.Millisecond)

			conn, err = gorm.Open("postgres", dbUri)

			if err == nil {
				break
			}

			retries++
		}

		if retries >= 10 {
			panic("Failed to connect to database")
		}

	}

	db = conn

	// auto migrate what needs to be migrated
	db.AutoMigrate(&User{}, &Nonce{})
}

func GetDB() *gorm.DB {
	return db
}
