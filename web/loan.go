package web

import (
	"aspire/constants"
	"aspire/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

	server, err := database.FindLoanById(idInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(server); err != nil {
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

	database := ctx.GetDB()

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

	loan, err = database.AddLoan(loan, 1)
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

/*func (siw *ServerInterfaceWrapper) UpdateLoan(w http.ResponseWriter, r *http.Request) {
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

	database := ctx.GetDB()


	if loan.Currency == "" {
		loan.Currency = "USD"
	}
	loan.State = constants.LoanStatusPending

	loan, err = database.AddLoan(loan)
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
}*/
