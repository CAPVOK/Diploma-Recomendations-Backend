package course

type CreateCourseDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCourseDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
