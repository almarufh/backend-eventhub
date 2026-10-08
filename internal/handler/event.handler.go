package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"backend/EventHub/internal/dto"
	msgerr "backend/EventHub/internal/message"
	"backend/EventHub/internal/middleware"
	"backend/EventHub/internal/service"

	"github.com/gin-gonic/gin"
)

type EventsHandler struct {
	es         *service.EventsService
	middleware *middleware.Middleware
}

func NewEventsHandler(es *service.EventsService,
	middleware *middleware.Middleware) *EventsHandler {
	return &EventsHandler{
		es:         es,
		middleware: middleware,
	}
}

// GetDetailEvents godoc
// @Summary      Get Detail Event
// @Description  Mengambil data detail event berdasarkan Event ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  dto.Response  "Detail event berhasil didapatkan"
// @Failure      400  {object}  dto.Response  "ID event tidak valid"
// @Failure      404  {object}  dto.Response  "Event tidak ditemukan"
// @Failure      500  {object}  dto.Response  "Internal server error"
// @Router       /events/{id} [get]
func (eh *EventsHandler) GetDetailEvents(ctx *gin.Context) {
	idParam := ctx.Param("id")
	eventID, err := strconv.Atoi(idParam)
	if err != nil || eventID <= 0 {
		log.Printf("[EventsHandler.GetDetailEvents] Invalid ID param: %v\n", idParam)
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID event tidak valid",
		})
		return
	}

	res, err := eh.es.GetDetailEvents(ctx.Request.Context(), int32(eventID))
	if err != nil {
		log.Printf("[EventsHandler.GetDetailEvents] Service error: %v\n", err.Error())

		if err.Error() == "event not found" {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: "Event tidak ditemukan",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal mengambil detail event",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil mengambil detail event",
		Data:    res,
	})
}

// JoinEvent
//
// @Summary      Join Event
// @Description  Join a event based on the provided event ID in path parameter
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Security        Bearer
// @Router       /events/{id}/join [post]
// @Success      200  {object}  dto.Response  "Berhasil bergabung dengan event"
// @Failure      400  {object}  dto.Response  "ID event tidak valid"
// @Failure      401  {object}  dto.Response  "Unauthorized"
// @Failure      404  {object}  dto.Response  "Event tidak ditemukan"
// @Failure      500  {object}  dto.Response  "Internal server error"
func (eh *EventsHandler) JoinEvent(ctx *gin.Context) {
	payload, err := eh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	eventIDParam := ctx.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil || eventID <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID event tidak valid",
		})
		return
	}

	if payload.Role != "attendee" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: fmt.Sprintf("%s cannot join event", payload.Role),
		})
		return
	}

	if err := eh.es.JoinEventService(ctx.Request.Context(), payload.ID, int32(eventID)); err != nil {
		log.Printf("[EventsHandler.JoinEvent] Error: %v\n", err)

		if errors.Is(err, msgerr.AlredyJoinedEvents) {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, msgerr.CapacityFulled) {
			ctx.JSON(http.StatusUnprocessableEntity, dto.ErrResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal bergabung dengan event",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil bergabung dengan event",
	})
}

// SaveEvent
//
// @Summary      Save Event
// @Description  Save a event based on the provided event ID in path parameter
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Security     Bearer
// @Router       /events/{id}/save [post]
// @Success      200  {object}  dto.Response  "Berhasil bergabung dengan event"
// @Failure      400  {object}  dto.Response  "ID event tidak valid"
// @Failure      401  {object}  dto.Response  "Unauthorized"
// @Failure      404  {object}  dto.Response  "Event tidak ditemukan"
// @Failure      500  {object}  dto.Response  "Internal server error"
func (eh *EventsHandler) SaveEvent(ctx *gin.Context) {
	payload, err := eh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	eventIDParam := ctx.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil || eventID <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID event tidak valid",
		})
		return
	}

	if payload.Role != "attendee" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: fmt.Sprintf("%s cannot save event", payload.Role),
		})
		return
	}

	if err := eh.es.SaveEventService(ctx.Request.Context(), payload.ID, int32(eventID)); err != nil {
		log.Printf("[EventsHandler.JoinEvent] Error: %v\n", err)

		if errors.Is(err, msgerr.AlredySavedEvent) {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal menyimpan event",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil menyimpan event",
	})
}

// LeaveEvent
//
// @Summary      Leave Event
// @Description  Leave a event based on the provided event ID in path parameter
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Security     Bearer
// @Router       /events/{id}/leave [delete]
// @Success      200  {object}  dto.Response  "Berhasil bergabung dengan event"
// @Failure      400  {object}  dto.Response  "ID event tidak valid"
// @Failure      401  {object}  dto.Response  "Unauthorized"
// @Failure      404  {object}  dto.Response  "Event tidak ditemukan"
// @Failure      500  {object}  dto.Response  "Internal server error"
func (eh *EventsHandler) LeaveEvent(ctx *gin.Context) {
	payload, err := eh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	eventIDParam := ctx.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil || eventID <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID event tidak valid",
		})
		return
	}

	if payload.Role != "attendee" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: fmt.Sprintf("%s cannot join event", payload.Role),
		})
		return
	}

	if err := eh.es.LeaveEventService(ctx.Request.Context(), payload.ID, int32(eventID)); err != nil {
		log.Printf("[EventsHandler.JoinEvent] Error: %v\n", err)

		if errors.Is(err, msgerr.NotJoinedEvents) {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil batal bergabung dengan event",
	})
}

// UnsaveEvent
//
// @Summary      Unsave Event
// @Description  Unsave a event based on the provided event ID in path parameter
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Security	 Bearer
// @Router       /events/{id}/unsave [delete]
// @Success      200  {object}  dto.Response  "Berhasil bergabung dengan event"
// @Failure      400  {object}  dto.Response  "ID event tidak valid"
// @Failure      401  {object}  dto.Response  "Unauthorized"
// @Failure      404  {object}  dto.Response  "Event tidak ditemukan"
// @Failure      500  {object}  dto.Response  "Internal server error"
func (eh *EventsHandler) UnsaveEvent(ctx *gin.Context) {
	payload, err := eh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	eventIDParam := ctx.Param("id")
	eventID, err := strconv.Atoi(eventIDParam)
	if err != nil || eventID <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "ID event tidak valid",
		})
		return
	}

	if payload.Role != "attendee" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: fmt.Sprintf("%s cannot unsave event", payload.Role),
		})
		return
	}

	if err := eh.es.UnsaveEventService(ctx.Request.Context(), payload.ID, int32(eventID)); err != nil {
		log.Printf("[EventsHandler.JoinEvent] Error: %v\n", err)

		if errors.Is(err, msgerr.NotSavedEvent) {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil menghapus event",
	})
}
