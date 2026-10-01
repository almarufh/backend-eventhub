package handler

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/middleware"
	"backend/EventHub/internal/service"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommunityHandler struct {
	cs         *service.CommunitieService
	middleware *middleware.Middleware
}

func NewCommunityHandler(m *middleware.Middleware, cs *service.CommunitieService) *CommunityHandler {
	return &CommunityHandler{
		cs:         cs,
		middleware: m,
	}
}

func (ch *CommunityHandler) JoinCommunity(ctx *gin.Context) {
	payload, err := ch.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[CommunityHandler.Logout] Unauthorized: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	param, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	joined, err := ch.cs.ToggleJoinCommunity(ctx.Request.Context(), int32(param), payload)
	if err != nil {
		log.Printf("[CommunityHandler.Logout] Service error: %s\n", err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: fmt.Sprintf("%s successfully", joined.Status),
		Data:    joined,
	})
}

func (ch *CommunityHandler) GetDetailCommunity(ctx *gin.Context) {
	param, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	communitiy, err := ch.cs.DetailCommunity(ctx.Request.Context(), int32(param))

	if err != nil {
		log.Printf("[CommunityHandler.Detail] Service error: %s\n", err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Succes get detail communitiy",
		Data:    communitiy,
	})

}

func (ch *CommunityHandler) GetMembersCommunity(ctx *gin.Context) {
	param, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	membersCommunities, err := ch.cs.GetMembersCommunity(ctx.Request.Context(), int32(param))

	if err != nil {
		log.Printf("[CommunityHandler.Members] Service error: %s\n", err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Succes get members communitiy",
		Data:    membersCommunities,
	})

}
