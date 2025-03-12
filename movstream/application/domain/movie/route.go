package movie

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"simple-go/application/domain"
)

type RouterHttp struct {
	router     *gin.RouterGroup
	handler    handler
	middleware domain.Middleware
}

func NewRouterHttp(router *gin.RouterGroup, db *gorm.DB, middle domain.Middleware) domain.HttpHandler {
	repository := NewRepository(db)
	service := NewService(repository)

	handler := NewHandler(&service)

	return &RouterHttp{
		router:     router,
		handler:    handler,
		middleware: middle,
	}
}

func (r RouterHttp) RegisterRoute() {
	r.router.GET("/movies", r.middleware.GetSessionCustomer(), r.handler.GetAllMovies)
	r.router.PATCH("/movies/watch/:movieId", r.middleware.GetSessionCustomer(), r.handler.MovieWatch)
	r.router.PATCH("/movies/vote/:movieId", r.middleware.GetSessionCustomer(), r.handler.MovieVote)
	r.router.PATCH("/movies/unvote/:movieId", r.middleware.GetSessionCustomer(), r.handler.MovieUnVote)

	r.router.POST("/admin/movies", r.middleware.GetSessionCustomer(), r.handler.AdminCreateMovie)
	r.router.PATCH("/admin/movies/:movieId", r.middleware.GetSessionCustomer(), r.handler.AdminUpdateMovie)
	r.router.DELETE("/admin/movies/:movieId", r.middleware.GetSessionCustomer(), r.handler.AdminDeleteMovie)

}
