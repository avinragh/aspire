package main

import (
	"aspire/context"
	"aspire/web"
	"net/http"
)

func main() {
	ctx := &context.Context{}
	ctx = ctx.Init()
	siw := &web.ServerInterfaceWrapper{}
	r := web.Handler(ctx, siw.Handler)
	http.ListenAndServe(":8080", r)

}
