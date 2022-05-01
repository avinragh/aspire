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
		r.Post(options.BaseURL+"/v1"+"/Loans", IsAuthorized(wrapper.AddLoan))
	})
	// r.Group(func(r chi.Router) {
	// 	r.Delete(options.BaseURL+"/Loans/{id}", wrapper.DeleteLoan)
	// })
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1"+"/Loans/{id}", wrapper.LoanById)
	})
	// r.Group(func(r chi.Router) {
	// 	r.Post(options.BaseURL+"/Loans/{id}", wrapper.AddLoanById)
	// })
	return r
}
