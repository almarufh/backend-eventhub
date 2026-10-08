package handler

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/internal/middleware"
	"backend/EventHub/internal/service"
	"backend/EventHub/pkg"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	as         *service.AuthService
	middleware *middleware.Middleware
}

func NewAuthHandler(as *service.AuthService, m *middleware.Middleware) *AuthHandler {
	return &AuthHandler{
		as:         as,
		middleware: m,
	}
}

// Register
//
// @Summary			Create New Account
// @Description		Create new account for full access
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Router			/auth/register	[post]
// @Param			data	body dto.ReqRegister true "Body data new user"
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrResponse
// @Failure			409		{object}	dto.ErrResponse
// @Failure			500		{object}	dto.ErrResponse
func (ah *AuthHandler) Register(ctx *gin.Context) {
	var body dto.ReqRegister
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[AuthHandler.Register] Bind error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	payload, err := ah.middleware.GetPayload(ctx)

	if err != nil || payload == nil || payload.Role != "admin" {
		body.Role = "attendee"
	}

	_, err = ah.as.RegisterAuthService(ctx.Request.Context(), body)
	if err != nil {
		log.Printf("[AuthHandler.Register] Service error: %s\n", err.Error())

		switch {
		case errors.Is(err, service.ErrFieldEmpty):
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: "Email, name, and password are required",
			})

		case errors.Is(err, pkg.ErrEmailExists):
			ctx.JSON(http.StatusConflict, dto.ErrResponse{
				Success: false,
				Message: "Email is already registered",
			})

		case errors.Is(err, pkg.ErrConflict):
			ctx.JSON(http.StatusConflict, dto.Response{
				Success: false,
				Message: "Resource already exists",
			})

		case errors.Is(err, pkg.ErrCheckViolation), errors.Is(err, pkg.ErrNotNullFail):
			ctx.JSON(http.StatusBadRequest, dto.ErrResponse{
				Success: false,
				Message: "Invalid input data",
			})

		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
				Success: false,
				Message: "Internal server error",
			})
		}
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "User registered successfully",
	})
}

// Login
//
// @Summary			Log In
// @Description		Log in for get full access
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Router			/auth/login	[post]
// @Param			data	body dto.ReqLogin true "Body data to access EventHub"
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrResponse
func (ah *AuthHandler) Login(ctx *gin.Context) {
	origin := ctx.GetHeader("Origin")
	log.Printf("Origin Auth : %s\n\n", origin)
	var body dto.ReqLogin
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[AuthHandler.Login] Bind error: %v\n", err.Error())
		parsedErrors := ah.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validasi gagal",
			Data:    parsedErrors,
		})
		return
	}

	client, _ := ah.middleware.GetClientInfo(ctx)
	body.Device = client.UserAgent

	auth, err := ah.as.LoginAuthService(ctx.Request.Context(), body)
	if err != nil {
		log.Printf("[AuthHandler.Login] Service error: %s\n", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.ErrResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "User login successfully",
		Data:    auth,
	})
}

func (ah *AuthHandler) FormatValidationError(err error) map[string]string {
	var ve validator.ValidationErrors
	errsMap := make(map[string]string)

	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				errsMap[field] = fmt.Sprintf("Field %s wajib diisi", field)
			case "email":
				errsMap[field] = fmt.Sprintf("Format %s tidak valid", field)
			case "min":
				errsMap[field] = fmt.Sprintf("Field %s minimal harus %s karakter", field, fe.Param())
			case "max":
				errsMap[field] = fmt.Sprintf("Field %s maksimal %s karakter", field, fe.Param())
			default:
				errsMap[field] = fmt.Sprintf("Field %s tidak valid (%s)", field, fe.Tag())
			}
		}
		return errsMap
	}

	errsMap["error"] = err.Error()
	return errsMap
}

// Logout
//
// @Summary			Log Out
// @Description		Log Out for delete acces
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Router			/auth/logout	[delete]
// @Security		Bearer
// @Success			200		{object}	dto.Response
// @Failure			400		{object}	dto.ErrResponse
// @Failure			401		{object}	dto.ErrResponse
// @Failure			500		{object}	dto.ErrResponse
func (ah *AuthHandler) Logout(ctx *gin.Context) {
	rawPayload, ok := ctx.Get("payload")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	payload, ok := rawPayload.(*middleware.Payload)
	if !ok {
		log.Println("[AuthHandler.Logout] Unauthorized: invalid payload type")
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	err := ah.as.LogoutAuthService(ctx.Request.Context(), payload)
	if err != nil {
		log.Printf("[AuthHandler.Logout] Service error: %s\n", err.Error())

		switch {
		case errors.Is(err, pkg.ErrNotFound):
			ctx.JSON(http.StatusNotFound, dto.ErrResponse{
				Success: false,
				Message: "Session not found or already logged out",
			})

		default:
			ctx.JSON(http.StatusInternalServerError, dto.ErrResponse{
				Success: false,
				Message: "Failed to logout",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Logged out successfully",
	})
}

// Create Password
//
// @Summary      Forgot Password
// @Description  Create new password when forgot old password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        data body dto.ReqChangePassword true "Body to Create New Password"
// @Success      200 {object} dto.Response
// @Failure      400 {object} dto.ErrResponse
// @Router       /auth/forgot/password [patch]
func (ah *AuthHandler) ChangePassword(ctx *gin.Context) {
	var body dto.ReqChangePassword
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Printf("[AuthHandler.ChangePassword] Bind error: %v\n", err.Error())
		parsedErrors := ah.FormatValidationError(err)
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validasi gagal",
			Data:    parsedErrors,
		})
		return
	}

	err := ah.as.ChangePassword(ctx.Request.Context(), body)
	if err != nil {
		log.Printf("[AuthHandler.Logout] Service error: %s\n", err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Reset new password successfully",
	})
}
