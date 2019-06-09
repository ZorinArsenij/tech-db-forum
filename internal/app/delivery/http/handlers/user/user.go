package user

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/message"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/user"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase"
	"github.com/jackc/pgx"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func GetUserByNickname(interactor *usecase.UserInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		nickname := ctx.UserValue("nickname").(string)

		received, err := interactor.GetUserByNickname(nickname)
		switch err {
		case nil:
			{
				if _, err := easyjson.MarshalToWriter(received, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusOK)
			}
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "User doesn't exist",
				}
				_, err = easyjson.MarshalToWriter(msg, ctx.Response.BodyWriter())
				if err != nil {
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

func UpdateUser(interactor *usecase.UserInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		data := &user.Update{}
		if err := data.UnmarshalJSON(ctx.PostBody()); err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		nickname := ctx.UserValue("nickname").(string)

		updated, err := interactor.UpdateUser(data, nickname)
		switch err {
		case nil:
			{
				if _, err := easyjson.MarshalToWriter(updated, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusOK)
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
				msg := message.Message{
					Description: "User with this email already exists",
				}
				if _, err = easyjson.MarshalToWriter(msg, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusConflict)
				return
			}
		}

	}
}

func CreateUser(interactor *usecase.UserInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		data := &user.User{}
		if err := data.UnmarshalJSON(ctx.PostBody()); err != nil || data.Email == "" || data.Fullname == "" {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		data.Nickname = ctx.UserValue("nickname").(string)

		users, err := interactor.CreateUser(data)
		switch {
		case err == nil:
			{

				if _, err := easyjson.MarshalToWriter((*users)[0], ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusCreated)
				return
			}
		case err.Error() == "conflict":
			{
				if _, err := easyjson.MarshalToWriter(users, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}
				ctx.SetStatusCode(fasthttp.StatusConflict)
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
