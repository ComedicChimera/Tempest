package models

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Token represents a standard JWT token
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// User represents a logged in user
type User struct {
	gorm.Model
	IPAddr string `gorm:"column:ipaddr"`
	Token  string `gorm:"column:token"`
}

func LoginUser(addr string) (string, error) {
	user := &User{}

	err := db.Table("users").Where("ipaddr = ?", addr).First(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			user.IPAddr = addr

			db.Create(user)

			if user.ID <= 0 {
				return "", errors.New("Failed to create new user")
			}
		} else {
			return "", errors.New("Unable to connect to database")
		}
	}

	tk := &Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("TEMPEST_TOKEN_PWD")))

	user.Token = tokenString

	return tokenString, nil
}
