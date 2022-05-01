package main

import (
	"aspire/context"
	"aspire/crons"
	"aspire/web"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

func main() {
	ctx := &context.Context{}
	ctx, err := ctx.Init()
	if err != nil {
		log.Fatal("unable to create context")
	}
	defer ctx.GetDB().Close()

	c := cron.New()
	c.AddFunc("@hourly", func() {
		crons.UpdatePaidLoans(ctx)
	})
	c.Start()

	siw := &web.ServerInterfaceWrapper{}
	r := web.Handler(ctx, siw.Handler)
	http.ListenAndServe(":8080", r)

}
