package db

import (
	"aspire/models"
	"log"
	"time"

	"github.com/go-openapi/strfmt"
)

func (db *DB) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	sqlFindByEmail := `
			SELECT id,username,password,email,role,created_on,modified_on FROM users WHERE email=$1`
	err := db.QueryRow(sqlFindByEmail, email).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.CreatedOn, &user.ModifiedOn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func (db *DB) AddUser(user *models.User) (*models.User, error) {
	currentDate := strfmt.DateTime(time.Now())
	user.CreatedOn = currentDate
	user.ModifiedOn = currentDate
	sqlInsert := `
		INSERT INTO users(username,password,email,role,created_on,modified_on)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var id int64
	err := db.QueryRow(sqlInsert, user.Username, user.Password, user.Email, user.Role, user.CreatedOn, user.ModifiedOn).Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID = &id
	return user, nil
}

func (db *DB) DeleteUser(Id int64) error {
	sqlDelete := `
	DELETE FROM users WHERE ID=$1;`
	_, err := db.Exec(sqlDelete, Id)
	if err != nil {
		return err
	}
	return nil
}
