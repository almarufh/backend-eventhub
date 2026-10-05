package router

import (
	"backend/EventHub/internal/handler"
	"backend/EventHub/internal/repo"
	"backend/EventHub/internal/service"
)

func (r *Router) userRouter() {
	ur := repo.NewUserRepo(r.pool)
	er := repo.NewEventsRepo(r.pool)
	ar := repo.NewAuthRepo()
	as := service.NewAuthService(ar, r.redis, r.pool)
	us := service.NewUserService(r.pool, ur, ar, as, r.redis, er)
	uh := handler.NewUserHandler(r.middleware, us)

	user := r.router.Group("user")
	user.Use(r.middleware.AuthMiddleware)
	{
		user.GET("profiles", r.middleware.AuthMiddleware, uh.MyProfile)

	}

	events := user.Group("events")
	{
		events.GET("joined", r.middleware.AuthMiddleware, uh.JoinedEvents)
		events.GET("saved", r.middleware.AuthMiddleware, uh.SavedEvents)
	}

	community := user.Group("community")
	{
		community.GET("joined", r.middleware.AuthMiddleware, uh.JoinedCommunities)
	}

	set := user.Group("set")
	{
		set.PATCH("password", r.middleware.AuthMiddleware, uh.NewPassword)
		set.PATCH("profiles", r.middleware.AuthMiddleware, uh.SetProfile)
	}
}
