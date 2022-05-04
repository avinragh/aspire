package web

import (
	"aspire/constants"
	aerrors "aspire/errors"
	"aspire/models"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/strfmt"
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
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: Id not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	database := ctx.GetDB()

	installment, err := database.FindLoanById(idInt)
	if err != nil {
		if err == sql.ErrNoRows {
			errorResponse := aerrors.New(aerrors.ErrNotFoundCode, aerrors.ErrNotFoundMessage, "")
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		} else {
			errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}
	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(installment); err != nil {
			errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
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
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	loanIdInt, err := strconv.ParseInt(loanId, 10, 64)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: Id not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = json.Unmarshal(reqBody, &installment)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = installment.Validate(strfmt.Default)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, err.Error())
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return

	}

	database := ctx.GetDB()

	installment, err = database.AddInstallment(installment, loanIdInt)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(installment); err != nil {
			errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
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
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	IdInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: id not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return

	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	err = json.Unmarshal(reqBody, &repaymentRequest)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = database.RepayInstallment(IdInt, *repaymentRequest.RepaymentAmount)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(""); err != nil {
			errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
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
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: param loanId not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "state", r.URL.Query(), &params.State)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: param state not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "sort", r.URL.Query(), &params.Sort)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: param sort not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: param loanId not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: param loanId not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	userRole := r.Context().Value(ContextRoleKey).(string)

	var loanIdInt int64

	if userRole == constants.RoleUser && params.LoanID == "" {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: invalid param loanId")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if params.LoanID != "" {
		loanIdInt, err = strconv.ParseInt(params.LoanID, 10, 64)
		if err != nil {
			errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: invalid param loanId")
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	var limit, page int64

	if params.Limit != "" {
		if params.Page != "" {
			limit, err = strconv.ParseInt(params.Limit, 10, 64)
			if err != nil {
				errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: param limit malformed")
				logger.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
		} else {
			err = errors.New("Cannot have limit wothout page")
			errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, err.Error())
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		page, err = strconv.ParseInt(params.Page, 10, 64)
	}

	installments, err = database.FindInstallments(loanIdInt, params.State, params.Sort, limit, page)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(installments); err != nil {
			errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(r.Context()))
}
