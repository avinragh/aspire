package web

import (
	"aspire/constants"
	"aspire/models"
	"aspire/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
)

// FindLoanByid operation middleware
func (siw *ServerInterfaceWrapper) LoanById(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	logger.Println("In Find Server By id")

	// ------------- Path parameter "id" -------------
	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader((http.StatusBadRequest))
	}

	database := ctx.GetDB()

	loan, err := database.FindLoanById(idInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(loan); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}

	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(r.Context()))
}

func (siw *ServerInterfaceWrapper) AddLoan(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	loan := &models.Loan{}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &loan)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
		return
	}
	if loan.Currency == "" {
		loan.Currency = "USD"
	}
	loan.State = constants.LoanStatusPending

	userIdString := r.Context().Value(ContextUserIdKey).(string)

	userId, err := strconv.ParseInt(userIdString, 16, 64)
	database := ctx.GetDB()

	spew.Dump("Hey")
	spew.Dump(userIdString)
	spew.Dump(userId)

	loan, err = database.AddLoan(loan, userId)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		logger.Println(err)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(loan); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			logger.Println(err)
			return
		}

	}
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)

	}

	handler(w, r.WithContext(r.Context()))
}

func (siw *ServerInterfaceWrapper) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()
	database := ctx.GetDB()

	logger.Println("In Find Server By id")

	// ------------- Path parameter "id" -------------
	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)

	//check user role. Allow only if admin
	userRole := r.Context().Value(ContextRoleKey).(string)

	var handler = func(w http.ResponseWriter, r *http.Request) {}

	if userRole == constants.RoleAdmin {
		loan, err := database.FindLoanById(idInt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Loan not found:"), http.StatusBadRequest)
			return
		}

		installments := util.GetInstallments(idInt, *loan.Amount, *loan.Term)
		database.ApproveLoan(idInt, installments)

		handler = func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(loan); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				logger.Println(err)
				return

			}
		}

	} else {
		handler = func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(r.Context()))

}
