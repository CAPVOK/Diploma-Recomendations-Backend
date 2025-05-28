package utils

import (
	"math/rand"
)

var encouragementMessages = []string{
	"Молодец!",
	"Отлично!",
	"Так держать!",
	"Ты справился!",
}

var supportMessages = []string{
	"Не переживай, получится в другой раз!",
	"Почти получилось, не сдавайся!",
	"Ошибки — это часть обучения!",
	"Давай попробуем ещё раз!",
}

func GenerateFeedbackMessage(correct bool) string {
	if correct {
		return encouragementMessages[rand.Intn(len(encouragementMessages))]
	}
	return supportMessages[rand.Intn(len(supportMessages))]
}
