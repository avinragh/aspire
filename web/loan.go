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

		if loan.State == constants.LoanStatusPending {
			installments := util.GetInstallments(idInt, *loan.Amount, *loan.Term)
			err = database.ApproveLoan(idInt, installments)
			if err != nil {
				http.Error(w, fmt.Sprintf("Could not Aprrove Loan: %s", err), http.StatusInternalServerError)
				return
			}

			handler = func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(loan); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					logger.Println(err)
					return

				}
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

func (siw *ServerInterfaceWrapper) FindLoans(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()

	database := ctx.GetDB()

	logger := ctx.GetLogger()

	loans := []*models.Loan{}

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params models.FindLoansParams

	// ------------- Optional query parameter "userId" -------------
	if paramValue := r.URL.Query().Get("userId"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	err = runtime.BindQueryParameter("form", true, false, "userId", r.URL.Query(), &params.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter userId: %s", err), http.StatusBadRequest)
		return
	}

	if paramValue := r.URL.Query().Get("state"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	err = runtime.BindQueryParameter("form", true, false, "state", r.URL.Query(), &params.State)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter state: %s", err), http.StatusBadRequest)
		return
	}

	if paramValue := r.URL.Query().Get("sort"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	err = runtime.BindQueryParameter("form", true, false, "sort", r.URL.Query(), &params.Sort)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter sort: %s", err), http.StatusBadRequest)
		return
	}

	if paramValue := r.URL.Query().Get("limit"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter limit: %s", err), http.StatusBadRequest)
		return
	}

	if paramValue := r.URL.Query().Get("page"); paramValue != "" {
		logger.Printf("%v", paramValue)
	}

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter page: %s", err), http.StatusBadRequest)
		return
	}

	userRole := r.Context().Value(ContextRoleKey).(string)
	userId := r.Context().Value(ContextUserIdKey).(string)

	var userIdInt int64

	if userRole == constants.RoleAdmin {
		if params.UserID != "" {
			userIdInt, err = strconv.ParseInt(params.UserID, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid parameter UserId: %s", err), http.StatusBadRequest)
				return
			}
		}
	} else if userRole == constants.RoleUser {
		userIdInt, err = strconv.ParseInt(userId, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot determine userId: %s", err), http.StatusInternalServerError)
			return
		}
	}

	var limit, page int64

	if params.Limit != "" {
		if params.Page != "" {
			limit, err = strconv.ParseInt(params.Limit, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid parameter limit: %s", err), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Cannot have limit wothout page", http.StatusBadRequest)
			return
		}
		page, err = strconv.ParseInt(params.Page, 10, 64)
	}

	loans, err = database.FindLoans(userIdInt, params.State, params.Sort, limit, page)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter username: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(loans); err != nil {
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
