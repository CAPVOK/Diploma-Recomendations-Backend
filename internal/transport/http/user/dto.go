package user

type CreateUserDTO struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required,min=8"`
	FirstName  string `json:"firstName" binding:"required"`
	LastName   string `json:"lastName" binding:"required"`
	Patronymic string `json:"patronymic,omitempty"`
}

type LoginUserDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type RefreshUserDTO struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
