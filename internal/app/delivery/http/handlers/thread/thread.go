package thread

import (
	"fmt"

	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/message"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/thread"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase"
	"github.com/jackc/pgx"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func GetThread(interactor *usecase.ThreadInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		slugOrId := ctx.UserValue("slug_or_id").(string)

		thread, err := interactor.GetThread(slugOrId)
		switch err {
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "Forum doesn't exist",
				}
				if _, err = easyjson.MarshalToWriter(msg, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
		case nil:
			{
				if _, err = easyjson.MarshalToWriter(thread, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusOK)
				return
			}
		default:
			{
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		}
	}
}

func GetThreads(interactor *usecase.ThreadInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		var limit *int
		var since *string

		slug := ctx.UserValue("slug").(string)
		if limitRaw := ctx.QueryArgs().GetUintOrZero("limit"); limitRaw != 0 {
			limit = &limitRaw
		}
		if exists := ctx.QueryArgs().Has("since"); exists {
			sinceRaw := string(ctx.QueryArgs().Peek("since"))
			since = &sinceRaw
		}
		orderDesc := ctx.QueryArgs().GetBool("desc")

		threads, err := interactor.GetThreads(slug, limit, since, orderDesc)
		switch err {
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "Forum doesn't exist",
				}
				if _, err = easyjson.MarshalToWriter(msg, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
		case nil:
			{
				if _, err = easyjson.MarshalToWriter(threads, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusOK)
				return
			}
		default:
			{
				fmt.Println(err)
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		}
	}
}

func CreateThread(interactor *usecase.ThreadInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		forumSlug := ctx.UserValue("slug").(string)

		data := &thread.Create{}
		if err := data.UnmarshalJSON(ctx.PostBody()); err != nil || data.Title == "" || data.Message == "" || data.UserNickname == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		data.ForumSlug = forumSlug

		received, err := interactor.CreateThread(data)
		switch err {
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "User or forum doesn't exist",
				}
				_, err = easyjson.MarshalToWriter(msg, ctx.Response.BodyWriter())
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
		case nil:
			{
				_, err = easyjson.MarshalToWriter(received, ctx.Response.BodyWriter())
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusCreated)
				return
			}
		default:
			{
				_, err = easyjson.MarshalToWriter(received, ctx.Response.BodyWriter())
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusConflict)
				return
			}
		}
	}
}

func UpdateThread(interactor *usecase.ThreadInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		var data thread.Update

		slugOrId := ctx.UserValue("slug_or_id").(string)
		err := data.UnmarshalJSON(ctx.PostBody())
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		updated, err := interactor.UpdateThread(&data, slugOrId)
		switch err {
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "Thread doesn't exist",
				}
				if _, err = easyjson.MarshalToWriter(msg, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
		case nil:
			{
				if _, err = easyjson.MarshalToWriter(updated, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusOK)
				return
			}
		default:
			{
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		}
	}
}
