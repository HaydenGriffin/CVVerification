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

	if err != nil {
		return user, err
	}

	result := db.QueryRow("SELECT u.id, u.username, u.full_name, u.password, u.email_address, u.user_role, u.profile_hash  FROM users u WHERE username = ?", username)
	err = result.Scan(&user.Id, &user.Username, &user.FullName, &user.Password, &user.EmailAddress, &user.UserRole, &user.ProfileHash)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func GetAllRatableCVHashes() (map[int] string, error) {

	ratableCVs := make(map[int] string)

	db, err := InitDB(dataSourceName)

	if err != nil {
		return ratableCVs, err
	}

	rows, err := db.Query("SELECT u.id, uc.cv_hash FROM user_cvs uc JOIN users u ON uc.user_id = u.id WHERE uc.cv_ratable = 1")

	for rows.Next() {
		var cvHash string
		var userID int
		err = rows.Scan(&userID, &cvHash)
		if err != nil {
			return ratableCVs, err
		}
		ratableCVs[userID] = cvHash
	}
	err = rows.Err()
	if err != nil {
		return ratableCVs, err
	}
	return ratableCVs, nil
}

func CreateNewUser(username, full_name, password, email_address, user_role, profile_hash string) (user models.User, error error) {

	db, err := InitDB(dataSourceName)

	if err != nil {
		return user, err
	}

	res, err := db.Exec("INSERT INTO users(username, full_name, password, email_address, user_role, profile_hash) VALUES (?, ?, ?, ?, ?, ?)", username, full_name, password, email_address, user_role, profile_hash)
	fmt.Println(res)

	if err != nil {
		return user, err
	}

	var selectedUser models.User

	selectedUser, err = GetUserFromUsername(username)

	if err != nil {
		return user, err
	}

	user = selectedUser

	return user, err
}

func GetCVInfoFromID(user_id int) (string, string, error) {
	db, err := InitDB(dataSourceName)

	if err != nil {
		return "", "", err
	}

	result := db.QueryRow("SELECT u.profile_hash, uc.cv_hash FROM users u JOIN user_cvs uc ON u.id = uc.user_id WHERE u.id = ? AND uc.status_control = 'C'", user_id)

	var profileHash, cvHash string

	err = result.Scan(&profileHash, &cvHash)

	if err != nil {
		return "", "", err
	}

	return profileHash, cvHash, nil
}

func CreateNewCV(user_id int, cv, cv_hash string) error {
	db, err := InitDB(dataSourceName)

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE user_cvs SET status_control = 'H' WHERE user_id = ?", user_id)

	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO user_cvs(user_id, timestamp, cv, cv_hash, cv_ratable, status_control) VALUES (?, CURRENT_TIMESTAMP, ?, ?, 0, 'C')", user_id, cv, cv_hash)

	return err
}

func UpdateCV(cv_hash string, ratable int) error {
	db, err := InitDB(dataSourceName)

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE user_cvs SET cv_ratable = ? WHERE cv_hash = ? AND status_control = 'C'", ratable, cv_hash)

	return err
}

func IsCVRatable(cv_hash string) (bool, error) {
	db, err := InitDB(dataSourceName)

	if err != nil {
		return false, err
	}

	result := db.QueryRow("SELECT cv_ratable FROM user_cvs WHERE cv_hash = ? AND status_control = 'C'", cv_hash)

	var cv_ratable int

	err = result.Scan(&cv_ratable)

	if err != nil {
		return false, err
	}

	if cv_ratable == 1 {
		return true, nil
	} else {
		return false, nil
	}
}