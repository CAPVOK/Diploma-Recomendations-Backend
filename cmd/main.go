package main

import (
	"diprec_api/cmd/application"
	"diprec_api/internal/config"
	"diprec_api/internal/infrastructure/db/postgres"
	"diprec_api/internal/infrastructure/kafka"
	"diprec_api/internal/pkg/logger"
	"diprec_api/internal/pkg/middleware"
	"diprec_api/internal/service"
	"fmt"
	"log"

	user_repo "diprec_api/internal/repository/user"
	user_handler "diprec_api/internal/transport/http/user"
	user_usecase "diprec_api/internal/usecase/user"

	course_repo "diprec_api/internal/repository/course"
	course_handler "diprec_api/internal/transport/http/course"
	course_usecase "diprec_api/internal/usecase/course"

	test_repo "diprec_api/internal/repository/test"
	test_handler "diprec_api/internal/transport/http/test"
	test_usecase "diprec_api/internal/usecase/test"

	question_repo "diprec_api/internal/repository/question"
	question_handler "diprec_api/internal/transport/http/question"
	question_usecase "diprec_api/internal/usecase/question"
)

func main() {
	cfg := config.MustLoad()

	fmt.Printf("Server starting on %s:%d\n", cfg.Server.Host, cfg.Server.Port)

	db, err := postgres.NewPostgresDB(postgres.Config{
		Host:            cfg.DB.Host,
		Port:            cfg.DB.Port,
		User:            cfg.DB.User,
		Password:        cfg.DB.Password,
		DBName:          cfg.DB.DBName,
		SSLMode:         cfg.DB.SSLMode,
		MaxIdleConns:    cfg.DB.MaxIdleConns,
		MaxOpenConns:    cfg.DB.MaxOpenConns,
		ConnMaxLifetime: cfg.DB.ConnMaxLifetime,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := postgres.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	custom_logger, err := logger.New(cfg.Logging)

	auth_service := service.NewAuthService(&service.JWTConfig{
		SecretKey:     cfg.Auth.JWTSecret,
		AccessExpiry:  cfg.Auth.AccessTokenExpire,
		RefreshExpiry: cfg.Auth.RefreshTokenExpire,
	})
	brokers := []string{cfg.KafkaProducer.Broker}
	kp := kafka.NewKafkaProducer(brokers, custom_logger)
	fmt.Printf("Kafka producer configured with brokers: %v\n", brokers)
	internalMW := middleware.Internal(cfg.InternalToken)
	fmt.Println("Internal token is:", cfg.InternalToken)
	ur := user_repo.NewUserRepository(db)
	uc := user_usecase.NewUserUseCase(ur, auth_service, custom_logger)
	uh := user_handler.NewUserHandler(uc, custom_logger)

	cr := course_repo.NewCourseRepository(db)
	cu := course_usecase.NewCourseUseCase(cr, custom_logger)
	ch := course_handler.NewCourseHandler(cu, custom_logger)

	tr := test_repo.NewTestRepository(db)
	tu := test_usecase.NewTestUsecase(tr, kp, custom_logger)
	th := test_handler.NewTestHandler(tu, custom_logger)

	qr := question_repo.NewQuestionRepository(db)
	qu := question_usecase.NewQuestionUsecase(qr, tr, kp, custom_logger)
	qh := question_handler.NewQuestionHandler(qu, custom_logger)

	app := application.NewApplication(cfg, custom_logger, db)

	app.Start(uh, ch, th, qh, auth_service, internalMW)
}
