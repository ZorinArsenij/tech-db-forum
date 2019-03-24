package forum

import (
	"errors"

	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/forum"
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/message"
	"github.com/ZorinArsenij/tech-db-forum/pkg/usecase"
	"github.com/jackc/pgx"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func GetForum(interactor *usecase.ForumInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		slug := ctx.UserValue("slug").(string)

		received, err := interactor.GetForum(slug)
		switch err {
		case nil:
			{
				if _, err = easyjson.MarshalToWriter(received, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusOK)
				return
			}
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
		default:
			{
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		}
	}
}

func CreateForum(interactor *usecase.ForumInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		data := &forum.Create{}
		if err := data.UnmarshalJSON(ctx.PostBody()); err != nil || data.Title == "" || data.Slug == "" || data.UserNickname == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		created, err := interactor.CreateForum(data)
		switch err {
		case nil:
			{
				if _, err := easyjson.MarshalToWriter(created, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusCreated)
				return
			}
		case errors.New("threadAlreadyExists"):
			{
				if _, err = easyjson.MarshalToWriter(created, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusConflict)
				return
			}
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "User doesn't exist",
				}
				if _, err = easyjson.MarshalToWriter(msg, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusNotFound)
				return
			}
		default:
			{
				if _, err = easyjson.MarshalToWriter(created, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusConflict)
				return
			}
		}
	}
}

func GetForumUsers(interactor *usecase.ForumInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		slug := ctx.UserValue("slug").(string)

		var limit *int
		if limitRaw := ctx.QueryArgs().GetUintOrZero("limit"); limitRaw != 0 {
			limit = &limitRaw
		}

		var since *string
		if exists := ctx.QueryArgs().Has("since"); exists {
			sinceRaw := string(ctx.QueryArgs().Peek("since"))
			since = &sinceRaw
		}

		orderDesc := ctx.QueryArgs().GetBool("desc")

		users, err := interactor.GetForumUsers(slug, limit, since, orderDesc)
		switch err {
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "Forum doesn't exist",
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
				_, err = easyjson.MarshalToWriter(users, ctx.Response.BodyWriter())
				if err != nil {
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
