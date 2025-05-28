package domain

type UserCourse struct {
	UserID   uint `gorm:"primary_key"`
	CourseID uint `gorm:"primary_key"`
}

type CourseTest struct {
	CourseID uint `gorm:"primary_key"`
	TestID   uint `gorm:"primary_key"`
}

type TestQuestion struct {
	TestID     uint `gorm:"primary_key"`
	QuestionID uint `gorm:"primary_key"`
}

type UserTests struct {
	TestID   uint `gorm:"primary_key"`
	UserID   uint `gorm:"primary_key"`
	Progress uint
	Status   UserTestStatus `gorm:"type:varchar(20);not null;status IN ('IN_PROGRESS', 'ENDED', 'REC_NEW');default:'IN_PROGRESS'"`
}

type UserTestStatus string

const (
	New        UserTestStatus = "REC_NEW"
	InProgress UserTestStatus = "IN_PROGRESS"
	Completed  UserTestStatus = "COMPLETED"
)

func (ut UserTestStatus) String() string {
	return string(ut)
}

type UserTestResponse struct {
	Progress uint   `json:"progress"`
	Status   string `json:"status"`
}

func (ut *UserTests) ToUserTestResponse() UserTestResponse {
	return UserTestResponse{
		Progress: ut.Progress,
		Status:   ut.Status.String(),
	}
}
