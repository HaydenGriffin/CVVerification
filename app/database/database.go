package database

import (
	"database/sql"
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

	result := db.QueryRow("SELECT u.id, u.username, u.title, u.first_name, u.surname, u.email_address, u.date_of_birth FROM users u WHERE username = ?", username)
	err := result.Scan(&user.Id, &user.Username, &user.Title, &user.FirstName, &user.Surname, &user.EmailAddress, &user.DateOfBirth)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func CreateNewUser(username, title, first_name, surname, email_address, date_of_birth, fabric_id string) (userDetails templateModel.UserDetails, error error) {

	_, err := db.Exec("INSERT INTO users(username, title, first_name, surname, email_address, date_of_birth, fabric_id) VALUES (?, ?, ?, ?, ?, ?, ?)", username, title, first_name, surname, email_address, date_of_birth, fabric_id)
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

func UpdateUser(username, title, first_name, surname, email_address, date_of_birth string) (userDetails templateModel.UserDetails, error error) {

	user, err := GetUserDetailsFromUsername(username)
	if err != nil {
		return templateModel.UserDetails{}, err
	}

	_, err = db.Exec("UPDATE users SET title = ?, first_name = ?, surname = ?, email_address = ?, date_of_birth = ? WHERE id = ?", title, first_name, surname, email_address, date_of_birth, user.Id)
	if err != nil {
		return templateModel.UserDetails{}, err
	} else {
		user.Title = title
		user.FirstName = first_name
		user.Surname = surname
		user.EmailAddress = email_address
		user.DateOfBirth = date_of_birth
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
