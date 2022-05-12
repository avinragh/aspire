package web

import (
	aerrors "aspire/errors"
	"encoding/json"
	"errors"
	"net/http"
)

func (siw *ServerInterfaceWrapper) HealthCheck(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()

	database := ctx.GetDB()

	logger := ctx.GetLogger()

	err := database.Ping()
	if err != nil {
		err := errors.New("Could not onnect to Database")
		errorResponse := aerrors.New(aerrors.ErrForbiddenCode, aerrors.ErrForbiddenMessage, err.Error())
		logger.Println(err)
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	logger.Println("Database connection established")

	w.WriteHeader(http.StatusOK)
}
