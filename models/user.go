package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`                       //用户ID
	Username string `json:"username" form:"username" gorm:"unique;not null;size:128" binding:"required"` //用户名称
	Password string `json:"password" form:"password" gorm:"not null;size:256" binding:"required"`        //用户密码
}

// CreateUserRecord creates a user record in the database
func (user *User) CreateUserRecord() error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) CreateDefaultUser() error {
	user.HashPassword("123456")
	var DefaultUsers = []User{
		{Username: "admin", Password: user.Password},
	}
	for _, d := range DefaultUsers{
		p, _ :=  user.GetUserByName(d.Username)
		if p != nil{
			continue
		}
		err := d.CreateUserRecord()
		if err != nil{
			continue
		}
	}
	return nil
}

func (user *User) GetUserByName(name string) (*User, error) {
	var ret User
	if err := db.Where("username = ?", name).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

// HashPassword encrypts user password
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword checks user password
func (user *User) CheckPassword(username, password string) error {
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

// CheckNewPassword checks new password
func (user *User) CheckNewPassword(username, newPassword string) bool {
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err != nil {
		return true
	}
	return false
}

func (user *User) GetPassword(username string) string {
	db.Where("username = ?", username).First(&user)
	return user.Password
}

func (user *User) UpdatePassword(useName, newPassword string) string {
	db.Model(user).Where("username = ?", useName).Update("password", newPassword)
	return user.Password
}

// RevokeAccountByName  RevokeAccount
func (user *User) RevokeAccountByName(username string) error {
	if err := db.Where("username = ?", username).Delete(user).Error; err != nil {
		return err
	}
	return nil
}
