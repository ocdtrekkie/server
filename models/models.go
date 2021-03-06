package models

import (
	"database/sql"

	// custom packages
	cError "github.com/joyread/server/error"
)

// CreateUser ...
func CreateUser(db *sql.DB) {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS account (id BIGSERIAL PRIMARY KEY, username VARCHAR(255) UNIQUE NOT NULL, email VARCHAR(255) UNIQUE NOT NULL, password_hash VARCHAR(255) NOT NULL, jwt_token VARCHAR(255) NOT NULL, is_admin BOOLEAN NOT NULL DEFAULT FALSE, storage VARCHAR(255) NOT NULL DEFAULT 'none')")
	cError.CheckError(err)
}

// SignUpModel struct
type SignUpModel struct {
	Username     string
	Email        string
	PasswordHash string
	Token        string
	IsAdmin      bool
}

// InsertUser ...
func InsertUser(db *sql.DB, signUpModel SignUpModel) int {
	var lastInsertID int
	err := db.QueryRow("INSERT INTO account (username, email, password_hash, jwt_token, is_admin) VALUES ($1, $2, $3, $4, $5) returning id", signUpModel.Username, signUpModel.Email, signUpModel.PasswordHash, signUpModel.Token, signUpModel.IsAdmin).Scan(&lastInsertID)
	cError.CheckError(err)

	return lastInsertID
}

// SelectAdmin ...
func SelectAdmin(db *sql.DB) int {
	// Check for Admin in the user table
	rows, err := db.Query("SELECT id FROM account WHERE is_admin = $1", true)
	cError.CheckError(err)

	var userID int

	if rows.Next() {
		err := rows.Scan(&userID)
		cError.CheckError(err)
	}
	rows.Close()

	return userID
}

// SelectPasswordHashAndJWTTokenModel struct
type SelectPasswordHashAndJWTTokenModel struct {
	UsernameOrEmail string
}

// SelectPasswordHashAndJWTTokenResponse struct
type SelectPasswordHashAndJWTTokenResponse struct {
	PasswordHash string
	Token        string
}

// SelectPasswordHashAndJWTToken ...
func SelectPasswordHashAndJWTToken(db *sql.DB, selectPasswordHashAndJWTTokenModel SelectPasswordHashAndJWTTokenModel) *SelectPasswordHashAndJWTTokenResponse {
	// Search for username in the 'account' table with the given string
	rows, err := db.Query("SELECT password_hash, jwt_token FROM account WHERE username = $1", selectPasswordHashAndJWTTokenModel.UsernameOrEmail)
	cError.CheckError(err)

	var selectPasswordHashAndJWTTokenResponse SelectPasswordHashAndJWTTokenResponse

	if rows.Next() {
		err := rows.Scan(&selectPasswordHashAndJWTTokenResponse.PasswordHash, &selectPasswordHashAndJWTTokenResponse.Token)
		cError.CheckError(err)
	} else {
		// if username doesn't exist, search for email in the 'account' table with the given string
		rows, err := db.Query("SELECT password_hash, jwt_token FROM account WHERE email = $1", selectPasswordHashAndJWTTokenModel.UsernameOrEmail)
		cError.CheckError(err)

		if rows.Next() {
			err := rows.Scan(&selectPasswordHashAndJWTTokenResponse.PasswordHash, &selectPasswordHashAndJWTTokenResponse.Token)
			cError.CheckError(err)
		}
		rows.Close()
	}
	rows.Close()

	return &selectPasswordHashAndJWTTokenResponse
}

// CreateSMTP ...
func CreateSMTP(db *sql.DB) {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS smtp (hostname VARCHAR(255) NOT NULL, port INTEGER NOT NULL, username VARCHAR(255) NOT NULL, password VARCHAR(255) NOT NULL)")
	cError.CheckError(err)
}

// SMTPModel struct
type SMTPModel struct {
	Hostname string
	Port     int
	Username string
	Password string
}

// InsertSMTP ...
func InsertSMTP(db *sql.DB, smtpModel SMTPModel) {
	_, err := db.Query("INSERT INTO smtp (hostname, port, username, password) VALUES ($1, $2, $3, $4)", smtpModel.Hostname, smtpModel.Port, smtpModel.Username, smtpModel.Password)
	cError.CheckError(err)
}

// CheckSMTP ...
func CheckSMTP(db *sql.DB) bool {

	// Check for Admin in the user table
	rows, err := db.Query("SELECT hostname FROM smtp")
	cError.CheckError(err)

	var isSMTPPresent = false

	if rows.Next() {
		isSMTPPresent = true
	}
	rows.Close()

	return isSMTPPresent
}

// CreateNextcloud ...
func CreateNextcloud(db *sql.DB) {
	_, err := db.Query("CREATE TABLE IF NOT EXISTS nextcloud (id BIGSERIAL, user_id INTEGER REFERENCES account(id), url VARCHAR(255) NOT NULL, client_id VARCHAR(1200) NOT NULL, client_secret VARCHAR(1200) NOT NULL, directory VARCHAR(255) NOT NULL, redirect_uri VARCHAR(255) NOT NULL, access_token VARCHAR(255), refresh_token VARCHAR(255), PRIMARY KEY (id, user_id))")
	cError.CheckError(err)
}

// NextcloudModel struct
type NextcloudModel struct {
	UserID       int
	URL          string
	ClientID     string
	ClientSecret string
	Directory    string
	RedirectURI  string
}

// InsertNextcloud ...
func InsertNextcloud(db *sql.DB, nextcloudModel NextcloudModel) {
	_, err := db.Query("INSERT INTO nextcloud (user_id, url, client_id, client_secret, directory, redirect_uri) VALUES ($1, $2, $3, $4, $5, $6)", nextcloudModel.UserID, nextcloudModel.URL, nextcloudModel.ClientID, nextcloudModel.ClientSecret, nextcloudModel.Directory, nextcloudModel.RedirectURI)
	cError.CheckError(err)

	_, err = db.Query("UPDATE account SET storage=$1 WHERE id=$2", "nextcloud", nextcloudModel.UserID)
	cError.CheckError(err)
}

// SelectNextcloudModel struct
type SelectNextcloudModel struct {
	UserID int
}

// SelectNextcloudResponse struct
type SelectNextcloudResponse struct {
	URL          string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// SelectNextcloud ...
func SelectNextcloud(db *sql.DB, selectNextcloudModel SelectNextcloudModel) *SelectNextcloudResponse {
	rows, err := db.Query("SELECT url, client_id, client_secret, redirect_uri FROM nextcloud WHERE user_id=$1", selectNextcloudModel.UserID)
	cError.CheckError(err)

	var selectNextcloudResponse SelectNextcloudResponse

	if rows.Next() {
		err := rows.Scan(&selectNextcloudResponse.URL, &selectNextcloudResponse.ClientID, &selectNextcloudResponse.ClientSecret, &selectNextcloudResponse.RedirectURI)
		cError.CheckError(err)
	}
	rows.Close()

	return &selectNextcloudResponse
}

// NextcloudTokenModel struct
type NextcloudTokenModel struct {
	AccessToken  string
	RefreshToken string
	UserID       int
}

// UpdateNextcloudToken ...
func UpdateNextcloudToken(db *sql.DB, nextcloudTokenModel NextcloudTokenModel) {
	_, err := db.Query("UPDATE nextcloud SET access_token=$1, refresh_token=$2 WHERE user_id=$3", nextcloudTokenModel.AccessToken, nextcloudTokenModel.RefreshToken, nextcloudTokenModel.UserID)
	cError.CheckError(err)
}

// CheckStorage ...
func CheckStorage(db *sql.DB, userID int) string {
	rows, err := db.Query("SELECT storage FROM account WHERE id=$1", userID)
	cError.CheckError(err)

	var storage string

	if rows.Next() {
		err := rows.Scan(&storage)
		cError.CheckError(err)
	}
	rows.Close()

	return storage
}

// CheckNextcloudToken ...
func CheckNextcloudToken(db *sql.DB, userID int) string {
	rows, err := db.Query("SELECT access_token FROM nextcloud WHERE user_id=$1", userID)
	cError.CheckError(err)

	var accessToken string

	if rows.Next() {
		err := rows.Scan(&accessToken)
		cError.CheckError(err)
	}
	rows.Close()

	return accessToken
}
