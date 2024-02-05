package database

import (
	"fmt"

	"github.com/lucaslucyk/krowi/pkg/config"
	"github.com/lucaslucyk/krowi/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() {
	var err error

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error to connect to database")
	}
	fmt.Println("DB connection opened!")

	if err = DB.AutoMigrate(&models.User{}); err != nil {
		panic("Error migrating data")
	}
	fmt.Println("Database migrated!")
}
