package domain

import (
	"time"

	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null;unique"`
	Description string
	Status      TestStatus  `json:"status" gorm:"type:varchar(20);not null;test_status IN ('DRAFT', 'PROGRESS', 'ENDED');default:'DRAFT'"`
	Assignee    Assignee    `json:"assignee" gorm:"type:varchar(20);not null;assignee IN ('TEACHER', 'RECOMMENDATION');default:'TEACHER'"`
	Deadline    time.Time   `gorm:"not null"`
	Courses     []*Course   `gorm:"many2many:course_tests;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Questions   []*Question `gorm:"many2many:test_questions;constraint:OnUpdate:CASCADE;OnDelete:CASCADE;"`
	UserTests   UserTests   `gorm:"foreignKey:test_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TestStatus string

const (
	Draft    TestStatus = "DRAFT"
	Progress TestStatus = "PROGRESS"
	Ended               = "ENDED"
)

func (t TestStatus) String() string {
	return string(t)
}

type TestResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Status           string    `json:"status"`
	Assignee         string    `json:"assignee"`
	Deadline         time.Time `json:"deadline"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	UserTestResponse `json:"result,omitempty"`
}

type TestResponseWithQuestions struct {
	TestResponse
	Questions []QuestionResponse `json:"questions"`
}

type Assignee string

const (
	Teacher        Assignee = "TEACHER"
	Recommendation Assignee = "RECOMMENDATION"
)

func (a Assignee) String() string {
	return string(a)
}

func (c *Test) ToTestResponse() TestResponse {
	return TestResponse{
		ID:               c.ID,
		Name:             c.Name,
		Description:      c.Description,
		Status:           c.Status.String(),
		Assignee:         c.Assignee.String(),
		Deadline:         c.Deadline,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
		UserTestResponse: c.UserTests.ToUserTestResponse(),
	}
}

func (c *Test) ToTestResponseWithQuestions(isTeacher bool) TestResponseWithQuestions {
	return TestResponseWithQuestions{
		TestResponse: c.ToTestResponse(),
		Questions:    ToQuestionsResponse(c.Questions, isTeacher),
	}
}

func ToTestsResponse(test []*Test) []TestResponse {
	response := make([]TestResponse, len(test))
	for i, test := range test {
		response[i] = test.ToTestResponse()
	}

	return response
}
