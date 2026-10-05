package router

import (
	"backend/EventHub/internal/handler"
	"backend/EventHub/internal/repo"
	"backend/EventHub/internal/service"
)

func (r *Router) eventRouter() {
	event := r.router.Group("events")

	er := repo.NewEventsRepo(r.pool)
	ar := repo.NewAuthRepo()
	es := service.NewEventsService(er, ar, r.pool)
	eh := handler.NewEventsHandler(es, r.middleware)

	{
		event.GET(":id", eh.GetDetailEvents)
		event.POST(":id/join", r.middleware.AuthMiddleware, eh.JoinEvent)
		event.POST(":id/save", r.middleware.AuthMiddleware, eh.SaveEvent)
		event.DELETE(":id/leave", r.middleware.AuthMiddleware, eh.LeaveEvent)
		event.DELETE(":id/unsave", r.middleware.AuthMiddleware, eh.UnsaveEvent)
	}
}
