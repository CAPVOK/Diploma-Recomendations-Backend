package test

import "time"

type CreateTestDTO struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}

type UpdateTestDTO struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
}

type AttachQuestionDTO struct {
	QuestionID int `json:"questionId"`
}

type FinishTestDTO struct {
	Progress uint `json:"progress"`
}

type RecommendTestDTO struct {
	Name        string    `json:"name"         binding:"required"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"     binding:"required"`
	UserID      uint      `json:"user_id"      binding:"required"`
	QuestionIDs []uint    `json:"question_ids" binding:"required"`
}
