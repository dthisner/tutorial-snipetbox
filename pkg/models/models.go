package models

import (
	"errors"
	"time"
)

var (
	// ErrNoRecord is a new error when no matching records were located
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials provieds an error when providing wrong login credentials
	ErrInvalidCredentials = errors.New("models: Invalid Credentials")
	// ErrDuplicateEmail provides error when user tries signing up with an existing email
	ErrDuplicateEmail = errors.New("models: duplicate email")
	// ErrDontExist provides error when user tries signing in with an email that dosn't exist
	ErrDontExist = errors.New("models: Email dosn't exist in the database")
	// ErrNotActive provides error when user tries signing in with an disabled account
	ErrNotActive = errors.New("models: Your account is deactivated, contact admin")
)

// Snippet has the structure for our snippets
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// User holds the struct for any signed up user
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
