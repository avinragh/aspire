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
func (siw *ServerInterfaceWrapper) InstallmentById(w http.ResponseWriter, r *http.Request) {
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

	installment, err := database.FindLoanById(idInt)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(installment); err != nil {
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

func (siw *ServerInterfaceWrapper) AddInstallment(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	installment := &models.Installment{}

	// ------------- Path parameter "id" -------------
	var loanId string

	err := runtime.BindStyledParameter("simple", false, "loanId", chi.URLParam(r, "loanId"), &loanId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	loanIdInt, err := strconv.ParseInt(loanId, 10, 64)
	if err != nil {
		w.WriteHeader((http.StatusBadRequest))
	}

	reqBody, _ := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(reqBody, &installment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
		return
	}
	// if loan.Currency == "" {
	// 	loan.Currency = "USD"
	// }
	installment.State = constants.InstallmentStatusPending

	database := ctx.GetDB()

	installment, err = database.AddInstallment(installment, loanIdInt)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		logger.Println(err)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(installment); err != nil {
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

func (siw *ServerInterfaceWrapper) RepayInstallment(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()
	database := ctx.GetDB()

	repaymentRequest := models.RepaymentRequest{}

	logger.Println("In Find Server By id")

	// ------------- Path parameter "id" -------------
	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	IdInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not Aprrove Loan: %s", err), http.StatusInternalServerError)
		return

	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(reqBody, &repaymentRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
		return
	}

	err = database.RepayInstallment(IdInt, *repaymentRequest.RepaymentAmount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Println(err)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(""); err != nil {
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

func (siw *ServerInterfaceWrapper) FindInstallments(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()

	database := ctx.GetDB()

	logger := ctx.GetLogger()

	installments := []*models.Installment{}

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params models.FindInstallmentsParams

	// ------------- Optional query parameter "userId" -------------

	err = runtime.BindQueryParameter("form", true, false, "loanId", r.URL.Query(), &params.LoanID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter userId: %s", err), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "state", r.URL.Query(), &params.State)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter state: %s", err), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "sort", r.URL.Query(), &params.Sort)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter sort: %s", err), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter limit: %s", err), http.StatusBadRequest)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter page: %s", err), http.StatusBadRequest)
		return
	}

	userRole := r.Context().Value(ContextRoleKey).(string)

	var loanIdInt int64

	if userRole == constants.RoleUser && params.LoanID == "" {
		http.Error(w, fmt.Sprintf("Invalid parameter loanId: %s", err), http.StatusBadRequest)
		return
	}

	if params.LoanID != "" {
		loanIdInt, err = strconv.ParseInt(params.LoanID, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid parameter loanId: %s", err), http.StatusBadRequest)
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

	installments, err = database.FindInstallments(loanIdInt, params.State, params.Sort, limit, page)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter username: %s", err), http.StatusBadRequest)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(installments); err != nil {
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
