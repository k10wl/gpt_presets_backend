package database







import (
	"gpt_presets_backend/internal/models"
	"log"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func Connect() *gorm.DB {
	var err error

	once.Do(func() {
		dsn := os.Getenv("DB_DSN")

		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	})

	if err != nil {
		log.Fatal("Failed to initialize database connection")
	}

	db.AutoMigrate(models.User{}, models.Tokens{})

	return db
}

func Disconnect(db *gorm.DB) {
	dbConnection, err := db.DB()
	if err == nil {
		dbConnection.Close()
		log.Println("Disconnected from database")
	}
	log.Println("Database connection not established")
}

func Instance() *gorm.DB {
	if db == nil {
		log.Fatal("DB connection not initialized")
	}

	return db
}
