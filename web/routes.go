package web

import (
	"aspire/context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(ctx *context.Context, si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	middlewares := []MiddlewareFunc{}

	options.Middlewares = middlewares

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
		r.Post(options.BaseURL+"/v1"+"/Signup", wrapper.Signup)
	})

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1"+"/Login", wrapper.Login)
	})

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1"+"/Loans/{id}", wrapper.IsAuthorized(wrapper.LoanById))
	})

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1"+"/Loans", wrapper.IsAuthorized(wrapper.Loans))
	})

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1"+"/Loans", wrapper.IsAuthorized(wrapper.AddLoan))
	})

	r.Group(func(r chi.Router) {
		r.Patch(options.BaseURL+"/v1"+"/Loans/{id}/Approve", wrapper.IsAuthorized(wrapper.ApproveLoan))
	})

	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/Loans/{id}", wrapper.IsAuthorized(wrapper.DeleteLoan))
	})

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1"+"/Installments/{id}", wrapper.IsAuthorized(wrapper.InstallmentById))
	})

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1"+"/Installments", wrapper.IsAuthorized(wrapper.Installments))
	})

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1"+"/Installments", wrapper.IsAuthorized(wrapper.AddInstallment))
	})

	r.Group(func(r chi.Router) {
		r.Patch(options.BaseURL+"/v1"+"/Installments/{id}/Repay", wrapper.IsAuthorized(wrapper.RepayInstallment))
	})

	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/Installments/{id}", wrapper.IsAuthorized(wrapper.DeleteInstallment))
	})

	return r
}
