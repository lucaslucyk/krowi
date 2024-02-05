package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lucaslucyk/krowi/pkg/database"
	"github.com/lucaslucyk/krowi/pkg/models"
)

func CreateUser(user *models.User) error {
	db := database.DB

	result := db.Create(&user)
	if result.Error != nil {
		errMsg := result.Error.Error()
		if strings.Contains(errMsg, "duplicate key value violates unique") {
			return errors.New(ALREADY_EXISTS)
		}
		return fmt.Errorf("Error creating user: %s", errMsg)
	}

	return nil
}

func GetUserByEmail(dest *models.User, email string) error {
	// var user models.User
	db := database.DB

	result := db.First(&dest, "email = ?", email)
	if result.Error != nil {
		return fmt.Errorf(NOT_FOUND)
	}
	return nil
}
