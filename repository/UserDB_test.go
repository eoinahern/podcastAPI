package repository

import (
	"database/sql"
	"testing"

	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var db *sql.DB
var mock sqlmock.Sqlmock
var gormDB *gorm.DB
var userDB UserDB

func init() {

	db, mock, err := sqlmock.New()

	if err != nil {
		panic(err)
	}

	gormDB, err = gorm.Open("mysql", db)

	if err != nil {
		panic(err)
	}

	userDB = UserDB{gormDB}

}

func tests(t *testing.T) {

}
