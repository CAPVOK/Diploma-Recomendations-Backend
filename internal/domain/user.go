package domain

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Username   string    `json:"username" gorm:"not null;unique"`
	Password   string    `json:"password" gorm:"not null"`
	FirstName  string    `json:"firstName" gorm:"not null"`
	LastName   string    `json:"lastName" gorm:"not null"`
	Patronymic string    `json:"patronymic,omitempty"`
	Role       Role      `json:"role" gorm:"type:varchar(20);not null;role IN ('STUDENT', 'TEACHER');default:'STUDENT'"`
	Courses    []*Course `gorm:"many2many:user_courses;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	Tests      []*Test   `gorm:"many2many:user_test;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}

type Role string

const (
	RoleTeacher Role = "TEACHER"
	RoleStudent Role = "STUDENT"
)

type UserResponse struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Patronymic string    `json:"patronymic"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type UserResponseWithCourses struct {
	UserResponse
	Courses []CourseResponse `json:"courses"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresAt    time.Time    `json:"expiresAt"`
}

func (r Role) String() string {
	return string(r)
}

type TokenPair struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) ToUserResponseWithCourses() UserResponseWithCourses {
	return UserResponseWithCourses{
		UserResponse: u.ToUserResponse(),
		Courses:      ToCoursesResponse(u.Courses),
	}
}

func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Patronymic: u.Patronymic,
		Role:       u.Role.String(),
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

func (u *User) ToAuthResponse(tokens *TokenPair) AuthResponse {
	return AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    tokens.ExpiresAt,
		User:         u.ToUserResponse(),
	}
}
