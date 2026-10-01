package router

import (
	"backend/EventHub/internal/handler"
	"backend/EventHub/internal/repo"
	"backend/EventHub/internal/service"
)

func (r *Router) authRouter() {
	auth := r.router.Group("/auth")

	ar := repo.NewAuthRepo(r.pool)
	as := service.NewAuthService(ar)
	ah := handler.NewAuthHandler(as, r.middleware)
	r.middleware.RepoAuthMiddleWare(ar)

	auth.GET("/cektoken", r.middleware.CekDong)

	{
		auth.POST("/register", ah.Register)
		auth.POST("/login", ah.Login)
		auth.DELETE("/logout", r.middleware.AuthMiddleware, ah.Logout)
		auth.PATCH("/forgot/password", ah.ChangePassword)
	}
}
