package store

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	dsn  string
	gorm *gorm.DB
}

func (postgresDB *PostgresDB) GenerateDSN(dbName string) {
	postgresDB.dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASS"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		dbName)
}

func (postgresDB *PostgresDB) ConnectORM() error {
	gormDb, err := gorm.Open(postgres.Open(postgresDB.dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	postgresDB.gorm = gormDb
	return nil

}

func (postgresDB *PostgresDB) GetGorm() *gorm.DB {
	return postgresDB.gorm
}
