package course

import (
	"context"
	"diprec_api/internal/domain"
	"diprec_api/internal/repository/course"

	"go.uber.org/zap"
)

type courseUsecase struct {
	repo   course.ICourseRepository
	logger *zap.Logger
}

type ICourseUsecase interface {
	Create(ctx context.Context, course *domain.Course) (*domain.Course, error)
	Update(ctx context.Context, course *domain.Course) (*domain.Course, error)
	Delete(ctx context.Context, id uint) error
	GetById(ctx context.Context, id, userID uint) (*domain.Course, error)
	Get(ctx context.Context) ([]*domain.Course, error)
	Enroll(ctx context.Context, courseID uint, userID uint) error
}

func NewCourseUseCase(repo course.ICourseRepository, logger *zap.Logger) ICourseUsecase {
	return &courseUsecase{
		repo:   repo,
		logger: logger.Named("CourseUsecase"),
	}
}

func (u *courseUsecase) Create(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	if err := u.repo.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

func (u *courseUsecase) Update(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	if err := u.repo.Update(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

func (u *courseUsecase) Delete(ctx context.Context, id uint) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (u *courseUsecase) GetById(ctx context.Context, id, userID uint) (*domain.Course, error) {
	course, err := u.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (u *courseUsecase) Get(ctx context.Context) ([]*domain.Course, error) {
	courses, err := u.repo.Get(ctx)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (u *courseUsecase) Enroll(ctx context.Context, courseID uint, userID uint) error {
	err := u.repo.EnrollUser(ctx, courseID, userID)
	if err != nil {
		return err
	}

	return nil
}
