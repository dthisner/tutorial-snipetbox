package mock

import (
	"time"

	"herobix.com/snipetbox/pkg/models"
)

var mockUser = &models.User{
	ID:      1,
	Name:    "Alice",
	Email:   "alice@example.com",
	Created: time.Now(),
	Active:  true,
}

// UserModel is a mock struct for testing the User model
type UserModel struct{}

// Insert is a mock call for testing the Insert of a user function
func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

// Authenticate is a mock call for testing authenticating the user
func (m *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "alice@example.com":
		return 1, nil
	case "notExisting@example.com":
		return 0, models.ErrDontExist
	case "notActive@example.com":
		return 0, models.ErrNotActive
	}
	return 0, models.ErrInvalidCredentials
}

// Get is a mock call for testing getting specific user
func (m *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// ChangePassword is a mock call for changing users password
func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	return nil
}
