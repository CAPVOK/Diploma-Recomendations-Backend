package question

type CreateQuestionDTO struct {
	Title    string                 `json:"title"`
	Type     string                 `json:"type" enums:"SINGLE,TEXT,NUMBER,MULTIPLE" example:"SINGLE"`
	Variants map[string]interface{} `json:"variants"`
	Answer   interface{}            `json:"answer"`
}

type UpdateQuestionDTO struct {
	Title    string                 `json:"title"`
	Type     string                 `json:"type" enums:"SINGLE,TEXT,NUMBER,MULTIPLE" example:"SINGLE"`
	Variants map[string]interface{} `json:"variants"`
	Answer   interface{}            `json:"answer"`
}

type CheckAnswerDTO struct {
	Answer interface{} `json:"answer"`
	TestId int         `json:"testId"`
}
