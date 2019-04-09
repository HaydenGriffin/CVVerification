package database

import (
	"database/sql"
	"fmt"
	"github.com/cvverification/chaincode/model"
	"github.com/cvverification/models"
)

var dataSourceName = "root:password@tcp(localhost:3306)/verification?parseTime=true"

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(100)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetUserDetailsFromUsername(username string) (models.UserDetails, error) {

	user := models.UserDetails{}

	db, err := InitDB(dataSourceName)

	if err != nil {
		return user, err
	}

	result := db.QueryRow("SELECT u.id, u.username, u.full_name, u.email_address, u.profile_id FROM users u WHERE username = ?", username)
	err = result.Scan(&user.Id, &user.Username, &user.FullName, &user.EmailAddress, &user.ID)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func GetAllReviewableCVHashes() (map[int] string, error) {

	ratableCVs := make(map[int] string)

	db, err := InitDB(dataSourceName)

	if err != nil {
		return ratableCVs, err
	}

	rows, err := db.Query("SELECT u.id, uc.cv_hash FROM user_cvs uc JOIN users u ON uc.user_id = u.id WHERE uc.cv_in_review = 1")

	for rows.Next() {
		var cvHash string
		var userID int
		err = rows.Scan(&userID, &cvHash)
		if err != nil {
			rows.Close()
			return ratableCVs, err
		}
		ratableCVs[userID] = cvHash
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return ratableCVs, err
	}
	return ratableCVs, nil
}

func GetUserCVDetails(user_id int) ([]model.CVHistoryInfo, error) {

	var historicalCVHistoryInfo []model.CVHistoryInfo

	db, err := InitDB(dataSourceName)

	if err != nil {
		return historicalCVHistoryInfo, err
	}

	rows, err := db.Query("SELECT uc.cv_hash, uc.cv_in_review, uc.timestamp FROM user_cvs uc WHERE uc.user_id = ? ORDER BY uc.timestamp ASC", user_id)

	var index = 1

	for rows.Next() {
		var cvHistoryInfo model.CVHistoryInfo
		err = rows.Scan(&cvHistoryInfo.CVHash, &cvHistoryInfo.CVInReview, &cvHistoryInfo.Timestamp)
		if err != nil {
			rows.Close()
			fmt.Println(err.Error())
			return historicalCVHistoryInfo, err
		}
		cvHistoryInfo.Index = index
		historicalCVHistoryInfo = append(historicalCVHistoryInfo, cvHistoryInfo)
		index++
	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		return historicalCVHistoryInfo, err
	}
	return historicalCVHistoryInfo, nil
}

func CreateNewUser(username, full_name, email_address, profile_id string) (userDetails models.UserDetails, error error) {

	db, err := InitDB(dataSourceName)

	if err != nil {
		return userDetails, err
	}

	res, err := db.Exec("INSERT INTO users(username, full_name, email_address, profile_id) VALUES (?, ?, ?, ?)", username, full_name, email_address, profile_id)
	fmt.Println(res)

	if err != nil {
		return userDetails, err
	}

	var selectedUser models.UserDetails

	selectedUser, err = GetUserDetailsFromUsername(username)

	if err != nil {
		return userDetails, err
	}

	userDetails = selectedUser

	return userDetails, err
}

func UpdateUser(username, full_name, email_address string) (userDetails models.UserDetails, error error) {

	db, err := InitDB(dataSourceName)

	if err != nil {
		return userDetails, err
	}

	user, err := GetUserDetailsFromUsername(username)

	if err != nil {
		return models.UserDetails{}, err
	}

	_, err = db.Exec("UPDATE users SET full_name = ?, email_address = ? WHERE id = ?", full_name, email_address, user.Id)

	if err != nil {
		return models.UserDetails{}, err
	} else {
		user.FullName = full_name
		user.EmailAddress = email_address
	}

	return userDetails, err
}

func GetCVInfoFromID(user_id int) (string, string, error) {
	db, err := InitDB(dataSourceName)

	if err != nil {
		return "", "", err
	}

	result := db.QueryRow("SELECT u.profile_id, uc.cv_hash FROM users u JOIN user_cvs uc ON u.id = uc.user_id WHERE u.id = ? AND uc.cv_in_review = 1", user_id)

	var ID, cvHash string

	err = result.Scan(&ID, &cvHash)

	if err != nil {
		return "", "", err
	}

	return ID, cvHash, nil
}

func CreateNewCV(user_id int, cv_hash string) error {
	db, err := InitDB(dataSourceName)

	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO user_cvs(user_id, timestamp, cv_hash, cv_in_review) VALUES (?, CURRENT_TIMESTAMP, ?, 0)", user_id, cv_hash)

	return err
}

func UpdateCV(cv_hash string, ratable int) error {
	db, err := InitDB(dataSourceName)

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE user_cvs SET cv_in_review = ? WHERE cv_hash = ?", ratable, cv_hash)

	return err
}

func UserHasCVInReview(user_id int) bool {
	db, err := InitDB(dataSourceName)

	if err != nil {
		return false
	}

	result := db.QueryRow("SELECT cv_in_review FROM user_cvs WHERE user_id = ? AND cv_in_review = 1", user_id)

	var cv_ratable int

	err = result.Scan(&cv_ratable)

	if err != nil {
		return false
	}

	if cv_ratable != 1 {
		return false
	} else {
		return true
	}
}
