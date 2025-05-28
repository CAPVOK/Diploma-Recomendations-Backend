package domain

import "errors"

type Error struct {
	Message string
}

var (
	/* common */
	ErrInternalServer     = errors.New("Внутренняя ошибка сервера")
	ErrUnauthorized       = errors.New("Неавторизован")
	ErrInvalidRequestBody = errors.New("Данные неверны")
	/* auth */
	ErrInvalidCredentials  = errors.New("Неверные имя пользователя или пароль")
	ErrUserExists          = errors.New("Пользователь с таким именем уже существует")
	ErrUserNotFound        = errors.New("Пользователь не найден")
	ErrInvalidRefreshToken = errors.New("Рефреш токен невалиден")
	ErrInvalidTokenType    = errors.New("Неверный тип токена")
	ErrInvalidRole         = errors.New("Данный функционал доступен только преподавателю!")
	/* course */
	ErrCourseNotFound = errors.New("Курс не найден")
	/* test */
	ErrTestNotFound = errors.New("Тест не найден")
	/* question */
	ErrQuestionNotFound = errors.New("Вопрос не найден")
)
