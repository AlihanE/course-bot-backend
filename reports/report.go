package reports

import "time"

type Report struct {
	Id         string    `json:"id" db:"id"`
	UserId     string    `json:"userId" db:"user_id"`
	CourseId   string    `json:"courseId" db:"course_id"`
	StepId     string    `json:"stepId" db:"step_id"`
	Text       string    `json:"text" db:"text"`
	Accepted   bool      `json:"accepted" db:"accepted"`
	CreateDate time.Time `json:"createDate" db:"create_date"`
}
