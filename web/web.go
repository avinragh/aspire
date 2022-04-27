// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.8.2 DO NOT EDIT.
package web

import (
	"net/http"

	"aspire/context"
	"github.com/go-chi/chi/v5"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// (GET /Loans)
	FindLoans(w http.ResponseWriter, r *http.Request, params string)

	// (POST /Loans)
	AddLoans(w http.ResponseWriter, r *http.Request)

	// (DELETE /Loans/{id})
	DeleteLoan(w http.ResponseWriter, r *http.Request, id string)

	// (GET /Loans/{id})
	FindLoanById(w http.ResponseWriter, r *http.Request, id int32)

	// (POST /Accounts/{id})
	AddLoanById(w http.ResponseWriter, r *http.Request, id string)
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

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(ctx *context.Context, si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	wrapper := &ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		Context:            ctx,
	}

	// r.Group(func(r chi.Router) {
	// 	r.Get(options.BaseURL+"/Loans", wrapper.FindLoans)
	// })
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/Loans", wrapper.AddLoans)
	})
	// r.Group(func(r chi.Router) {
	// 	r.Delete(options.BaseURL+"/Loans/{id}", wrapper.DeleteLoan)
	// })
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/Loans/{id}", wrapper.LoanById)
	})
	// r.Group(func(r chi.Router) {
	// 	r.Post(options.BaseURL+"/Loans/{id}", wrapper.AddLoanById)
	// })
	return r
}
