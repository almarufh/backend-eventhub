package router

import (
	"backend/EventHub/internal/handler"
	"backend/EventHub/internal/repo"
	"backend/EventHub/internal/service"
)

func (r *Router) userRouter() {
	user := r.router.Group("user")
	user.Use(r.middleware.AuthMiddleware)
	set := user.Group("set")

	ur := repo.NewUserRepo(r.pool)
	ar := repo.NewAuthRepo(r.pool)
	as := service.NewAuthService(ar)
	us := service.NewUserService(r.pool, ur, ar, as, r.redis)
	uh := handler.NewUserHandler(r.middleware, us)

	{
		user.GET("profiles", uh.MyProfile)
	}

	{
		set.PATCH("password", uh.NewPassword)
		set.PATCH("profiles", uh.SetProfile)
	}
}
