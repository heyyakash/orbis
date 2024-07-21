package db

import (
	"fmt"
	"log"
	"time"

	"github.com/heyyakash/orbis/helpers"
	"github.com/heyyakash/orbis/modals"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	Host     = helpers.GetString("POSTGRES_HOST")
	Database = helpers.GetString("POSTGRES_DB")
	User     = helpers.GetString("POSTGRES_USER")
	Password = helpers.GetString("POSTGRES_PASSWORD")
	Port     = helpers.GetString("POSTGRES_PORT")
)

type PostgresStore struct {
	DB *gorm.DB
}

var Store PostgresStore

func (p *PostgresStore) CreateTable() {
	if err := p.DB.AutoMigrate(&modals.CronJob{}); err != nil {
		log.Fatal("Couldn't Migrate : ", err)
	}
}

func Init() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", Host, User, Password, Database, Port)

	// start connection to db
	Store.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Couldn't connect to DB : ", err)
	}
	log.Print("DB Connected")

	// create table if not exists
	Store.CreateTable()
}

func (p *PostgresStore) DeleteRowById(id uint, modal interface{}) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&modals.CronJob{}).Error; err != nil {
			return err
		}

		if err := tx.Where("job_id = ?", id).Delete(&modals.CronJob{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (p *PostgresStore) UpdateTimeById(key string, job *modals.CronJob, next_time time.Time, modal interface{}) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(job).Error; err != nil {
			return err
		}

		job.NextRun = next_time
		if err := tx.Save(job).Error; err != nil {
			return err
		}
		return nil
	})
}
