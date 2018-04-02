package repository

import (
	"fmt"
	"log"

	"github.com/eoinahern/podcastAPI/models"

	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
)

//UserDB : used to do CRUD operations on the users DB table.
type UserDB struct {
	*gorm.DB
}

//CheckExist : check user exists in table by users email address.
func (DB *UserDB) CheckExist(email string) bool {

	var count int
	DB.Model(&models.User{}).Where("user_name = ?", email).Count(&count)

	if count >= 1 {
		return true
	}

	return false
}

//ValidateUserPlusRegToken : check if user with specific registration exists in table.
func (DB *UserDB) ValidateUserPlusRegToken(email string, regToken string) bool {

	var count int
	DB.Model(&models.User{}).Where("user_name = ? AND reg_token = ?", email, regToken).Count(&count)

	if count == 1 {
		return true
	}

	return false
}

//SetVerified : set user with specific token and email to be verified in table.
func (DB *UserDB) SetVerified(username string, token string) {

	var user models.User
	DB.Where("user_name = ? AND reg_token = ?", username, token).First(&user)
	user.Verified = true
	db := DB.Save(&user)

	if db.Error != nil {
		log.Println(db.Error)
	}

}

//ValidatePasswordAndUser : check user exists with specific password.
func (DB *UserDB) ValidatePasswordAndUser(email string, password string) bool {

	var user models.User
	DB.Where("user_name = ? AND password = ?", email, password).First(&user)

	fmt.Println(user.UserName)

	if user.UserName == email {
		return true
	}

	return false
}

//Insert : Add new user to the users table.
func (DB *UserDB) Insert(user *models.User) {
	DB.Save(user)
}

//GetUser returns a user based on its email.
func (DB *UserDB) GetUser(email string) models.User {

	var user models.User
	DB.Where("user_name = ?", email).First(&user)
	return user
}

func (DB *UserDB) delete(email string) bool {
	return true
}
