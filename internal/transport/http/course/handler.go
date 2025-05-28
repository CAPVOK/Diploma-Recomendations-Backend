package course

import (
	"diprec_api/internal/domain"
	"diprec_api/internal/usecase/course"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CourseHandler struct {
	cu     course.ICourseUsecase
	logger *zap.Logger
}

func NewCourseHandler(cu course.ICourseUsecase, logger *zap.Logger) *CourseHandler {
	return &CourseHandler{
		cu:     cu,
		logger: logger.Named("CourseHandler"),
	}
}

// Create godoc
// @Summary Создать курс
// @Tags Course
// @Security BearerAuth
// @Produce json
// @Param input body CreateCourseDTO true "Название и описание курса"
// @Success 201 {object} domain.CourseResponse
// @Error 400 {object} domain.Error
// @Error 401 {object} domain.Error
// @Error 500 {object} domain.Error
// @Router /course [post]
func (h *CourseHandler) Create(c *gin.Context) {
	var req CreateCourseDTO
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	course, err := h.cu.Create(c.Request.Context(), &domain.Course{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := course.ToCourseResponse()
	c.JSON(http.StatusCreated, response)
}

// Get godoc
// @Summary Получить курсы
// @Tags Course
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.CourseResponse[]
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /course [get]
func (h *CourseHandler) Get(c *gin.Context) {
	courses, err := h.cu.Get(c.Request.Context())
	if err != nil {
		h.logger.Error("Get courses failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := domain.ToCoursesResponse(courses)
	c.JSON(http.StatusOK, response)
}

// GetByID Get godoc
// @Summary Получить курс по ID
// @Tags Course
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID курса"
// @Success 200 {object} domain.CourseResponseWithTests
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /course/{id} [get]
func (h *CourseHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	course, err := h.cu.GetById(c.Request.Context(), uint(id), userID)
	if err != nil {
		if errors.Is(err, domain.ErrCourseNotFound) {
			h.logger.Error("Get course failed", zap.Error(err))
			c.JSON(http.StatusNotFound, domain.Error{Message: err.Error()})
			return
		}
		h.logger.Error("Get course failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := course.ToCourseResponseWithTests()
	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary Обновить курс
// @Tags Course
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID курса"
// @Param input body UpdateCourseDTO true "Название и описание курса"
// @Success 200 {object} domain.CourseResponseWithTests
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /course/{id} [put]
func (h *CourseHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	var req UpdateCourseDTO
	if err := c.ShouldBind(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	course, err := h.cu.Update(c.Request.Context(), &domain.Course{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, domain.ErrCourseNotFound) {
			h.logger.Error("Get course failed", zap.Error(err))
			c.JSON(http.StatusNotFound, domain.Error{Message: err.Error()})
			return
		}
		h.logger.Error("Update course failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := course.ToCourseResponseWithTests()
	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Удалить курс
// @Tags Course
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID курса"
// @Success 200 "Курс успешно удалён"
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 404 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /course/{id} [delete]
func (h *CourseHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid course ID", zap.String("id", idStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	err = h.cu.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrCourseNotFound) {
			h.logger.Error("Get course failed", zap.Error(err))
			c.JSON(http.StatusNotFound, domain.Error{Message: err.Error()})
			return
		}

		h.logger.Error("Delete course failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// Enroll godoc
// @Summary Записаться на курс
// @Tags Course
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID курса"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /course/{id}/enroll [post]
func (h *CourseHandler) Enroll(c *gin.Context) {
	idStr := c.Param("id")
	courseID, err := strconv.Atoi(idStr)

	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	userID := c.GetUint("userID")

	err = h.cu.Enroll(c.Request.Context(), uint(courseID), userID)
	if err != nil {
		h.logger.Error("Enroll course failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
