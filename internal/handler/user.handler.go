package handler

import (
	"backend/EventHub/internal/dto"
	msgerr "backend/EventHub/internal/message"
	"backend/EventHub/internal/middleware"
	"backend/EventHub/internal/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserHandler struct {
	us         *service.UserService
	middleware *middleware.Middleware
}

func NewUserHandler(m *middleware.Middleware, us *service.UserService) *UserHandler {
	return &UserHandler{
		us:         us,
		middleware: m,
	}
}

// Change Password
//
// @Summary			Change Password
// @Description		Change password
// @Tags			User
// @Accept			json
// @Produce			json
// @Security		Bearer
// @Router			/user/set/password	[patch]
// @Param			data	body 	dto.ReqNewPassword true "Body to Change Password"
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrResponse
// @Failure			401		{object}	dto.ErrResponse
// @Failure			500		{object}	dto.ErrResponse
func (uh *UserHandler) NewPassword(ctx *gin.Context) {
	payload, err := uh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[UserHandler.Logout] Unauthorized: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, dto.ErrResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var body dto.ReqNewPassword
	body.ID = payload.ID

	err = ctx.ShouldBindWith(&body, binding.JSON)
	if err != nil {
		log.Printf("[UserHandler.ChangePassword] Bind error: %v\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrResponse{
			Success: false,
			Message: "Validasi gagal",
			Data:    err.Error(),
		})
		return
	}

	err = uh.us.NewPassword(ctx.Request.Context(), body)
	if err != nil {
		log.Printf("[UserHandler.Logout] Service error: %s\n", err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Changed password successfully",
	})
}

// Get Profiles
//
// @Summary         Get Profiles
// @Description     for get detail information profiles
// @Tags            User
// @Produce         json
// @Security        Bearer
// @Router          /user/profiles  [get]
// @Success         200     {object}    dto.Response
// @Failure         401     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (uh *UserHandler) MyProfile(ctx *gin.Context) {
	payload, err := uh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[UserHandler.Logout] Unauthorized: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, dto.ErrResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	user, err := uh.us.MyProfile(ctx.Request.Context(), payload.ID)
	if err != nil {
		log.Printf("[UserHandler.Logout] Service error: %s\n", err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "get profile successfully",
		Data:    user,
	})
}

// / SetProfile
//
// @Summary         Update Profile
// @Description     Update user profile information and upload a profile picture (max 2MB, formats: jpg/jpeg/png)
// @Tags            User
// @Accept          multipart/form-data
// @Produce         json
// @Param name        formData string false "User Full Name" default(Alma'ruf Hidayat)
// @Param address     formData string false "User Address"   default(Masamba, Luwu Utara, Sulawesi Selatan)
// @Param job         formData string false "User Job Title" default(Fullstack Developer)
// @Param office      formData string false "User Office"    default(PT. Best Life Ummah)
// @Param description formData string false "Description"    default(Senior Fullstack Developer)
// @Param image       formData file   false "Profile Picture File"
// @Security        Bearer
// @Router          /user/set/profiles [patch]
// @Success         200     {object}    dto.Response{data=dto.ResSetProfile}
// @Failure         400     {object}    dto.ErrResponse
// @Failure         401     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (uh *UserHandler) SetProfile(ctx *gin.Context) {
	payload, err := uh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[UserHandler.SetProfile] Unauthorized: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, dto.ErrResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var body dto.ReqSetProfile

	if err := ctx.ShouldBindWith(&body, binding.FormMultipart); err != nil {
		log.Printf("[UserHandler.SetProfile] Bind error: %v\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	profile, err := uh.us.SetProfileService(ctx.Request.Context(), payload.ID, body)
	if err != nil {
		log.Printf("[UserHandler.SetProfile] Service error: %s\n", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Profil berhasil diperbarui",
		Data:    profile,
	})
}

// JoinedEvents
//
// @Summary         Get Joined Events
// @Description     for get list of events joined by the authenticated user
// @Tags            User
// @Produce         json
// @Security        Bearer
// @Router          /user/events/joined  [get]
// @Success         200     {object}    dto.Response
// @Failure         401     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (uh *UserHandler) JoinedEvents(ctx *gin.Context) {
	payload, err := uh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[UserHandler.JoinedEvents] Unauthorized: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, dto.ErrResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	events, err := uh.us.JoinedEvents(ctx.Request.Context(), payload.ID)
	if err != nil {
		log.Printf("[UserHandler.JoinedEvents] Error: %v\n", err)
		if errors.Is(err, msgerr.SavedEventsUserNotFound) {
			ctx.JSON(http.StatusOK, dto.ErrResponse{
				Success: false,
				Message: "tidak memiliki daftar event yang diikuti",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
			Success: false,
			Message: "Gagal mengambil daftar event yang diikuti",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil mengambil daftar event yang diikuti",
		Data:    events,
	})
}

// SavedEvents
//
// @Summary         Get Saved Events
// @Description     for get list of events joined by the authenticated user
// @Tags            User
// @Produce         json
// @Security        Bearer
// @Router          /user/events/saved  [get]
// @Success         200     {object}    dto.Response
// @Failure         401     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (uh *UserHandler) SavedEvents(ctx *gin.Context) {
	payload, err := uh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[UserHandler.JoinedEvents] Unauthorized: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, dto.ErrResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	events, err := uh.us.SavedEvents(ctx.Request.Context(), payload.ID)
	if err != nil {
		log.Printf("[UserHandler.JoinedEvents] Error: %v\n", err)
		if errors.Is(err, msgerr.SavedEventsUserNotFound) {
			ctx.JSON(http.StatusOK, dto.ErrResponse{
				Success: false,
				Message: "tidak memiliki daftar event yang disimpan",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
			Success: false,
			Message: "Gagal mengambil daftar event yang disimpan",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil mengambil daftar event yang disimpan",
		Data:    events,
	})
}

// JoinedCommunities
//
// @Summary         Get Joined Communities
// @Description     for get list of events joined by the authenticated user
// @Tags            User
// @Produce         json
// @Security        Bearer
// @Router          /user/community/joined  [get]
// @Success         200     {object}    dto.Response
// @Failure         401     {object}    dto.ErrResponse
// @Failure         500     {object}    dto.ErrResponse
func (uh *UserHandler) JoinedCommunities(ctx *gin.Context) {
	payload, err := uh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[UserHandler.JoinedEvents] Unauthorized: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, dto.ErrResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	events, err := uh.us.JoinedCommunities(ctx.Request.Context(), payload.ID)
	if err != nil {
		log.Printf("[UserHandler.JoinedCommunities] Error: %v\n", err)
		if errors.Is(err, msgerr.SavedEventsUserNotFound) {
			ctx.JSON(http.StatusOK, dto.ErrResponse{
				Success: false,
				Message: "tidak memiliki daftar communities yang diikuti",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
			Success: false,
			Message: "Gagal mengambil daftar communities yang diikuti",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Berhasil mengambil daftar communities yang diikuti",
		Data:    events,
	})
}
