package question

import (
	"context"
	"diprec_api/internal/domain"
	"diprec_api/internal/infrastructure/kafka"
	"diprec_api/internal/pkg/utils"
	"diprec_api/internal/repository/question"
	"diprec_api/internal/repository/test"

	"strconv"
	"time"

	"go.uber.org/zap"
)

type questionUsecase struct {
	repo     question.IQuestionRepository
	testRepo test.ITestRepository
	producer kafka.IKafkaProducer
	logger   *zap.Logger
}

type IQuestionUsecase interface {
	Create(ctx context.Context, question *domain.Question) (*domain.Question, error)
	GetAll(ctx context.Context) ([]*domain.Question, error)
	GetByID(ctx context.Context, id uint) (*domain.Question, error)
	Update(ctx context.Context, question *domain.Question) (*domain.Question, error)
	Delete(ctx context.Context, id uint) error
	Check(ctx context.Context, id, userID uint, answer interface{}, testId int) (*domain.QuestionAnswer, error)
}

func NewQuestionUsecase(repo question.IQuestionRepository, testRepo test.ITestRepository, producer kafka.IKafkaProducer, logger *zap.Logger) IQuestionUsecase {
	return &questionUsecase{repo, testRepo, producer, logger}
}

func (u *questionUsecase) Create(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	if err := u.repo.Create(ctx, question); err != nil {
		return nil, err
	}

	_ = u.producer.Send(
		ctx,
		domain.TopicCreateQuestion,
		strconv.Itoa(int(question.ID)),
		question,
	)

	return question, nil
}

func (u *questionUsecase) GetAll(ctx context.Context) ([]*domain.Question, error) {
	questions, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (u *questionUsecase) GetByID(ctx context.Context, id uint) (*domain.Question, error) {
	question, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return question, nil
}

func (u *questionUsecase) Update(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	if err := u.repo.Update(ctx, question); err != nil {
		return nil, err
	}

	_ = u.producer.Send(
		ctx,
		domain.TopicEditQuestion,
		strconv.Itoa(int(question.ID)),
		question,
	)

	return question, nil
}

func (u *questionUsecase) Delete(ctx context.Context, id uint) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	msg := map[string]interface{}{
		"id": id,
	}
	_ = u.producer.Send(
		ctx,
		domain.TopicDeleteQuestion,
		strconv.Itoa(int(id)),
		msg,
	)

	return nil
}

func (u *questionUsecase) Check(ctx context.Context, id, userID uint, answer interface{}, testId int) (*domain.QuestionAnswer, error) {
	question, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	isCorrect := question.CheckAnswer(answer)

	courseID, err := u.testRepo.GetCourseIDByTestID(ctx, uint(testId))
	if err != nil {
		u.logger.Warn("cannot lookup course for test", zap.Uint("testID", id), zap.Error(err))
		courseID = 0
	} else {
		// --- добавляем проверку на рекомендованный тест ---
		test, err := u.testRepo.GetByID(ctx, uint(testId), userID)
		if err != nil {
			u.logger.Warn("cannot get test by ID", zap.Uint("testID", uint(testId)), zap.Error(err))
		} else if test.Assignee != domain.Recommendation {
			userAnswer := &domain.UserAnswerCheck{
				QuestionID: id,
				CourseID:   courseID,
				Title:      question.Title,
				Type:       question.Type.String(),
				Variants:   question.Variants,
				Answer:     question.Answer,
				IsCorrect:  isCorrect,
				UserID:     userID,
				Timestamp:  time.Now(),
			}

			_ = u.producer.Send(
				ctx,
				domain.TopicUserAnswers,
				strconv.Itoa(int(userID)),
				userAnswer,
			)
		}
	}

	return &domain.QuestionAnswer{
		IsCorrect: isCorrect,
		Message:   utils.GenerateFeedbackMessage(isCorrect),
		Answer:    question.Answer,
	}, nil
}
