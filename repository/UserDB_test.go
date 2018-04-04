package repository

import (
	"testing"

	"github.com/eoinahern/podcastAPI/models"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestExist(t *testing.T) {

	db, mock, err := sqlmock.New()
	defer db.Close()

	if err != nil {
		panic(err)
	} //helper func

	userDb := UserDB{db}
	rows := sqlmock.NewRows([]string{"user_name"}).AddRow(1)

	mock.ExpectQuery(`SELECT`).WithArgs("hello").WillReturnRows(rows)
	val := userDb.CheckExist("hello")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, true, val)

}

func TestSetVerified(t *testing.T) {

	t.Parallel()

	db, mock, err := sqlmock.New()

	defer db.Close()

	if err != nil {
		panic(err)
	}

	userDB := UserDB{db}
	rows := sqlmock.NewRows([]string{"user_name", "verified", "password", "reg_token"}).AddRow("eoin", true, "pass", "12345")
	mock.ExpectQuery(`SELECT \* FROM users`).WithArgs("eoin", "12345").WillReturnRows(rows)
	mock.ExpectPrepare("UPDATE users SET")
	mock.ExpectExec("UPDATE users SET").WithArgs(true, "eoin", "12345").WillReturnResult(sqlmock.NewResult(1, 1))
	userDB.SetVerified("eoin", "12345")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

}

func TestValidatePassAndUser(t *testing.T) {

	t.Parallel()

	db, mock, err := sqlmock.New()

	defer db.Close()

	if err != nil {
		panic(err)
	}

	userDB := UserDB{db}

	row := sqlmock.NewRows([]string{"user_name", "verified", "password", "reg_token"}).AddRow("eoin", true, "pass", "token")
	mock.ExpectQuery(`SELECT \* FROM users WHERE user_name`).WithArgs("eoin", "pass").WillReturnRows(row)
	val := userDB.ValidatePasswordAndUser("eoin", "pass")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, true, val)

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

func TestGetUser(t *testing.T) {

	t.Parallel()

	db, mock, err := sqlmock.New()

	defer db.Close()

	if err != nil {
		panic(err)
	}

	userDB := UserDB{db}
	row := sqlmock.NewRows([]string{"user_name", "verified", "password", "reg_token"}).AddRow("eoin", true, "pass", "token")
	mock.ExpectQuery(`SELECT \* FROM users`).WithArgs("eoin").WillReturnRows(row)
	user := userDB.GetUser("eoin")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, user.UserName, "eoin")
	assert.Equal(t, user.Password, "pass")

}

func TestValidateUserPlusRegToken(t *testing.T) {

	t.Parallel()

	db, mock, err := sqlmock.New()

	defer db.Close()

	if err != nil {
		panic(err)
	}

	userDB := UserDB{db}
	rows := sqlmock.NewRows([]string{"user_name"}).AddRow(1)

	mock.ExpectQuery(`SELECT`).WithArgs("hello", "token").WillReturnRows(rows)
	val := userDB.ValidateUserPlusRegToken("hello", "token")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("err %s", err)
	}

	assert.Equal(t, true, val)

}
