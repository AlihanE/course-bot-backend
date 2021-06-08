package step

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
)

type Step struct {
	Id         string        `json:"id" db:"id"`
	Name       string        `json:"name" db:"name"`
	Type       string        `json:"type" db:"type"`
	File       string        `json:"file" db:"file"`
	Text       string        `json:"text" db:"text"`
	StepNumber sql.NullInt32 `json:"stepNumber" db:"step_num"`
}

func (s Step) MarshalJSON() ([]byte, error) {
	var stepNumber *int32
	if s.StepNumber.Valid {
		stepNumber = &s.StepNumber.Int32
	} else {
		stepNumber = nil
	}
	return json.Marshal(&struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		Type       string `json:"type"`
		File       string `json:"file"`
		Text       string `json:"text"`
		StepNumber *int32 `json:"stepNumber"`
	}{
		Id:         s.Id,
		Name:       s.Name,
		Type:       s.Type,
		File:       s.File,
		Text:       s.Text,
		StepNumber: stepNumber,
	})
}

func (s *Step) UnmarshalJSON(data []byte) error {
	var x map[string]interface{}
	log.Println(string(data))
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	s.Id, _ = x["id"].(string)
	s.Name, _ = x["name"].(string)
	s.Type, _ = x["type"].(string)
	s.File, _ = x["file"].(string)
	s.Text, _ = x["text"].(string)

	var stepNum string
	stepNum, _ = x["stepNumber"].(string)
	if len(stepNum) > 0 {
		if val, err := strconv.Atoi(stepNum); err == nil {
			s.StepNumber.Valid = true
			s.StepNumber.Int32 = int32(val)
		} else {
			s.StepNumber.Valid = false
		}
	}

	return nil
}
