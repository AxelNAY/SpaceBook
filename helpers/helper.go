package helpers

import (
	"testing"

	"spacebook/config"
	"spacebook/models"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func SetupTestDB() {
	godotenv.Load("../.env")
	if config.DB == nil {
		config.ConnectDatabase()
	}
}

func CreateTestCompany(t *testing.T, name string) models.Company {
	t.Helper()
	company := models.Company{
		ID:   uuid.New(),
		Name: name,
	}
	config.DB.Create(&company)
	return company
}

func CreateTestPlace(t *testing.T, companyID uuid.UUID, name, description string) models.Place {
	t.Helper()
	place := models.Place{
		ID:          uuid.New(),
		CompanyID:   companyID,
		Name:        name,
		Description: description,
	}
	config.DB.Create(&place)
	return place
}

func CreateTestUser(t *testing.T, email, username string) models.User {
	t.Helper()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		Password: hashedPassword,
		Role:     "user",
	}
	config.DB.Create(&user)
	return user
}

func CreateTestResource(t *testing.T, name string) models.Resource {
	t.Helper()
	resource := models.Resource{
		Name:   name,
		Status: "available",
	}
	config.DB.Create(&resource)
	return resource
}

func CreateTestNotification(t *testing.T, userID *uuid.UUID, message string) models.Notification {
	t.Helper()
	notification := models.Notification{
		UserID:  userID,
		Type:    "reservation",
		Message: message,
		IsRead:  false,
	}
	config.DB.Create(&notification)
	return notification
}
