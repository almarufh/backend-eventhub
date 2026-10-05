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

// JoinCommunity handles joining a community
//
// @Summary         Join Community
// @Description     Join a community based on the provided community ID in path parameter
// @Tags            Community
// @Produce         json
// @Param           id      path        int     true   "Community ID"                   default(1)
// @Security        Bearer
// @Router          /communities/{id}/join [post]
// @Success         200     {object}    dto.Response
// @Failure         400     {object}    dto.ErrResponse
// @Failure         401     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (ch *CommunityHandler) JoinCommunity(ctx *gin.Context) {
	payload, err := ch.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[CommunityHandler.JoinCommunity] GetPayload: %s\n", err.Error())
		ctx.JSON(http.StatusUnauthorized, dto.ErrResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	param, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("[CommunityHandler.JoinCommunity] Atoi: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrResponse{
			Success: false,
			Message: "Invalid community ID parameter",
		})
		return
	}

	joined, err := ch.cs.ToggleJoinCommunity(ctx.Request.Context(), int32(param), payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
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

// LeaveCommunity handles leaving a community
//
// @Summary         LeaveCommunity
// @Description     leave a community based on the provided community ID in path parameter
// @Tags            Community
// @Produce         json
// @Param           id      path        int     true   "Community ID"                   default(1)
// @Security        Bearer
// @Router          /communities/{id}/leave [delete]
// @Success         200     {object}    dto.Response
// @Failure         400     {object}    dto.ErrResponse
// @Failure         401     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (ch *CommunityHandler) LeaveCommunity(ctx *gin.Context) {
	payload, err := ch.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[CommunityHandler.LeaveCommunity] GetPayload: %s\n", err.Error())
		ctx.JSON(http.StatusUnauthorized, dto.ErrResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	param, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("[CommunityHandler.LeaveCommunity] Atoi: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrResponse{
			Success: false,
			Message: "Invalid community ID parameter",
		})
		return
	}

	if payload.Role != "attendee" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: fmt.Sprintf("%s cannot join communities", payload.Role),
		})
		return
	}

	joined, err := ch.cs.ToggleJoinCommunity(ctx.Request.Context(), int32(param), payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
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

// GetDetailCommunity handles fetching detail information of a specific community
//
// @Summary         Get Detail Community
// @Description     Retrieve detail information of a community by its ID
// @Tags            Community
// @Produce         json
// @Param           id      path        int     true   "Community ID"                   default(1)
// @Router          /communities/{id} [get]
// @Success         200     {object}    dto.Response
// @Failure         400     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (ch *CommunityHandler) GetDetailCommunity(ctx *gin.Context) {
	param, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("[CommunityHandler.GetDetailCommunity] Atoi: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrResponse{
			Success: false,
			Message: "Invalid community ID parameter",
		})
		return
	}

	communitiy, err := ch.cs.DetailCommunity(ctx.Request.Context(), int32(param))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
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

// GetMembersCommunity handles fetching the list of members in a specific community
//
// @Summary         Get Community Members
// @Description     Retrieve all members associated with a specific community ID
// @Tags            Community
// @Produce         json
// @Param           id      path        int     true   "Community ID"                   default(1)
// @Router          /communities/{id}/members [get]
// @Success         200     {object}    dto.Response
// @Failure         400     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (ch *CommunityHandler) GetMembersCommunity(ctx *gin.Context) {
	param, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("[CommunityHandler.GetMembersCommunity] Atoi: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrResponse{
			Success: false,
			Message: "Invalid community ID parameter",
		})
		return
	}

	membersCommunities, err := ch.cs.GetMembersCommunity(ctx.Request.Context(), int32(param))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
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
