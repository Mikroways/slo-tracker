package store

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDB struct {
	dsn  string
	gorm *gorm.DB
}

func (mysqlDB *MysqlDB) GenerateDSN(dbName string) {
	mysqlDB.dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASS"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		dbName,
	)
}

func (mysqlDB *MysqlDB) ConnectORM() error {
	gormDb, err := gorm.Open(mysql.Open(mysqlDB.dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	mysqlDB.gorm = gormDb
	return nil
}

func (mysqlDB *MysqlDB) GetGorm() *gorm.DB {
	return mysqlDB.gorm
}
