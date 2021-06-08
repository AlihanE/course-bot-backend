package asset

type Asset struct {
	Id       string `json:"id" db:"id"`
	CourseId string `json:"courseId" db:"course_id"`
	StepId   string `json:"stepId" db:"step_id"`
	Link     string `json:"link" db:"link"`
	Text     string `json:"text" db:"text"`
	Picture  string `json:"picture" db:"picture"`
	Data     string `json:"data,omitempty"`
	Type     string `json:"type" db:"type"`
}
