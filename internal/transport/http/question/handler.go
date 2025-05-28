package question

import (
	"diprec_api/internal/domain"
	"diprec_api/internal/pkg/utils"
	"diprec_api/internal/usecase/question"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type QuestionHandler struct {
	qu     question.IQuestionUsecase
	logger *zap.Logger
}

func NewQuestionHandler(qu question.IQuestionUsecase, logger *zap.Logger) *QuestionHandler {
	return &QuestionHandler{
		qu:     qu,
		logger: logger.Named("QuestionHandler"),
	}
}

// Create godoc
// @Summary Создать вопрос
// @Tags Question
// @Security BearerAuth
// @Produce json
// @Param input body CreateQuestionDTO true "ДТО создания вопроса"
// @Success 201 {object} domain.QuestionResponse
// @Error 400 {object} domain.Error
// @Error 400 {object} domain.Error
// @Error 401 {object} domain.Error
// @Error 500 {object} domain.Error
// @Router /question [post]
func (h *QuestionHandler) Create(c *gin.Context) {
	var req CreateQuestionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	question, err := h.qu.Create(c.Request.Context(), &domain.Question{
		Title:    req.Title,
		Type:     domain.Type(req.Type),
		Variants: utils.ParseMapToJSON(req.Variants),
		Answer:   utils.ParseToJSON(req.Answer),
	})
	if err != nil {
		h.logger.Warn("Create error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := question.ToQuestionResponse(true)
	c.JSON(http.StatusCreated, response)
}

// GetAll godoc
// @Summary Получить все вопросы
// @Tags Question
// @Security BearerAuth
// @Produce json
// @Success 200 {array} domain.QuestionResponse
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /question [get]
func (h *QuestionHandler) GetAll(c *gin.Context) {
	questions, err := h.qu.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Warn("GetAll error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := domain.ToQuestionsResponse(questions, true)
	c.JSON(http.StatusOK, response)
}

// GetByID Get godoc
// @Summary Получить вопрос по ID
// @Tags Question
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID вопроса"
// @Success 200 {object} domain.QuestionResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /question/{id} [get]
func (h *QuestionHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	question, err := h.qu.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Warn("GetByID error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := question.ToQuestionResponse(true)
	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary Обновить вопрос
// @Tags Question
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID вопроса"
// @Param input body UpdateQuestionDTO true "ДТО обновления вопроса"
// @Success 200 {object} domain.QuestionResponse
// @Success 400 {object} domain.Error
// @Success 400 {object} domain.Error
// @Success 500 {object} domain.Error
// @Router /question/{id} [put]
func (h *QuestionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	var req UpdateQuestionDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	question, err := h.qu.Update(c.Request.Context(), &domain.Question{
		ID:       uint(id),
		Title:    req.Title,
		Type:     domain.Type(req.Type),
		Variants: utils.ParseMapToJSON(req.Variants),
		Answer:   utils.ParseToJSON(req.Answer),
	})
	if err != nil {
		h.logger.Warn("Update error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := question.ToQuestionResponse(true)
	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Удалить вопрос
// @Tags Question
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID вопроса"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /question/{id} [delete]
func (h *QuestionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	err = h.qu.Delete(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Warn("Delete error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// Check godoc
// @Summary Проверить вопрос
// @Tags Question
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID вопроса"
// @Param input body CheckAnswerDTO true "ДТО ответа на вопрос"
// @Success 200 {object} domain.QuestionAnswer
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /question/{id}/check [post]
func (h *QuestionHandler) Check(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	var req CheckAnswerDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	userID := c.GetUint("userID")

	result, err := h.qu.Check(c.Request.Context(), uint(id), userID, req.Answer, req.TestId)
	if err != nil {
		h.logger.Warn("Check error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
