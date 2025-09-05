package configs

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Env.DBHost,
		Env.DBUser,
		Env.DBPass,
		Env.DBName,
		Env.DBPort,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,    // Slow SQL threshold
			LogLevel:                  logger.Silent,  // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      IsProduction(), // Don't include params in the SQL log
			Colorful:                  true,           // Disable color
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   false, // skip the snake_casing of names
		},
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect PostgreSQL: %v", err)
	}
	sDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}
	sDB.SetMaxIdleConns(10)
	sDB.SetMaxOpenConns(100)
	DB = db
	fmt.Println("✅ Database connected successfully!")
}
