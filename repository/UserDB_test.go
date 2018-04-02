package repository

import (
	"testing"

	"github.com/eoinahern/podcastAPI/models"

	_ "github.com/go-sql-driver/mysql"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestExist(t *testing.T) {

	db, mock, err := sqlmock.New()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	userDb := UserDB{db}
	row := sqlmock.NewRows([]string{"user_name"})

	row.AddRow("hello")
	mock.ExpectQuery(`SELECT`).WithArgs("hello")
	userDb.CheckExist("hello")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}
}

func TestSetVerified(t *testing.T) {

	db, mock, err := sqlmock.New()

	defer db.Close()

	if err != nil {
		panic(err)
	}

	userDB := UserDB{db}
	mock.ExpectQuery(`SELECT \* FROM users`).WithArgs("eoin", "12345")
	userDB.SetVerified("eoin", "12345")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

}

func TestValidatePassAndUser(t *testing.T) {

}

func TestInsert(t *testing.T) {

	db, mock, err := sqlmock.New()

	defer db.Close()

	if err != nil {
		panic(err)
	}

	userDB := UserDB{db}
	mock.ExpectPrepare("INSERT into users")
	mock.ExpectExec("INSERT into users").WithArgs("eoin", true, "pass", "boo").WillReturnResult(sqlmock.NewResult(1, 1))

	user := &models.User{UserName: "eoin", Verified: true, Password: "pass", RegToken: "boo"}
	userDB.Insert(user)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

}
