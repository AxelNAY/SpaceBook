package helpers

import (
	"testing"
	"time"

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
		Status:   "Actif",
	}
	config.DB.Create(&user)
	return user
}

func CreateTestPendingUser(t *testing.T, email, username string) models.User {
	t.Helper()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		Password: hashedPassword,
		Role:     "user",
		Status:   "En attente",
	}
	config.DB.Create(&user)
	return user
}

func CreateTestAdmin(t *testing.T, email, username string) models.User {
	t.Helper()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		Password: hashedPassword,
		Role:     "admin",
		Status:   "Actif",
	}
	config.DB.Create(&user)
	return user
}

func CreateTestResource(t *testing.T, name string) models.Resource {
	t.Helper()
	resource := models.Resource{
		ID:     uuid.New(),
		Name:   name,
		Status: "available",
	}
	config.DB.Create(&resource)
	return resource
}

func CreateTestResourceInPlace(t *testing.T, placeID uuid.UUID, name string) models.Resource {
	t.Helper()
	resource := models.Resource{
		ID:      uuid.New(),
		PlaceID: &placeID,
		Name:    name,
		Status:  "available",
	}
	config.DB.Create(&resource)
	return resource
}

func CreateTestCategory(t *testing.T, name string) models.Category {
	t.Helper()
	category := models.Category{
		ID:   uuid.New(),
		Name: name,
	}
	config.DB.Create(&category)
	return category
}

func CreateTestPlaceUser(t *testing.T, placeID, userID uuid.UUID) models.PlaceUser {
	t.Helper()
	placeUser := models.PlaceUser{
		ID:      uuid.New(),
		PlaceID: placeID,
		UserID:  userID,
	}
	config.DB.Create(&placeUser)
	return placeUser
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

func CreateTestReservation(t *testing.T, userID, resourceID uuid.UUID) models.Reservation {
	t.Helper()
	now := time.Now()
	start := now.Add(24 * time.Hour)
	end := now.Add(25 * time.Hour)
	reservation := models.Reservation{
		ID:            uuid.New(),
		UserID:        userID,
		RessourceID:   &resourceID,
		StartDatetime: start,
		EndDatetime:   end,
		Status:        "pending",
	}
	config.DB.Create(&reservation)
	return reservation
}
