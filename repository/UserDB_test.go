package repository

import (
	"database/sql"
	"testing"

	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var db *sql.DB
var mocksss sqlmock.Sqlmock
var gormDB *gorm.DB
var userDB UserDB

func init() {

	db, mocksss, _ = sqlmock.New()

	gormDB, err := gorm.Open("mysql", db)

	if err != nil {
		panic(err)
	}

	userDB = UserDB{gormDB}

}

func TestExist(t *testing.T) {

	row := sqlmock.NewRows([]string{"user_name"}).AddRow("hello")
	mocksss.ExpectQuery(`SELECT count(\\*)`).WithArgs("hello").WillReturnRows(row)
	userDB.CheckExist("hello")

	if err := mocksss.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

}
