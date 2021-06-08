package step

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

func (s *Store) GetAll() ([]Step, error) {
	var steps []Step

	err := s.conn.Select(&steps, "select id, name, type, file, text, step_num from steps")
	if err != nil {
		s.logger.Error("failed to get all steps", err)
		return nil, err
	}

	return steps, nil
}

func mapStepToMap(step *Step) map[string]interface{} {
	return map[string]interface{}{
		"id":      step.Id,
		"name":    step.Name,
		"type":    step.Type,
		"file":    step.File,
		"text":    step.Text,
		"stepNum": step.StepNumber,
	}
}

func (s *Store) execute(sql string, params map[string]interface{}) error {
	res, err := s.conn.NamedExec(sql, params)
	if err != nil {
		s.logger.Error("failed to execute sql for step", err)
		return err
	}

	if cnt, err := res.RowsAffected(); err != nil {
		if cnt <= 0 {
			return errors.New("no rows affected")
		}
	} else {
		return err
	}

	return nil
}

func (s *Store) Update(step *Step) error {
	sql := `UPDATE steps set name=:name, type=:type, file=:file, text=:text, update_date=:updateDate, step_num=:stepNum ` +
		`where id=:id`

	params := mapStepToMap(step)
	params["updateDate"] = time.Now()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Insert(step *Step) error {
	sql := `INSERT into steps (id, name, type, file, text, create_date, step_num) ` +
		`values (:id, :name, :type, :file, :text, :createDate, :stepNum)`

	params := mapStepToMap(step)
	params["createDate"] = time.Now()
	params["id"] = utils.GetGuid()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}
