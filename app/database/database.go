package database

import (
	"database/sql"
	"fmt"
	templateModel "github.com/cvverification/app/model"
)

var DataSourceName = "root:password@tcp(localhost:3306)/verification?parseTime=true"

type MySQL struct {

}

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(100)

	return nil
}

func GetUserDetailsFromUsername(username string) (templateModel.UserDetails, error) {

	user := templateModel.UserDetails{}


	result := db.QueryRow("SELECT u.id, u.username, u.full_name, u.email_address FROM users u WHERE username = ?", username)
	err := result.Scan(&user.Id, &user.Username, &user.FullName, &user.EmailAddress)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func CreateNewUser(username, full_name, email_address, fabric_id string) (userDetails templateModel.UserDetails, error error) {

	res, err := db.Exec("INSERT INTO users(username, full_name, email_address, fabric_id) VALUES (?, ?, ?, ?)", username, full_name, email_address, fabric_id)
	fmt.Println(res)

	if err != nil {
		return userDetails, err
	}

	var selectedUser templateModel.UserDetails

	selectedUser, err = GetUserDetailsFromUsername(username)

	if err != nil {
		return userDetails, err
	}

	userDetails = selectedUser

	return userDetails, err
}

func UpdateUser(username, full_name, email_address string) (userDetails templateModel.UserDetails, error error) {

	user, err := GetUserDetailsFromUsername(username)

	if err != nil {
		return templateModel.UserDetails{}, err
	}

	_, err = db.Exec("UPDATE users SET full_name = ?, email_address = ? WHERE id = ?", full_name, email_address, user.Id)

	if err != nil {
		return templateModel.UserDetails{}, err
	} else {
		user.FullName = full_name
		user.EmailAddress = email_address
	}

	return userDetails, err
}

func GetFabricIDFromCVID(cvID string) (string, error) {

	result := db.QueryRow("SELECT u.fabric_id FROM users u JOIN user_cvs uc ON u.id = uc.user_id WHERE uc.cv_id = ?", cvID)

	var fabricID string

	err := result.Scan(&fabricID)

	if err != nil {
		return "", err
	}

	return fabricID, nil
}

func CreateNewCV(user_id int, cv_id string) error {

	_, err := db.Exec("INSERT INTO user_cvs(user_id, timestamp, cv_id) VALUES (?, CURRENT_TIMESTAMP, ?)", user_id, cv_id)

	return err
}