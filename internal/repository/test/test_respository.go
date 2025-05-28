package test

import (
	"context"
	"diprec_api/internal/domain"
	"diprec_api/internal/pkg/validator"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type testRepository struct {
	db *gorm.DB
}

type ITestRepository interface {
	Create(ctx context.Context, test *domain.Test, courseID uint) error
	Get(ctx context.Context, courseID, userID uint) ([]*domain.Test, error)
	GetByID(ctx context.Context, id, userID uint) (*domain.Test, error)
	Update(ctx context.Context, test *domain.Test) error
	Delete(ctx context.Context, id uint) error
	AttachQuestion(ctx context.Context, testID uint, questionID uint) error
	DetachQuestion(ctx context.Context, testID uint, questionID uint) error
	UpdateUserTest(ctx context.Context, userTest *domain.UserTests) error
	CreateUserTest(ctx context.Context, userTest *domain.UserTests) error
	GetCourseIDByTestID(ctx context.Context, testID uint) (uint, error)
	CreateWithUser(ctx context.Context, test *domain.Test, courseID uint, userID uint) error
}

func NewTestRepository(db *gorm.DB) ITestRepository { return &testRepository{db: db} }

func (r *testRepository) Create(ctx context.Context, test *domain.Test, courseID uint) error {
	if err := r.db.Create(test).Error; err != nil {
		return err
	}

	return r.db.Model(test).Association("Courses").Append(&domain.Course{ID: courseID})
}

func (r *testRepository) Get(ctx context.Context, courseID, userID uint) ([]*domain.Test, error) {
	var course domain.Course

	err := r.db.Preload("Tests").First(&course, courseID, userID).Error
	if err != nil {
		return nil, err
	}

	return course.Tests, nil
}

func (r *testRepository) GetByID(ctx context.Context, id, userID uint) (*domain.Test, error) {
	var test domain.Test

	err := r.db.
		Preload("Questions", "deleted_at IS NULL").
		Preload("UserTests", "user_id = ?", userID). // ← фильтрация по userID
		Where("id = ?", id).
		First(&test).Error

	if err != nil {
		return nil, err
	}

	return &test, nil
}

func (r *testRepository) Update(ctx context.Context, test *domain.Test) error {
	updates := validator.BuildUpdates(test)

	result := r.db.Model(&domain.Test{}).Where("id = ?", test.ID).Updates(updates).First(&test)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrTestNotFound
		}

		return result.Error
	}

	return nil
}

func (r *testRepository) Delete(ctx context.Context, id uint) error {
	var test domain.Test

	if err := r.db.First(&test, id).Error; err != nil {
		return err
	}

	if err := r.db.Model(&test).Association("Courses").Clear(); err != nil {
		return err
	}

	return r.db.Delete(&test).Error
}

func (r *testRepository) AttachQuestion(ctx context.Context, testID uint, questionID uint) error {
	err := r.db.Model(&domain.Test{ID: testID}).Association("Questions").Append(&domain.Question{ID: questionID})
	if err != nil {
		return err
	}

	return nil
}

func (r *testRepository) DetachQuestion(ctx context.Context, testID uint, questionID uint) error {
	err := r.db.
		Model(&domain.Test{ID: testID}).
		Association("Questions").
		Delete(&domain.Question{ID: questionID})
	if err != nil {
		return err
	}
	return nil
}

func (r *testRepository) CreateUserTest(ctx context.Context, ut *domain.UserTests) error {
	// Убедимся, что ut.Status == domain.InProgress, Progress == 0
	ut.Status = domain.InProgress
	ut.Progress = 0

	return r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "test_id"}, {Name: "user_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status":   ut.Status,
				"progress": ut.Progress,
			}),
		}).
		Create(ut).
		Error
}

// CreateWithUser - создаём тест, цепляем к курсу и сразу добавляем
// запись в user_tests.
func (r *testRepository) CreateWithUser(
	ctx context.Context,
	test *domain.Test,
	courseID uint,
	userID uint,
) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. сам тест
		if err := tx.Create(test).Error; err != nil {
			return err
		}
		// 2. связь «курс-тест»
		if err := tx.Model(test).
			Association("Courses").
			Append(&domain.Course{ID: courseID}); err != nil {
			return err
		}
		// 3. запись в user_tests
		ut := &domain.UserTests{
			UserID: userID,
			TestID: test.ID,
			Status: domain.New,
		}
		if err := tx.Create(ut).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *testRepository) UpdateUserTest(ctx context.Context, userTest *domain.UserTests) error {
	updates := validator.BuildUpdates(userTest)

	result := r.db.Model(&domain.UserTests{}).Where("user_id = ? AND test_id = ?", userTest.UserID, userTest.TestID).Updates(updates).First(&userTest)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.ErrTestNotFound
		}

		return result.Error
	}

	return nil
}

func (r *testRepository) GetCourseIDByTestID(ctx context.Context, testID uint) (uint, error) {
	var tst domain.Test
	if err := r.db.WithContext(ctx).
		Preload("Courses").
		First(&tst, testID).
		Error; err != nil {
		return 0, err
	}
	if len(tst.Courses) == 0 {
		return 0, fmt.Errorf("no course associated with test %d", testID)
	}
	return tst.Courses[0].ID, nil
}
