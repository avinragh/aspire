package web

import (
	"aspire/models"
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
	if err != nil {
		err = errors.New("Email already in use")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return

	}
	hashPassword, err := GeneratehashPassword(*user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}

	user.Password = &hashPassword

	//insert user details in database
	user, err = database.AddUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateJWT(userId int, role string) (string, error) {
	var mySigningKey = []byte("secretkey")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userId"] = strconv.Itoa(userId)
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 90).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
