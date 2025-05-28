package question

import (
	"context"
	"diprec_api/internal/domain"
	"diprec_api/internal/pkg/validator"
	"errors"

	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

type IQuestionRepository interface {
	Create(ctx context.Context, question *domain.Question) error
	GetByID(ctx context.Context, id uint) (*domain.Question, error)
	GetAll(ctx context.Context) ([]*domain.Question, error)
	Update(ctx context.Context, question *domain.Question) error
	Delete(ctx context.Context, id uint) error
}

func NewQuestionRepository(db *gorm.DB) IQuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, question *domain.Question) error {
	if err := r.db.Create(question).Error; err != nil {
		return err
	}

	return nil
}

func (r *questionRepository) GetAll(ctx context.Context) ([]*domain.Question, error) {
	var questions []*domain.Question

	err := r.db.Find(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (r *questionRepository) GetByID(ctx context.Context, id uint) (*domain.Question, error) {
	var question *domain.Question

	err := r.db.First(&question, id).Error
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (r *questionRepository) Update(ctx context.Context, question *domain.Question) error {
	updates := validator.BuildUpdates(question)

	result := r.db.Model(&domain.Question{}).Where("id = ?", question.ID).Updates(updates)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrQuestionNotFound
		}

		return result.Error
	}

	return nil
}

func (r *questionRepository) Delete(ctx context.Context, id uint) error {
	var question domain.Question

	if err := r.db.First(&question, id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&question).Association("Tests").Clear(); err != nil {
		return err
	}

	return r.db.Delete(&question).Error
}
