package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitDB(logger *logrus.Logger) (*gorm.DB, error) {
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	database := viper.GetString("mysql.database")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	charset := viper.GetString("mysql.charset")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", user, password, host, port, database, charset)
	logger.Infof("dsn info %s\n", dsn)
	return initDatabase(dsn)
}

func getDBLogger() logger.Interface {
	var Logger logger.Interface
	if gin.Mode() != gin.DebugMode {
		Logger = logger.Default.LogMode(logger.Info)
	} else {
		Logger = logger.Default
	}
	return Logger
}

func initDatabase(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: getDBLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	// set db attributes
	if sqlConn, err := db.DB(); err != nil {
		return nil, err
	} else {
		sqlConn.SetMaxIdleConns(10)
		sqlConn.SetMaxOpenConns(100)
		sqlConn.SetConnMaxLifetime(30 * time.Second)
	}
	// create the table User automatically if it doesn't exist
	if err = db.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return db, nil
}
