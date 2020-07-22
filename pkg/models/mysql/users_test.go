package mysql

import (
	"reflect"
	"testing"
	"time"

	"herobix.com/snipetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: Skipping integration test")
	}

	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:  true,
			},
			wantError: nil,
		},
		{
			name:      "Zero ID",
			userID:    0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{
			name:      "Non-existent ID",
			userID:    -1,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}

			user, err := m.Get(tt.userID)
			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
		})
	}
}

func TestUserModelChangePassword(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: Skipping integration test")
	}

	tests := []struct {
		name            string
		userID          int
		currentPassword string
		newPassword     string
		wantError       error
	}{
		{
			name:            "Valid Password",
			userID:          1,
			currentPassword: "Min21Pirate!",
			newPassword:     "NewAmazingPassword",
			wantError:       nil,
		},
		{
			name:            "Not Matching Password",
			userID:          1,
			currentPassword: "WrongPassword",
			newPassword:     "NewAmazingPassword",
			wantError:       models.ErrInvalidCredentials,
		},
		{
			name:            "No Matching Records",
			userID:          -1,
			currentPassword: "Min21Pirate!",
			newPassword:     "NewAmazingPassword",
			wantError:       models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}

			err := m.ChangePassword(tt.userID, tt.currentPassword, tt.newPassword)
			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}
		})
	}
}

func TestUserModelAuthenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: Skipping integration test")
	}

	tests := []struct {
		name       string
		email      string
		password   string
		wantError  error
		wantUserID int
	}{
		{
			name:       "Valid Password",
			email:      "alice@example.com",
			password:   "Min21Pirate!",
			wantError:  nil,
			wantUserID: 1,
		},
		{
			name:       "Invalid Password",
			email:      "alice@example.com",
			password:   "NotValidAtAll",
			wantError:  models.ErrInvalidCredentials,
			wantUserID: 0,
		},
		{
			name:       "Not Active",
			email:      "notactive@example.com",
			password:   "Min21Pirate!",
			wantError:  models.ErrNotActive,
			wantUserID: 0,
		},
		{
			name:       "Dosn't Exist",
			email:      "notExisting@example.com",
			password:   "Min21Pirate!",
			wantError:  models.ErrDontExist,
			wantUserID: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}

			id, err := m.Authenticate(tt.email, tt.password)
			if err != tt.wantError {
				t.Errorf("want %v; got %v", tt.wantError, err)
			}

			if id != tt.wantUserID {
				t.Errorf("want %v; got %v", tt.wantUserID, id)
			}
		})
	}
}
