package vote

import (
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/message"
	"github.com/ZorinArsenij/tech-db-forum/pkg/domain/vote"
	"github.com/ZorinArsenij/tech-db-forum/pkg/usecase"
	"github.com/jackc/pgx"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func CreateVote(interactor *usecase.VoteInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		slugOrId := ctx.UserValue("slug_or_id").(string)

		data := &vote.Vote{}
		err := data.UnmarshalJSON(ctx.PostBody())
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		switch data.Rating {
		case 1:
			data.Voice = true
		case -1:
			data.Voice = false
		default:
			{
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
		}

		thread, err := interactor.CreateVote(data, slugOrId)
		switch err {
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "User or thread doesn't exist",
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
