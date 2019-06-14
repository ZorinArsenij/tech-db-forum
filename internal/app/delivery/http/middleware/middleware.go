package middleware

import (
	"github.com/valyala/fasthttp"
	"log"
)

func LoggerMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		log.Print("[RECEIVED]", string(ctx.Path()))
		next(ctx)
		log.Println("[FINISHED]", string(ctx.Path()))
	}
}
