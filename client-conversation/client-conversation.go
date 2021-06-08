package client_conversation

import "time"

type ClientConversation struct {
	Id         string    `json:"id" db:"id"`
	UserId     string    `json:"userId" db:"user_id"`
	ReportId   string    `json:"reportId" db:"report_id"`
	Text       string    `json:"text" db:"text"`
	CreateDate time.Time `json:"createDate" db:"create_date"`
	UpdateDate time.Time `json:"updateDate" db:"update_date"`
}
