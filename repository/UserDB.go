package repository

import (
	"database/sql"
	"log"

	"github.com/eoinahern/podcastAPI/models"
)

//UserDB : used to do CRUD operations on the users DB table.
type UserDB struct {
	*sql.DB
}

//CheckExist : check user exists in table by users email address.
func (DB *UserDB) CheckExist(email string) bool {

	var count int
	rows := DB.QueryRow("SELECT count(*) FROM users WHERE user_name = ?", email)
	rows.Scan(&count)

	if count == 1 {
		return true
	}

	return false

}

//ValidateUserPlusRegToken : check if user with specific registration exists in table.
func (DB *UserDB) ValidateUserPlusRegToken(email string, regToken string) bool {

	/*var count int
	DB.Model(&models.User{}).Where("user_name = ? AND reg_token = ?", email, regToken).Count(&count)

	if count == 1 {
		return true
	}*/

	return false
}

//SetVerified : set user with specific token and email to be verified in table.
func (DB *UserDB) SetVerified(username string, token string) {

	var user models.User
	row := DB.QueryRow("SELECT * FROM users WHERE user_name = ? AND reg_token = ?", username, token)
	err := row.Scan(&user.UserName, &user.Verified, &user.Password, &user.RegToken)

	if err != nil {
		log.Println(err)
	}

	user.Verified = true

	stmt, err := DB.Prepare("UPDATE users SET verified = ? WHERE user_name= ? AND reg_token = ?")

	if err != nil {
		log.Println(err)
	}

	_, err = stmt.Exec(true, username, token)

	if err != nil {
		log.Println(err)
	}

}

//ValidatePasswordAndUser : check user exists with specific password.
func (DB *UserDB) ValidatePasswordAndUser(email string, password string) bool {

	var user models.User
	//DB.Where("user_name = ? AND password = ?", email, password).First(&user)

	row := DB.QueryRow("SELECT * FROM users WHERE user_name = ? AND password = ?", email, password)
	row.Scan(&user.UserName, &user.Verified, &user.Password, &user.RegToken)

	if user.UserName == email {
		return true
	}

	return false
}

//Insert : Add new user to the users table.
func (DB *UserDB) Insert(user *models.User) {
	//DB.Save(user)

	stmt, err := DB.Prepare("INSERT into users(username, verified, password, reg_token) VALUES(?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(user.UserName, user.Verified, user.Password, user.RegToken)

	if err != nil {
		log.Fatal(err)
	}

}

//GetUser returns a user based on its email.
func (DB *UserDB) GetUser(email string) models.User {

	/*var user models.User
	DB.Where("user_name = ?", email).First(&user)
	return user*/
	return models.User{}
}

func (DB *UserDB) delete(email string) bool {
	return true
}
