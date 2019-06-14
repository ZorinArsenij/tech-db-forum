package http

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/forum"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/post"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/service"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/thread"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/user"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/vote"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/middleware"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/usecase"
	"github.com/buaazp/fasthttprouter"
)

type Api struct {
	Router *fasthttprouter.Router
}

func NewRestApi(
	userInteractor *usecase.UserInteractor,
	forumInteractor *usecase.ForumInteractor,
	threadInteractor *usecase.ThreadInteractor,
	postInteractor *usecase.PostInteractor,
	voteInteractor *usecase.VoteInteractor,
	serviceInteractor *usecase.ServiceInteractor,
) *Api {
	router := fasthttprouter.New()

	//User routes
	router.POST("/api/user/:nickname/create", middleware.LoggerMiddleware(user.CreateUser(userInteractor)))
	router.GET("/api/user/:nickname/profile", middleware.LoggerMiddleware(user.GetUserByNickname(userInteractor)))
	router.POST("/api/user/:nickname/profile", middleware.LoggerMiddleware(user.UpdateUser(userInteractor)))

	//Forum routes
	router.POST("/api/forum/:slug", middleware.LoggerMiddleware(forum.CreateForum(forumInteractor)))
	router.GET("/api/forum/:slug/details", middleware.LoggerMiddleware(forum.GetForum(forumInteractor)))
	router.GET("/api/forum/:slug/users", middleware.LoggerMiddleware(forum.GetForumUsers(forumInteractor)))

	//Thread routes
	router.GET("/api/thread/:slug_or_id/details", middleware.LoggerMiddleware(thread.GetThread(threadInteractor)))
	router.GET("/api/forum/:slug/threads", middleware.LoggerMiddleware(thread.GetThreads(threadInteractor)))
	router.POST("/api/forum/:slug/create", middleware.LoggerMiddleware(thread.CreateThread(threadInteractor)))
	router.POST("/api/thread/:slug_or_id/details", middleware.LoggerMiddleware(thread.UpdateThread(threadInteractor)))

	//Post routes
	router.GET("/api/post/:id/details", middleware.LoggerMiddleware(post.GetPost(postInteractor)))
	router.POST("/api/post/:id/details", middleware.LoggerMiddleware(post.UpdatePost(postInteractor)))
	router.POST("/api/thread/:slug_or_id/create", middleware.LoggerMiddleware(post.CreatePosts(postInteractor)))
	router.GET("/api/thread/:slug_or_id/posts", middleware.LoggerMiddleware(post.GetPosts(postInteractor)))

	//Vote routes
	router.POST("/api/thread/:slug_or_id/vote", middleware.LoggerMiddleware(vote.CreateVote(voteInteractor)))

	//Service routes
	router.GET("/api/service/status", middleware.LoggerMiddleware(service.GetStatus(serviceInteractor)))
	router.POST("/api/service/clear", middleware.LoggerMiddleware(service.Clear(serviceInteractor)))

	return &Api{
		Router: router,
	}
}
