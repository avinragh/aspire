package web

import (
	"aspire/constants"
	aerrors "aspire/errors"
	"aspire/models"
	"aspire/util"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/davecgh/go-spew/spew"
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

	installment, err := database.FindInstallmentById(idInt)
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

	userRole := r.Context().Value(ContextRoleKey).(string)

	if userRole == constants.RoleUser {
		err := errors.New("User does not have permission to this method")
		errorResponse := aerrors.New(aerrors.ErrForbiddenCode, aerrors.ErrForbiddenMessage, err.Error())
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

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

	repaymentRequest := &models.RepaymentRequest{}
	installment := &models.Installment{}

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

	userRole := r.Context().Value(ContextRoleKey).(string)
	userId := r.Context().Value(ContextUserIdKey).(string)

	loanIdsSlice := []int64{}
	if userRole == constants.RoleUser {
		userIdInt, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: user Id not in correct format")
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		loanIdsSlice, err = database.FindLoanIdsForUser(userIdInt)
		if err != nil {
			errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		installment, err = database.FindInstallmentById(IdInt)
		if err != nil {
			errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		if !util.Contains(loanIdsSlice, installment.LoanID) {
			err = errors.New("User does not have access to the resource")
			errorResponse := aerrors.New(aerrors.ErrForbiddenCode, aerrors.ErrForbiddenMessage, err.Error())
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return

		}

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

	spew.Dump(repaymentRequest)

	if repaymentRequest.RepaymentAmount == nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: No repaymentAmount in request body	")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	installment, err = database.FindInstallmentById(IdInt)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if installment.State != constants.InstallmentStatusPaid {
		installment, err = database.RepayInstallment(IdInt, *repaymentRequest.RepaymentAmount, loanIdsSlice)
		if err != nil {
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

func (siw *ServerInterfaceWrapper) Installments(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()

	database := ctx.GetDB()

	logger := ctx.GetLogger()

	installments := []*models.Installment{}

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params FindInstallmentsParams

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
	userId := r.Context().Value(ContextUserIdKey).(string)

	loanIdsSlice := []int64{}

	//If user role is user, then get list of accesible loanIds for the user
	if userRole == constants.RoleUser {
		userIdInt, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: user Id not in correct format")
			logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		loanIdsSlice, err = database.FindLoanIdsForUser(userIdInt)
		if err != nil {
			errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
			logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

	}

	if params.Limit != nil && params.Page == nil {
		page := int64(1)
		params.Page = &page

	}

	installments, err = database.FindInstallments(params.LoanID, params.State, params.Sort, params.Limit, params.Page, loanIdsSlice)
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
