package web

import (
	"aspire/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (siw *ServerInterfaceWrapper) Signup(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()
	logger := ctx.GetLogger()
	database := ctx.GetDB()
	user := &models.User{}
	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
		return
	}

	_, err = database.FindUserByEmail(*user.Email)
	if err != nil && err == sql.ErrNoRows {
		hashPassword, err := GeneratehashPassword(*user.Password)
		if err != nil {
			log.Fatalln("error in password hash")
		}
		user.Password = &hashPassword

		user, err = database.AddUser(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	}
	logger.Println(err)
	err = errors.New("Email already in use")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
	return

}

func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()
	database := ctx.GetDB()

	authdetails := &models.Authentication{}
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		err = errors.New("Error in reading body")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	user, err := database.FindUserByEmail(*authdetails.Email)
	if err != nil && err == sql.ErrNoRows {
		err = errors.New("Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	check := CheckPasswordHash(*authdetails.Password, *user.Password)

	if !check {
		err = errors.New("Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	validToken, err := GenerateJWT(*user.ID, user.Role)
	if err != nil {
		err = errors.New("Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var token models.Token
	token.Email = user.Email
	token.Role = user.Role
	token.TokenString = &validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userId int64, role string) (string, error) {
	var mySigningKey = []byte("secretkey")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["userId"] = strconv.FormatInt(userId, 10)
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 90).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
