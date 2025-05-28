package test

import (
	"diprec_api/internal/domain"
	"diprec_api/internal/usecase/test"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TestHandler struct {
	tu     test.ITestUsecase
	logger *zap.Logger
}

func NewTestHandler(tu test.ITestUsecase, logger *zap.Logger) *TestHandler {
	return &TestHandler{
		tu:     tu,
		logger: logger.Named("TestHandler"),
	}
}

// Create godoc
// @Summary Создать тест
// @Tags Test
// @Security BearerAuth
// @Produce json
// @Param input body CreateTestDTO true "Название, описание и дедлайн теста"
// @Param id path int true "ID курса"
// @Success 201 {object} domain.TestResponse
// @Error 400 {object} domain.Error
// @Error 400 {object} domain.Error
// @Error 401 {object} domain.Error
// @Error 500 {object} domain.Error
// @Router /test/{id} [post]
func (h *TestHandler) Create(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	var req CreateTestDTO
	if err := c.ShouldBind(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	test, err := h.tu.Create(c.Request.Context(), &domain.Test{
		Name:        req.Name,
		Description: req.Description,
		Deadline:    req.Deadline,
	}, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := test.ToTestResponse()
	c.JSON(http.StatusCreated, response)
}

// GetByID Get godoc
// @Summary Получить тест по ID
// @Tags Test
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID теста"
// @Success 200 {object} domain.TestResponseWithQuestions
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /test/{id} [get]
func (h *TestHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	test, err := h.tu.GetByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		h.logger.Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	isTeacher := c.GetString("role")
	response := test.ToTestResponseWithQuestions(isTeacher == domain.RoleTeacher.String())
	c.JSON(http.StatusOK, response)
}

// Update godoc
// @Summary Обновить тест
// @Tags Test
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID теста"
// @Param input body UpdateTestDTO true "Название, описание и дедлайн теста"
// @Success 200 {object} domain.TestResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /test/{id} [put]
func (h *TestHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	var req UpdateTestDTO
	if err := c.ShouldBind(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	test, err := h.tu.Update(c.Request.Context(), &domain.Test{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		Deadline:    req.Deadline,
	})
	if err != nil {
		h.logger.Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := test.ToTestResponse()
	c.JSON(http.StatusOK, response)
}

// Delete godoc
// @Summary Удалить тест
// @Tags Test
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID теста"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /test/{id} [delete]
func (h *TestHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	if err := h.tu.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// AttachQuestion godoc
// @Summary Прикрепить вопрос к тесту
// @Tags Test
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID теста"
// @Param input body AttachQuestionDTO true "ID вопроса"
// @Success 200 "Вопрос прикреплен"
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /test/{id}/question [post]
func (h *TestHandler) AttachQuestion(c *gin.Context) {
	idStr := c.Param("id")
	testID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	var req AttachQuestionDTO
	if err := c.ShouldBind(&req); err != nil {
		h.logger.Warn("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	err = h.tu.AttachQuestion(c.Request.Context(), uint(testID), uint(req.QuestionID))
	if err != nil {
		h.logger.Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DetachQuestion godoc
// @Summary Открепить вопрос от теста
// @Tags Test
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param testID path int true "ID теста"
// @Param questionID path int true "ID вопроса"
// @Success 200 "Вопрос откреплён"
// @Failure 400 {object} domain.Error "Неверный запрос"
// @Failure 401 {object} domain.Error "Unauthorized"
// @Failure 500 {object} domain.Error "Internal Server Error"
// @Router /test/delete/{testId}/{questionId} [delete]
func (h *TestHandler) DetachQuestion(c *gin.Context) {
	// Парсим ID теста
	testID, err := strconv.Atoi(c.Param("testId"))
	if err != nil {
		h.logger.Warn("Invalid test ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: "Неверный ID теста"})
		return
	}

	questionID, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		h.logger.Warn("Invalid question ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: "Неверный ID вопроса"})
		return
	}

	// Выполняем детач
	if err := h.tu.DetachQuestion(c.Request.Context(), uint(testID), uint(questionID)); err != nil {
		h.logger.Error("DetachQuestion failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// StartTest godoc
// @Summary Запустить тест (учитель)
// @Tags Test
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID теста"
// @Success 200 {object} domain.TestResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /test/{id}/start [put]
func (h *TestHandler) StartTest(c *gin.Context) {
	idStr := c.Param("id")
	testID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error, invalid test ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	test, err := h.tu.Update(c.Request.Context(), &domain.Test{
		ID:     uint(testID),
		Status: domain.Progress,
	})
	if err != nil {
		h.logger.Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := test.ToTestResponse()

	c.JSON(http.StatusOK, response)
}

// StopTest godoc
// @Summary Остановить тест (учитель)
// @Tags Test
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID теста"
// @Success 200 {object} domain.TestResponse
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /test/{id}/stop [put]
func (h *TestHandler) StopTest(c *gin.Context) {
	idStr := c.Param("id")
	testID, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Validation error, invalid test ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	test, err := h.tu.Update(c.Request.Context(), &domain.Test{
		ID:     uint(testID),
		Status: domain.Ended,
	})
	if err != nil {
		h.logger.Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	response := test.ToTestResponse()
	c.JSON(http.StatusOK, response)
}

// BeginTest godoc
// @Summary Приступить к тесту (студент)
// @Tags Test
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID теста"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /test/{id}/begin [post]
func (h *TestHandler) BeginTest(c *gin.Context) {
	userID := c.GetUint("userID")

	testID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Validation error, invalid test ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	err = h.tu.StartTest(c.Request.Context(), &domain.UserTests{
		UserID: userID,
		TestID: uint(testID),
	})
	if err != nil {
		h.logger.Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

// POST /test/{course_id}/recommend   (внутренний вызов от Python)
func (h *TestHandler) CreateRecommend(c *gin.Context) {
	courseID, err := strconv.Atoi(c.Param("course_id"))
	if err != nil {
		h.logger.Warn("invalid course_id", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: "invalid id"})
		return
	}

	var req RecommendTestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("bind error", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	newTest, err := h.tu.CreateRecommendTest(
		c.Request.Context(),
		&domain.Test{
			Name:        req.Name,
			Description: req.Description,
			Deadline:    req.Deadline,
		},
		uint(courseID),
		req.UserID,
		req.QuestionIDs,
	)
	if err != nil {
		h.logger.Error("create recommend test", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTest.ToTestResponse())
}

// FinishTest godoc
// @Summary Завершить тест (студент)
// @Tags Test
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID теста"
// @Param input body FinishTestDTO true "Результат по тесту в процентах"
// @Success 200
// @Failure 400 {object} domain.Error
// @Failure 401 {object} domain.Error
// @Failure 500 {object} domain.Error
// @Router /test/{id}/finish [put]
func (h *TestHandler) FinishTest(c *gin.Context) {
	userID := c.GetUint("userID")
	testID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Warn("Validation error, invalid test ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	var req FinishTestDTO
	if err := c.ShouldBind(&req); err != nil {
		h.logger.Warn("Validation error, invalid body", zap.Error(err))
		c.JSON(http.StatusBadRequest, domain.Error{Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	err = h.tu.EndTest(c.Request.Context(), &domain.UserTests{
		UserID:   userID,
		TestID:   uint(testID),
		Progress: req.Progress,
		Status:   domain.Ended,
	})
	if err != nil {
		h.logger.Warn("Internal error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, domain.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
