package router

import (
	"backend/EventHub/internal/handler"
	"backend/EventHub/internal/repo"
	"backend/EventHub/internal/service"
)

func (r *Router) communitieRouter() {
	com := r.router.Group("communities")

	cr := repo.NewCommunityRepo(r.pool)
	cs := service.NewCommunitieService(r.pool, cr)
	ch := handler.NewCommunityHandler(r.middleware, cs)

	{
		com.GET(":id", ch.GetDetailCommunity)
		com.GET(":id/members", ch.GetMembersCommunity)
		com.POST(":id/join", ch.JoinCommunity)
		com.DELETE(":id/leave", ch.JoinCommunity)
	}

}
