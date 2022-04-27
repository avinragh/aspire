package db

import (
	"aspire/models"
	"time"

	"github.com/go-openapi/strfmt"
)

func (db *DB) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	sqlFindByEmail := `
			SELECT id,username,password,email,role,createdOn,modifiedOn FROM users WHERE email=$1`
	err := db.QueryRow(sqlFindByEmail, email).Scan(user.ID, user.Username, user.Password, user.Email, user.Role, user.CreatedOn, user.ModifiedOn)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) AddUser(user *models.User) (*models.User, error) {
	currentDate := strfmt.DateTime(time.Now())
	sqlInsert := `
		INSERT INTO users(username,password,email,role,createdOn,modifiedOn)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var id int64
	err := db.QueryRow(sqlInsert, user.Username, user.Password, user.Email, user.Role, currentDate, currentDate).Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID = &id
	return user, nil

}
