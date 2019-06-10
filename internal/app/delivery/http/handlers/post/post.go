package post

import (
	"strings"

	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/message"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/post"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase"

	"github.com/jackc/pgx"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func GetPost(interactor *usecase.PostInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		id := ctx.UserValue("id").(string)
		relatedRaw := strings.Split(string(ctx.QueryArgs().Peek("related")), ",")

		related := make(map[string]bool)
		for _, elem := range relatedRaw {
			related[elem] = true
		}

		info, err := interactor.GetPost(id, related)
		switch err {
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "Post doesn't exist",
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
				if _, err = easyjson.MarshalToWriter(info, ctx.Response.BodyWriter()); err != nil {
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

func UpdatePost(interactor *usecase.PostInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		id := ctx.UserValue("id").(string)

		data := &post.Update{}
		err := data.UnmarshalJSON(ctx.PostBody())
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		data.ID = id

		updated, err := interactor.UpdatePost(data)
		switch err {
		case pgx.ErrNoRows:
			{
				msg := message.Message{
					Description: "Post doesn't exist",
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

func CreatePosts(interactor *usecase.PostInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		slugOrId := ctx.UserValue("slug_or_id").(string)

		newPosts := &post.PostsCreate{}
		err := newPosts.UnmarshalJSON(ctx.PostBody())
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		posts, err := interactor.CreatePosts(newPosts, slugOrId)
		switch {
		case err == nil:
			{
				if _, err := easyjson.MarshalToWriter(posts, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusCreated)
				return
			}
		case err == pgx.ErrNoRows:
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
		case err.Error() == "postParentDoesNotExist":
			{
				msg := message.Message{
					Description: "Post parent doesn't exist",
				}
				if _, err = easyjson.MarshalToWriter(msg, ctx.Response.BodyWriter()); err != nil {
					ctx.SetStatusCode(fasthttp.StatusInternalServerError)
					return
				}

				ctx.SetStatusCode(fasthttp.StatusConflict)
				return
			}
		default:
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
	}
}

func GetPosts(interactor *usecase.PostInteractor) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")

		slugOrId := ctx.UserValue("slug_or_id").(string)

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

		sort := string(ctx.QueryArgs().Peek("sort"))

		var posts *post.Posts
		var err error
		switch sort {
		case "flat":
			{
				posts, err = interactor.GetPostsFlat(slugOrId, limit, since, orderDesc)
			}
		case "tree":
			{
				posts, err = interactor.GetPostsTree(slugOrId, limit, since, orderDesc)
			}
		case "parent_tree":
			{
				posts, err = interactor.GetPostsParentTree(slugOrId, limit, since, orderDesc)
			}
		default:
			{
				posts, err = interactor.GetPosts(slugOrId, limit, since, orderDesc)
			}
		}

		switch {
		case err == pgx.ErrNoRows:
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
		case err != nil:
			{
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		}

		if posts == nil {
			posts = &post.Posts{}
		}

		if _, err = easyjson.MarshalToWriter(posts, ctx.Response.BodyWriter()); err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}
