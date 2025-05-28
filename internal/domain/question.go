package domain

import (
	"diprec_api/internal/pkg/utils"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	ID       uint           `gorm:"primaryKey;autoIncrement"`
	Title    string         `gorm:"not null,unique"`
	Type     Type           `gorm:"type:varchar(20);not null;type IN ('SINGLE', 'MULTIPLE', 'TEXT', 'NUMBER');default:'SINGLE'"`
	Variants datatypes.JSON `gorm:"type:jsonb"`
	Answer   datatypes.JSON `gorm:"type:jsonb"`
	Tests    []Test         `gorm:"many2many:test_questions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Type string

const (
	Single   Type = "SINGLE"
	Multiple Type = "MULTIPLE"
	Text     Type = "TEXT"
	Number   Type = "NUMBER"
)

func (t Type) String() string {
	return string(t)
}

type QuestionResponse struct {
	ID       uint                   `json:"id"`
	Title    string                 `json:"title"`
	Type     string                 `json:"type" enums:"SINGLE,MULTIPLE,TEXT,NUMBER" example:"SINGLE"`
	Variants map[string]interface{} `json:"variants"`
	Answer   interface{}            `json:"answer"`
}

type QuestionAnswer struct {
	IsCorrect bool        `json:"isCorrect"`
	Message   string      `json:"message"`
	Answer    interface{} `json:"answer"`
}

type UserAnswer struct {
	QuestionID uint        `json:"question_id"`
	Title      string      `json:"question_title"`
	Type       string      `json:"question_type"`
	Variants   interface{} `json:"question_variants"`
	Answer     interface{} `json:"question_answer"`
	UserID     uint        `json:"user_id"`
	IsCorrect  bool        `json:"is_correct"`
	Timestamp  time.Time   `json:"timestamp"`
}

type UserAnswerCheck struct {
	QuestionID uint        `json:"question_id"`
	CourseID   uint        `json:"course_id"`
	Title      string      `json:"question_title"`
	Type       string      `json:"question_type"`
	Variants   interface{} `json:"question_variants"`
	Answer     interface{} `json:"question_answer"`
	UserID     uint        `json:"user_id"`
	IsCorrect  bool        `json:"is_correct"`
	Timestamp  time.Time   `json:"timestamp"`
}

func (q *Question) CheckAnswer(userAnswer interface{}) bool {
	var correctAnswer interface{}
	if err := json.Unmarshal(q.Answer, &correctAnswer); err != nil {
		return false
	}

	switch q.Type {
	case Single:
		correctStr, ok1 := correctAnswer.(string)
		userStr, ok2 := userAnswer.(string)
		return ok1 && ok2 && userStr == correctStr

	case Multiple:
		correctSlice, ok1 := utils.ToStringSlice(correctAnswer)
		userSlice, ok2 := utils.ToStringSlice(userAnswer)
		if !ok1 || !ok2 {
			return false
		}
		return utils.EqualStringSlices(correctSlice, userSlice)

	case Text:
		correctStr, ok1 := correctAnswer.(string)
		userStr, ok2 := userAnswer.(string)
		return ok1 && ok2 && strings.TrimSpace(strings.ToLower(userStr)) == strings.TrimSpace(strings.ToLower(correctStr))

	case Number:
		correctNum, ok1 := correctAnswer.(float64)
		userNum, ok2 := userAnswer.(float64)
		return ok1 && ok2 && userNum == correctNum
	}

	return false
}

func (c *Question) ToQuestionResponse(isTeacher bool) QuestionResponse {
	if isTeacher {
		return QuestionResponse{
			ID:       c.ID,
			Title:    c.Title,
			Type:     c.Type.String(),
			Variants: utils.ParseJSONToMap(c.Variants),
			Answer:   utils.ParseJSONInterface(c.Answer),
		}
	}

	return QuestionResponse{
		ID:       c.ID,
		Title:    c.Title,
		Type:     c.Type.String(),
		Variants: utils.ParseJSONToMap(c.Variants),
	}
}

func ToQuestionsResponse(questions []*Question, isTeacher bool) []QuestionResponse {
	responses := make([]QuestionResponse, len(questions))
	for i, question := range questions {
		responses[i] = question.ToQuestionResponse(isTeacher)
	}

	return responses
}
