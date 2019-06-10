package service

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func GetStatus(interactor *usecase.ServiceInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		status, _ := interactor.GetStatus()
		if _, err := easyjson.MarshalToWriter(status, ctx.Response.BodyWriter()); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}

func Clear(interactor *usecase.ServiceInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		if err := interactor.Clear(); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}
