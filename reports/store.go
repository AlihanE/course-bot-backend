package reports

import (
	"backend/utils"
	"errors"
	"github.com/hashicorp/go-hclog"
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	conn   *sqlx.DB
	logger hclog.Logger
}

func NewStore(l hclog.Logger, c *sqlx.DB) *Store {
	return &Store{
		conn:   c,
		logger: l,
	}
}

func (s *Store) GetAll() ([]Report, error) {
	var reports []Report

	err := s.conn.Select(&reports, "select id, user_id, course_id, step_id, text, accepted, create_date from reports")
	if err != nil {
		s.logger.Error("failed to get all reports", err)
		return nil, err
	}

	return reports, nil
}

func mapReportToMap(report *Report) map[string]interface{} {
	return map[string]interface{}{
		"id":       report.Id,
		"userId":   report.UserId,
		"courseId": report.CourseId,
		"stepId":   report.StepId,
		"text":     report.Text,
		"accepted": report.Accepted,
	}
}

func (s *Store) execute(sql string, params map[string]interface{}) error {
	res, err := s.conn.NamedExec(sql, params)
	if err != nil {
		s.logger.Error("failed to execute sql for report", err)
		return err
	}

	if cnt, err := res.RowsAffected(); err != nil {
		if cnt <= 0 {
			return errors.New("clients no rows affected")
		}
	} else {
		return err
	}

	return nil
}

func (s *Store) Update(report *Report) error {

	sql := `UPDATE reports set user_id=:userId, course_id=:courseId, step_id=:stepId,
		text=:text, accepted=:accepted, update_date=:updateDate where id=:id`

	params := mapReportToMap(report)
	params["updateDate"] = time.Now()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Insert(report *Report) error {
	sql := `INSERT into reports (id, user_id, course_id, step_id, text, accepted, create_date)
		values (:id, :userId, :courseId, :stepId, :text, :accepted, :createDate)`

	params := mapReportToMap(report)
	params["createDate"] = time.Now()
	params["id"] = utils.GetGuid()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}
