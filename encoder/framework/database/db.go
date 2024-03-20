package database

import (
	"log"

	"github.com/alvesrafa/video-encoder/domain"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type Database struct {
	Db            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDatabase() *Database {
	return &Database{}
}

func NewDatabaseTest() *gorm.DB {
	dbInstance := NewDatabase()
	dbInstance.Env = "test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.debug = true

	connection, err := dbInstance.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (database *Database) Connect() (*gorm.DB, error) {

	var err error

	if database.Env != "test" {
		database.Db, err = gorm.Open(database.DbType, database.Dsn)
	} else {
		database.Db, err = gorm.Open(database.DbTypeTest, database.DsnTest)
	}

	if err != nil {
		return nil, err
	}

	if database.debug {
		database.Db.LogMode(true)
	}

	if database.AutoMigrateDb {
		database.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
		database.Db.Model(domain.Job{}).AddForeignKey("video_id", "videos {id}", "CASCADE", "CASCADE")
	}

	return database.Db, nil

}
