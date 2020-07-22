package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"herobix.com/snipetbox/pkg/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// UserModel holds database connection to SQL
type UserModel struct {
	DB *sql.DB
}

// Insert adds  a user into the DB, returns an error, if any
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Authenticate checks if the provided email and password macth any records
// and returns relevant User ID if they do
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	var active bool

	// Check if user exist
	stmt := "SELECT id, hashed_password, active FROM users WHERE email = ?"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword, &active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrDontExist
		}
		return 0, err
	}

	// Check if users account is disabled or not
	if !active {
		return 0, models.ErrNotActive
	}

	// Check if password matches
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

// Get lets you fetch the users information, based on the id
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}

	stmt := `SELECT id, name, email, created, active FROM users WHERE id =?`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return u, nil
}

// ChangePassword checks user provided correct password and if so, update password with
// with a new one
func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
	var currentHassedPassword []byte

	row := m.DB.QueryRow("SELECT hashed_password FROM users WHERE id = ?", id)
	err := row.Scan(&currentHassedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNoRecord
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword(currentHassedPassword, []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return models.ErrInvalidCredentials
		}
		return err
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 15)
	if err != nil {
		return err
	}

	stmt := "UPDATE users SET hashed_password = ? WHERE id = ?"
	_, err = m.DB.Exec(stmt, string(newHashedPassword), id)
	if err != nil {
		return err
	}

	return nil
}
