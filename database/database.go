package database

import (
	"database/sql"
	"fmt"
	"github.com/cvtracker/models"
)

var dataSourceName = "root:password@tcp(localhost:3306)/verification"

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetUserFromUsername(username string) (models.User, error) {

	user := models.User{}

	db, err := InitDB(dataSourceName)

	result := db.QueryRow("SELECT u.id, u.username, u.full_name, u.password, u.email_address, u.user_role  FROM users u WHERE username = ?", username)
	err = result.Scan(&user.Id, &user.Username, &user.FullName, &user.Password, &user.EmailAddress, &user.UserRole)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func GetUserFromId(id string) (models.User, error) {

	user := models.User{}

	db, err := InitDB(dataSourceName)

	result := db.QueryRow("SELECT u.id, u.username, u.full_name, u.password, u.email_address, u.user_role  FROM users u WHERE id = ?", id)
	err = result.Scan(&user.Id, &user.Username, &user.FullName, &user.Password, &user.EmailAddress, &user.UserRole)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func GetCVHashFromUserID(id int) (string, error) {

	var cvHash string

	db, err := InitDB(dataSourceName)

	result := db.QueryRow("SELECT uc.cv_hash FROM user_cvs uc WHERE user_id = ?", id)
	err = result.Scan(&cvHash)

	if err != nil {
		return "", err
	} else {
		return cvHash, nil
	}
}

func CreateNewUser(username, full_name, password, email_address string) error {

	db, err := InitDB(dataSourceName)
	res, err := db.Exec("INSERT INTO users(username, full_name, password, email_address, user_role) VALUES (?, ?, ?, ?, ?)", username, full_name, password, email_address, "APPLICANT")
	fmt.Println(res)

	return err
}