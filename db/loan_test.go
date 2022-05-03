package db

import (
	"aspire/models"
	"aspire/util"
	"database/sql"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFindLoanByID(t *testing.T) {
	//setup
	db, err := Init()
	if err != nil {
		t.Errorf("Not able to connect to db")
	}
	user := &models.User{Username: util.GetStringPointer("avinragh"), Password: util.GetStringPointer("Siri2109!"), Role: "admin", Email: util.GetStringPointer("avinragh@gmail.com")}
	user, err = db.AddUser(user)
	if err != nil {
		t.Errorf("Not able to Add user to db for testing")
	}

	spew.Dump(user)

	loan := &models.Loan{Amount: util.GetFloat64Pointer(100.0), Term: util.GetInt64Pointer(3)}
	loan, err = db.AddLoan(loan, *user.ID)
	if err != nil {
		t.Errorf("Not able to Add loan to db for testing")
	}

	spew.Dump(loan)

	spew.Dump(*loan.ID)

	//test cases
	t.Run("SUCCESS", func(t *testing.T) {
		got, err := db.FindLoanById(*loan.ID)
		if err != nil {
			t.Errorf("Not able to Find loan")
		}
		assert.EqualValues(t, loan.ID, got.ID, "The loan result was not as expected")
		assert.EqualValues(t, loan.Amount, got.Amount, "The loan result was not as expected")
		assert.EqualValues(t, loan.Term, got.Term, "The loan result was not as expected")
		assert.EqualValues(t, loan.State, got.State, "The loan result was not as expected")
		assert.EqualValues(t, loan.UserID, got.UserID, "The loan result was not as expected")
	})

	t.Run("FAILURE NOT FOUND", func(t *testing.T) {
		_, err := db.FindLoanById(2)
		assert.EqualErrorf(t, err, sql.ErrNoRows.Error(), "The expected error was not returned")

	})

	// teardown
	db.DeleteUser(*user.ID)
	db.DeleteLoan(*loan.ID)
	db.Close()
}

// func TestAddUser(t *testing.T) {
// 	//setup
// 	db, err := Init()
// 	if err != nil {
// 		t.Errorf("Not able to connect to db")
// 	}
// 	user := &models.User{Username: util.GetStringPointer("avinragh"), Password: util.GetStringPointer("Siri2109!"), Role: "admin", Email: util.GetStringPointer("avinragh@gmail.com")}
// 	failUser := &models.User{Role: "admin", Email: util.GetStringPointer("bob@gmail.com")}
// 	//test cases
// 	t.Run("SUCCESS", func(t *testing.T) {
// 		got, err := db.AddUser(user)
// 		if err != nil {
// 			t.Errorf("Not able to Find user: %s", err.Error())
// 		}
// 		assert.NotEmpty(t, got.ID, "The user result was not as expected")
// 		assert.EqualValues(t, user.Username, got.Username, "The user result was not as expected")
// 		assert.EqualValues(t, user.Role, got.Role, "The user result was not as expected")
// 		assert.EqualValues(t, user.Email, got.Email, "The user result was not as expected")
// 		assert.EqualValues(t, user.Password, got.Password, "The user result was not as expected")

// 	})

// 	t.Run("FAILURE DUPLICATE KEY", func(t *testing.T) {
// 		_, err := db.AddUser(user)
// 		assert.EqualErrorf(t, err, `pq: duplicate key value violates unique constraint "users_email_key"`, "Error Not as Expected")

// 	})

// 	t.Run("FAILURE MISSING FIELDS", func(t *testing.T) {
// 		_, err := db.AddUser(failUser)
// 		assert.EqualErrorf(t, err, `pq: null value in column "password" violates not-null constraint`, "Error Not as Expected")

// 	})

// 	// teardown
// 	user, err = db.FindUserByEmail("avinragh@gmail.com")
// 	if err != nil {
// 		t.Errorf("Unable to teardown AddUserSetup")
// 	}
// 	db.DeleteUser(*user.ID)
// 	db.Close()

// }

// func TestDeleteUser(t *testing.T) {
// 	//setup
// 	db, err := Init()
// 	if err != nil {
// 		t.Errorf("Not able to connect to db")
// 	}
// 	user := &models.User{Username: util.GetStringPointer("avinragh"), Password: util.GetStringPointer("Siri2109!"), Role: "admin", Email: util.GetStringPointer("avinragh@gmail.com")}
// 	user, err = db.AddUser(user)
// 	if err != nil {
// 		t.Errorf("Not able to Add user to db for testing")
// 	}

// 	//run tests
// 	t.Run("SUCCESS", func(t *testing.T) {
// 		err := db.DeleteUser(*user.ID)
// 		if err != nil {
// 			t.Errorf("Not able to Find user: %s", err.Error())
// 		}
// 		assert.Empty(t, err, "Could not delete user")

// 	})

// }
