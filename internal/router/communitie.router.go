package router

import (
	"backend/EventHub/internal/handler"
	"backend/EventHub/internal/repo"
	"backend/EventHub/internal/service"
)

func (r *Router) communitieRouter() {
	com := r.router.Group("communities")

	cr := repo.NewCommunityRepo(r.pool)
	ar := repo.NewAuthRepo()
	cs := service.NewCommunitieService(r.pool, cr, ar)
	ch := handler.NewCommunityHandler(r.middleware, cs)

	{
		com.GET(":id", ch.GetDetailCommunity)
		com.GET(":id/members", ch.GetMembersCommunity)
	}

	auth := com.Group("")
	auth.Use(r.middleware.AuthMiddleware)

	{
		auth.POST(":id/join", r.middleware.AuthMiddleware, ch.JoinCommunity)
		auth.DELETE(":id/leave", r.middleware.AuthMiddleware, ch.LeaveCommunity)
	}

}
