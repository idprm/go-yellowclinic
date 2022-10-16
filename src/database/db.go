package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/idprm/go-yellowclinic/src/config"
	"github.com/idprm/go-yellowclinic/src/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Datasource *NewDatasource

type NewDatasource struct {
	db    *gorm.DB
	sqlDb *sql.DB
}

func (d NewDatasource) DB() *gorm.DB {
	return d.db
}

func (d NewDatasource) SqlDB() *sql.DB {
	return d.sqlDb
}

func Connect() {

	var db *gorm.DB
	var sqlDb *sql.DB

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		config.ViperEnv("DB_USER"),
		config.ViperEnv("DB_PASS"),
		config.ViperEnv("DB_HOST"),
		config.ViperEnv("DB_PORT"),
		config.ViperEnv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		panic("Could not connect with the database!")
	}

	sqlDb, _ = db.DB()
	sqlDb.SetConnMaxLifetime(time.Minute * 2)
	sqlDb.SetMaxOpenConns(10000)
	sqlDb.SetMaxIdleConns(10000)

	// try to establish connection
	if sqlDb != nil {
		err = sqlDb.Ping()
		if err != nil {
			log.Fatal("cannot connect to db:", err.Error())
		}
	}

	log.Println("Connected to database successfully")

	if err != nil {
		log.Fatal("Error! \n", err.Error())
	}

	// DEBUG ON CONSOLE
	db.Logger = logger.Default.LogMode(logger.Info)

	// TODO: Add migrations
	db.AutoMigrate(
		&model.User{},
		&model.Config{},
		&model.Doctor{},
		&model.Chat{},
	)

	// TODO: Add seeders
	var config []model.Config
	var doctor []model.Doctor

	resultConfig := db.Find(&config)
	resultDoctor := db.Find(&doctor)

	if resultConfig.RowsAffected == 0 {
		for i, _ := range configs {
			db.Model(&model.Config{}).Create(&configs[i])
		}
	}

	if resultDoctor.RowsAffected == 0 {
		for i, _ := range doctors {
			db.Model(&model.Doctor{}).Create(&doctors[i])
		}
	}

	Datasource = &NewDatasource{db: db, sqlDb: sqlDb}
}
