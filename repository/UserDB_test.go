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

	db, mock, _ = sqlmock.New()
	gormDB, err := gorm.Open("mysql", db)

	if err != nil {
		panic(err)
	}

	userDB = UserDB{gormDB}

}

func TestExist(t *testing.T) {

	row := sqlmock.NewRows([]string{"user_name"}).AddRow("hello")
	mock.ExpectQuery(`SELECT count\(\*\) FROM`).WithArgs("hello").WillReturnRows(row)
	userDB.CheckExist("hello")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

}
