package user_step

import "time"

type UserStep struct {
	Id         string    `json:"id" db:"id"`
	UserId     string    `json:"userId" db:"user_id"`
	StepId     string    `json:"stepId" db:"step_id"`
	Completed  bool      `json:"completed" db:"completed"`
	Sent       bool      `json:"sent" db:"sent"`
	CreateDate time.Time `json:"createDate" db:"create_date"`
	UpdateDate time.Time `json:"updateDate" db:"update_date"`
}
