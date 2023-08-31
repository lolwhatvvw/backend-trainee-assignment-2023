package postgres

import (
	"fmt"
	"log"

	"github.com/lolwhatvvw/backend-trainee-assignment-2023/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(conf config.Config) (*gorm.DB, error) {
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		conf.Database.Host,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.DbName,
		conf.Database.Port,
	)

	log.Println(connString)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
