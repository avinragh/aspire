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

// FindLoanByid Handler
func (siw *ServerInterfaceWrapper) LoanById(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err = errors.New("ID not right format")
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, err.Error())
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	database := ctx.GetDB()

	loan, err := database.FindLoanById(idInt)
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
		if err := json.NewEncoder(w).Encode(loan); err != nil {
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

//AddLoan Handler
func (siw *ServerInterfaceWrapper) AddLoan(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()
	logger := ctx.GetLogger()

	loan := &models.Loan{}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return

	}
	err = json.Unmarshal(reqBody, &loan)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	loan.Validate(strfmt.Default)

	userIdString := r.Context().Value(ContextUserIdKey).(string)

	userId, err := strconv.ParseInt(userIdString, 16, 64)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	database := ctx.GetDB()

	loan, err = database.AddLoan(loan, userId)
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
		if err := json.NewEncoder(w).Encode(loan); err != nil {
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

func (siw *ServerInterfaceWrapper) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	ctx := siw.GetContext()
	logger := ctx.GetLogger()
	database := ctx.GetDB()

	logger.Println("In Find Server By id")

	var id string

	err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: id not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	//check user role. Allow only if admin
	userRole := r.Context().Value(ContextRoleKey).(string)

	var handler = func(w http.ResponseWriter, r *http.Request) {}

	if userRole == constants.RoleAdmin {
		loan, err := database.FindLoanById(idInt)
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

		if loan.State == constants.LoanStatusPending {
			err = database.ApproveLoan(idInt)
			if err != nil {
				errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
				logger.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(errorResponse)
				return
			}

			handler = func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(loan); err != nil {
					errorResponse := aerrors.New(aerrors.ErrInternalServerCode, aerrors.ErrInternalServerMessage, "")
					logger.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(errorResponse)
					return

				}
			}
		}

	} else {
		handler = func(w http.ResponseWriter, r *http.Request) {
			err = errors.New("Require Admin role to perform action")
			errorResponse := aerrors.New(aerrors.ErrForbiddenCode, aerrors.ErrForbiddenMessage, err.Error())
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

func (siw *ServerInterfaceWrapper) FindLoans(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()

	database := ctx.GetDB()

	logger := ctx.GetLogger()

	loans := []*models.Loan{}

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params models.FindLoansParams

	// -------------  query parameters -------------

	err = runtime.BindQueryParameter("form", true, false, "userId", r.URL.Query(), &params.UserID)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: userId not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "state", r.URL.Query(), &params.State)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: state nt in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "sort", r.URL.Query(), &params.Sort)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: sort param not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: sort param not in cporrect format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = runtime.BindQueryParameter("form", true, false, "page", r.URL.Query(), &params.Page)
	if err != nil {
		errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: page param not in correct format")
		logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	userRole := r.Context().Value(ContextRoleKey).(string)
	userId := r.Context().Value(ContextUserIdKey).(string)

	var userIdInt int64

	if userRole == constants.RoleAdmin {
		if params.UserID != "" {
			userIdInt, err = strconv.ParseInt(params.UserID, 10, 64)
			if err != nil {
				errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: user Id not in correct format")
				logger.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
		}
	} else if userRole == constants.RoleUser {
		userIdInt, err = strconv.ParseInt(userId, 10, 64)
		if err != nil {
			errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: user Id not in correct format")
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
				errorResponse := aerrors.New(aerrors.ErrInputValidationCode, aerrors.ErrInputValidationMessage, "Bad Request: limit not in correct format")
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

	loans, err = database.FindLoans(userIdInt, params.State, params.Sort, limit, page)
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
		if err := json.NewEncoder(w).Encode(loans); err != nil {
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
