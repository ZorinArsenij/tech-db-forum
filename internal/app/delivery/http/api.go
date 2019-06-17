package http

import (
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/forum"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/post"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/service"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/thread"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/user"
	"github.com/ZorinArsenij/tech-db-forum/internal/app/delivery/http/handlers/vote"
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
	router.POST("/api/user/:nickname/create", user.CreateUser(userInteractor))
	router.GET("/api/user/:nickname/profile", user.GetUserByNickname(userInteractor))
	router.POST("/api/user/:nickname/profile", user.UpdateUser(userInteractor))

	//Forum routes
	router.POST("/api/forum/:slug", forum.CreateForum(forumInteractor))
	router.GET("/api/forum/:slug/details", forum.GetForum(forumInteractor))
	router.GET("/api/forum/:slug/users", forum.GetForumUsers(forumInteractor))

	//Thread routes
	router.GET("/api/thread/:slug_or_id/details", thread.GetThread(threadInteractor))
	router.GET("/api/forum/:slug/threads", thread.GetThreads(threadInteractor))
	router.POST("/api/forum/:slug/create", thread.CreateThread(threadInteractor))
	router.POST("/api/thread/:slug_or_id/details", thread.UpdateThread(threadInteractor))

	//Post routes
	router.GET("/api/post/:id/details", post.GetPost(postInteractor))
	router.POST("/api/post/:id/details", post.UpdatePost(postInteractor))
	router.POST("/api/thread/:slug_or_id/create", post.CreatePosts(postInteractor))
	router.GET("/api/thread/:slug_or_id/posts", post.GetPosts(postInteractor))

	//Vote routes
	router.POST("/api/thread/:slug_or_id/vote", vote.CreateVote(voteInteractor))

	//Service routes
	router.GET("/api/service/status", service.GetStatus(serviceInteractor))
	router.POST("/api/service/clear", service.Clear(serviceInteractor))

	return &Api{
		Router: router,
	}
}
