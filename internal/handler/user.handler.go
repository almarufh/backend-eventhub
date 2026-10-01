package handler

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/middleware"
	"backend/EventHub/internal/service"
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
// @Summary			Get Profiles
// @Description		for get detail information profiles
// @Tags			User
// @Produce			json
// @Security		Bearer
// @Router			/user/profiles	[get]
// @Success			200		{object}	dto.Response
// @Failure			401		{object}	dto.ErrResponse
// @Failure			500		{object}	dto.ErrResponse
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

func (uh *UserHandler) SetProfile(ctx *gin.Context) {
	payload, err := uh.middleware.GetPayload(ctx)
	if err != nil || payload == nil {
		log.Printf("[UserHandler.Logout] Unauthorized: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var body dto.DtoSetProfile
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[UserHandler.SetProfile] Bind error: %v\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	profile, err := uh.us.SetProfileService(ctx.Request.Context(), (payload.ID), body)
	if err != nil {
		log.Printf("[UserHandler.SetProfile] Service error: %s\n", err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Gagal memperbarui profil",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Profil berhasil diperbarui",
		Data:    profile,
	})
}
