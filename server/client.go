package server

import (
	"github.com/ZorinArsenij/tech-db-forum/models"
	"github.com/jackc/pgx"
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func (server *Server) CreateUser(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	client := &models.Client{}
	err := client.UnmarshalJSON(ctx.PostBody())
	if err != nil || client.Email == "" || client.Fullname == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	client.Nickname = ctx.UserValue("nickname").(string)

	clients, err := server.Dbm.CreateUser(client)
	if clients == nil && err == nil {
		_, err := easyjson.MarshalToWriter(client, ctx.Response.BodyWriter())
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetStatusCode(fasthttp.StatusCreated)
		return
	}

	_, err = easyjson.MarshalToWriter(clients, ctx.Response.BodyWriter())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusConflict)
}

func (server *Server) GetUser(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	nickname := ctx.UserValue("nickname").(string)

	client, err := server.Dbm.GetUser(nickname)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	_, err = easyjson.MarshalToWriter(client, ctx.Response.BodyWriter())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (server *Server) UpdateUser(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	client := &models.Client{}
	err := client.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	client.Nickname = ctx.UserValue("nickname").(string)

	err = server.Dbm.UpdateUser(client)
	switch {
	case err == pgx.ErrNoRows:
		{
			message := models.RequestError{
				Message: "User not found",
			}
			_, err = easyjson.MarshalToWriter(message, ctx.Response.BodyWriter())
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
	case err != nil:
		{
			message := models.RequestError{
				Message: "User with this email already exists",
			}
			_, err = easyjson.MarshalToWriter(message, ctx.Response.BodyWriter())
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			ctx.SetStatusCode(fasthttp.StatusConflict)
			return
		}
	}

	_, err = easyjson.MarshalToWriter(client, ctx.Response.BodyWriter())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
}
