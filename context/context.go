package context

import (
	"aspire/db"
	"log"
	"os"
)

type Context struct {
	DB     *db.DB
	Logger *log.Logger
}

func (ctx *Context) Init() *Context {
	var err error
	database := &db.DB{}
	database, err = db.Init()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	ctx.DB = database
	file, err := os.OpenFile("aspire.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Openfile error %s", err)
	}
	// messageQueue := &amqp.AMQP{}
	// messageQueue, err = amqp.Init()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx.AMQP = messageQueue
	logger := &log.Logger{}
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ctx.Logger = logger
	return ctx
}

func (ctx *Context) GetDB() (db *db.DB) {
	db = ctx.DB
	return
}

func (ctx *Context) GetLogger() (logger *log.Logger) {
	logger = ctx.Logger
	return
}
