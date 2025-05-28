package domain

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null;unique"`
	Description string
	Users       []*User `gorm:"many2many:user_courses;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tests       []*Test `gorm:"many2many:course_tests;constraint:OnUpdate:CASCADE;OnDelete:CASCADE;"`
}

type CourseResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CourseResponseWithTests struct {
	CourseResponse
	Tests []TestResponse `json:"tests"`
}

func (c *Course) ToCourseResponse() CourseResponse {
	return CourseResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func (c *Course) ToCourseResponseWithTests() CourseResponseWithTests {
	return CourseResponseWithTests{
		CourseResponse: c.ToCourseResponse(),
		Tests:          ToTestsResponse(c.Tests),
	}
}

func ToCoursesResponse(courses []*Course) []CourseResponse {
	responses := make([]CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = course.ToCourseResponse()
	}
	return responses
}
