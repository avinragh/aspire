package web

import (
	"fmt"
	"log"
	"net/http"
)

func (siw *ServerInterfaceWrapper) HealthCheck(w http.ResponseWriter, r *http.Request) {

	ctx := siw.GetContext()

	database := ctx.GetDB()

	err := database.Ping()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not connect to DB: %s", err), http.StatusInternalServerError)
	}
	log.Println("Database connection established")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
