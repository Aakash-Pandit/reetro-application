package routes

import (
	"log"
	"net/http"

	"github.com/Aakash-Pandit/reetro-golang/common"
	"github.com/Aakash-Pandit/reetro-golang/core"
	"github.com/Aakash-Pandit/reetro-golang/middlewares"
	"github.com/Aakash-Pandit/reetro-golang/services"
	"github.com/Aakash-Pandit/reetro-golang/storages"
	"github.com/gorilla/mux"
)

type Router struct {
	Route           *mux.Router
	Port            string
	BasicService    services.BasicService
	UserService     services.UserService
	BoardService    services.BoardService
	FeedbackService services.FeedbackService
	Middleware      middlewares.Middleware
}

func NewRouter(route *mux.Router, port string, db *storages.PostgresStore, client *storages.RedisStore, emailInterface common.EmailInterface) *Router {
	requestUser := services.RequestUser{
		Store: services.ServiceStorage(db),
	}

	return &Router{
		Route: route,
		Port:  port,
		BasicService: services.BasicService{
			RedisClient: client,
		},
		UserService: services.UserService{
			Store:       storages.Storage(db),
			User:        requestUser,
			RedisClient: client,
			Email:       emailInterface,
		},
		BoardService: services.BoardService{
			Store:       storages.Storage(db),
			User:        requestUser,
			RedisClient: client,
		},
		FeedbackService: services.FeedbackService{
			Store:       storages.Storage(db),
			User:        requestUser,
			RedisClient: client,
		},
		Middleware: middlewares.Middleware{
			Store: middlewares.MiddlewareInterface(db),
		},
	}
}

func (r *Router) Run() {
	r.Route.HandleFunc(
		"/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(services.HomeHandler),
		),
	).Methods(http.MethodGet)

	r.Route.HandleFunc(
		"/about/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(services.AboutHandler),
		),
	).Methods(http.MethodGet)

	r.Route.HandleFunc(
		"/clear_redis/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.BasicService.ClearRedisCache),
		),
	).Methods(http.MethodPost)

	r.Route.HandleFunc(
		"/login/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.UserService.LoginHandler),
		),
	).Methods(http.MethodPost)

	r.Route.HandleFunc(
		"/users/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.UserService.GetAllUsersHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodGet)

	r.Route.HandleFunc(
		"/users/{id}/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.UserService.GetUserByIdHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodGet)

	r.Route.HandleFunc(
		"/users/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.UserService.CreateUserHandler),
		),
	).Methods(http.MethodPost)

	r.Route.HandleFunc(
		"/signup/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.UserService.CreateUserHandler),
		),
	).Methods(http.MethodPost)

	r.Route.HandleFunc(
		"/users/{id}/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.UserService.DeleteUserHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodDelete)

	r.Route.HandleFunc(
		"/forgot_password/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.UserService.ForgotPasswordHandler),
		),
	).Methods(http.MethodPost)

	r.Route.HandleFunc(
		"/reset_password/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.UserService.ResetPasswordHandler),
		),
	).Methods(http.MethodPost)

	r.Route.HandleFunc(
		"/boards/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.BoardService.GetAllBoardsHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodGet)

	r.Route.HandleFunc(
		"/boards/{id}/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.BoardService.GetBoardByIdHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodGet)

	r.Route.HandleFunc(
		"/boards/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.BoardService.CreateBoardHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodPost)

	r.Route.HandleFunc(
		"/boards/{id}/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.BoardService.UpdateBoardHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodPatch)

	r.Route.HandleFunc(
		"/boards/{id}/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.BoardService.DeleteBoardHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodDelete)

	r.Route.HandleFunc(
		"/feedbacks/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.FeedbackService.GetAllFeedbacksHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodGet)

	r.Route.HandleFunc(
		"/feedbacks/{id}/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.FeedbackService.GetFeedbackByIdHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodGet)

	r.Route.HandleFunc(
		"/feedbacks/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.FeedbackService.CreateFeedbackHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodPost)

	r.Route.HandleFunc(
		"/feedbacks/{id}/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.FeedbackService.UpdateFeedbackHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodPatch)

	r.Route.HandleFunc(
		"/feedbacks/{id}/",
		middlewares.ChainOfMiddleware(
			core.HTTPHandleFunc(r.FeedbackService.DeleteFeedbackHandler),
			r.Middleware.JWTAuthentication,
		),
	).Methods(http.MethodDelete)

	http.Handle("/", r.Route)
	log.Fatal(http.ListenAndServe(r.Port, nil))
}
