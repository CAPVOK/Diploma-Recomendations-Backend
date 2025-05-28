package test

import (
	"context"
	"diprec_api/internal/domain"
	"diprec_api/internal/infrastructure/kafka"
	"diprec_api/internal/repository/test"
	"strconv"

	"go.uber.org/zap"
)

type testUsecase struct {
	repo     test.ITestRepository
	producer kafka.IKafkaProducer
	logger   *zap.Logger
}

type ITestUsecase interface {
	Create(ctx context.Context, test *domain.Test, courseID uint) (*domain.Test, error)
	Get(ctx context.Context, courseID, userID uint) ([]*domain.Test, error)
	GetByID(ctx context.Context, id, userID uint) (*domain.Test, error)
	Update(ctx context.Context, test *domain.Test) (*domain.Test, error)
	Delete(ctx context.Context, id uint) error
	AttachQuestion(ctx context.Context, testID uint, questionID uint) error
	DetachQuestion(ctx context.Context, testID uint, questionID uint) error
	StartTest(ctx context.Context, userTests *domain.UserTests) error
	EndTest(ctx context.Context, userTest *domain.UserTests) error
	CreateRecommendTest(
		ctx context.Context,
		test *domain.Test,
		courseID uint,
		userID uint,
		questionIDs []uint,
	) (*domain.Test, error)
}

func NewTestUsecase(repo test.ITestRepository, producer kafka.IKafkaProducer, logger *zap.Logger) ITestUsecase {
	return &testUsecase{repo, producer, logger.Named("TestUsecase")}
}

func (u *testUsecase) Create(ctx context.Context, test *domain.Test, courseID uint) (*domain.Test, error) {
	if err := u.repo.Create(ctx, test, courseID); err != nil {
		return nil, err
	}

	return test, nil
}

func (u *testUsecase) Get(ctx context.Context, courseID, userID uint) ([]*domain.Test, error) {
	tests, err := u.repo.Get(ctx, courseID, userID)
	if err != nil {
		return nil, err
	}

	return tests, nil
}

func (u *testUsecase) GetByID(ctx context.Context, id, userID uint) (*domain.Test, error) {
	test, err := u.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (u *testUsecase) Update(ctx context.Context, test *domain.Test) (*domain.Test, error) {
	if err := u.repo.Update(ctx, test); err != nil {
		return nil, err
	}

	return test, nil
}

func (u *testUsecase) Delete(ctx context.Context, id uint) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (u *testUsecase) AttachQuestion(ctx context.Context, testID uint, questionID uint) error {
	if err := u.repo.AttachQuestion(ctx, testID, questionID); err != nil {
		return err
	}

	return nil
}

func (u *testUsecase) DetachQuestion(ctx context.Context, testID uint, questionID uint) error {
	if err := u.repo.DetachQuestion(ctx, testID, questionID); err != nil {
		return err
	}

	return nil
}

func (u *testUsecase) StartTest(ctx context.Context, userTests *domain.UserTests) error {
	if err := u.repo.CreateUserTest(ctx, userTests); err != nil {
		return err
	}

	return nil
}

func (u *testUsecase) CreateRecommendTest(
	ctx context.Context,
	test *domain.Test,
	courseID uint,
	userID uint,
	questionIDs []uint,
) (*domain.Test, error) {

	test.Assignee = domain.Recommendation
	test.Status = domain.Progress

	if err := u.repo.CreateWithUser(ctx, test, courseID, userID); err != nil {
		return nil, err
	}
	// прикрепляем вопросы
	for _, qid := range questionIDs {
		if err := u.repo.AttachQuestion(ctx, test.ID, qid); err != nil {
			return nil, err
		}
	}
	return test, nil
}

func (u *testUsecase) EndTest(ctx context.Context, userTest *domain.UserTests) error {
	if err := u.repo.UpdateUserTest(ctx, userTest); err != nil {
		return err
	}

	test, err := u.repo.GetByID(ctx, userTest.TestID, userTest.UserID)
	if err != nil {
		return err
	}

	if test.Assignee == domain.Recommendation {
		return nil
	}

	courseID, err := u.repo.GetCourseIDByTestID(ctx, userTest.TestID)
	if err != nil {
		u.logger.Warn("cannot lookup course for test", zap.Uint("testID", userTest.TestID), zap.Error(err))
	} else {
		msg := map[string]interface{}{
			"user_id":   int(userTest.UserID),
			"test_id":   int(userTest.TestID),
			"course_id": int(courseID),
		}
		_ = u.producer.Send(
			ctx,
			domain.TopicUserTest,
			strconv.Itoa(int(userTest.UserID)),
			msg,
		)
	}

	return nil
}
