package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"aspire/context"

	packContext "context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// (GET /Loans)
	LoanById(w http.ResponseWriter, r *http.Request)

	// (POST /Loans)
	AddLoan(w http.ResponseWriter, r *http.Request)

	// (PUT /Loans/{id})
	UpdateLoan(w http.ResponseWriter, r *http.Request)

	// (POST /Signup)
	Signup(w http.ResponseWriter, r *http.Request)

	// (POST /Login)
	Login(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	Context            *context.Context
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

func (siw *ServerInterfaceWrapper) GetContext() *context.Context {
	return siw.Context
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(ctx *context.Context, si ServerInterface) http.Handler {
	return HandlerWithOptions(ctx, si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL     string
	BaseRouter  chi.Router
	Middlewares []MiddlewareFunc
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(ctx *context.Context, si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(ctx, si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(ctx *context.Context, si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(ctx, si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			err := errors.New("No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte("secretkey")

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			err := errors.New("Your Token has been expired")
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := claims["userId"]
			role := claims["role"]

			ctx := packContext.WithValue(r.Context(), "userId", userId)
			ctx = packContext.WithValue(ctx, "role", role)

			r = r.WithContext(ctx)

			handler.ServeHTTP(w, r)
		}
		err = errors.New("Not Authorized")
		json.NewEncoder(w).Encode(err)
	}
}
