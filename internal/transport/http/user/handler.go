package user

import (
	"diprec_api/internal/domain"
	"diprec_api/internal/usecase/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	uc     user.IUserUseCase
	logger *zap.Logger
}

func NewUserHandler(uc user.IUserUseCase, logger *zap.Logger) *UserHandler {

	return &UserHandler{
		uc:     uc,
		logger: logger.Named("UserHandler"),
	}
}

// Register godoc
// @Summary Зарегистрировать нового пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body CreateUserDTO true "Данные пользователя"
// @Success 201 {object} domain.AuthResponse "Пользователь успешно зарегистрирован"
// @Failure 400 {object} domain.Error "Неверный формат запроса / тело запроса"
// @Failure 401 {object} domain.Error "Ошибка авторизации"
// @Failure 409 {object} domain.Error "Пользователь с таким именем уже существует"
// @Failure 500 {object} domain.Error "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req CreateUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	user, err := h.uc.Register(c.Request.Context(), &domain.User{
		Username:   req.Username,
		Password:   req.Password,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Patronymic: req.Patronymic,
	})
	if err != nil {
		h.logger.Error("Register error", zap.Error(err))
		c.JSON(errorStatusCode(err), domain.Error{Message: err.Error()})
		return
	}

	response := user.ToUserResponse()

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Аутентификация пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body LoginUserDTO true "Данные для авторизации"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginUserDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	user, err := h.uc.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.Warn("Authenticate error", zap.Error(err))
		c.JSON(errorStatusCode(err), domain.Error{Message: err.Error()})
		return
	}

	tokens, err := h.uc.GenerateTokens(user)
	if err != nil {
		h.logger.Warn("GenerateTokens error", zap.Error(err))
		c.JSON(errorStatusCode(err), domain.Error{Message: err.Error()})
		return
	}

	response := user.ToAuthResponse(tokens)
	c.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Summary Обновить токены авторизации
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body RefreshUserDTO true "Refresh Token"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Router /auth/refresh [post]
func (h *UserHandler) Refresh(c *gin.Context) {
	var req RefreshUserDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	pair, err := h.uc.RefreshTokens(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Warn("Refresh error", zap.Error(err))
		c.JSON(errorStatusCode(err), domain.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pair)
}

// Me godoc
// @Summary Получение информации о текущем пользователе
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.UserResponseWithCourses
// @Failure 401 {object} domain.Error
// @Router /user/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.uc.GetMe(c.Request.Context(), userID)
	if err != nil {
		h.logger.Warn("GetMe", zap.Error(err))
		c.JSON(errorStatusCode(err), domain.Error{Message: err.Error()})
		return
	}

	response := user.ToUserResponseWithCourses()
	c.JSON(http.StatusOK, response)
}

func errorStatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrInvalidRefreshToken):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrUserExists):
		return http.StatusConflict
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
