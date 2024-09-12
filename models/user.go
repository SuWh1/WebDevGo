package models // This layer is responsible for defining the structure of
// data (such as users, products, etc.),
// interacting with the database, and implementing
// the core logic associated with this data.

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) { // creating user
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // hashed password of the userr
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}
	// we inserting into database the ready signup data
	row := us.DB.QueryRow(`
		INSERT INTO users(email, password_hash)
		VALUES ($1, $2) RETURNING id`, email, passwordHash)

	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	// data ready to query
	row := us.DB.QueryRow(`SELECT id, password_hash FROM users WHERE email=$1`, email)

	// int row query returns id and hashed password if email matches, no we need to scan data
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}

	// we checked the email but we need to verify wheather it is real user by password
	// compare input password with password in database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("auth: %w", err)
	}
	// if we have valid password and do not have error upper:
	return &user, nil
}
