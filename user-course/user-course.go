package user_course

type UserCourse struct {
	Id       string `json:"id" db:"id"`
	UserId   string `json:"userId" db:"user_id"`
	CourseId string `json:"courseId" db:"course_id"`
	Finished bool   `json:"finished" db:"finished"`
}
