package user_step

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

func (s *Store) GetAllActive() ([]UserStep, error) {
	var userSteps []UserStep

	err := s.conn.Select(&userSteps, "select * from user_step where completed = false and sent = false")
	if err != nil {
		s.logger.Error("failed to get all user courses")
		return nil, err
	}

	return userSteps, nil
}

func mapUserStepToMap(uc *UserStep) map[string]interface{} {
	return map[string]interface{}{
		"id":        uc.Id,
		"userId":    uc.UserId,
		"stepId":    uc.StepId,
		"completed": uc.Completed,
		"sent":      uc.Sent,
	}
}

func (s *Store) execute(sql string, params map[string]interface{}) error {
	res, err := s.conn.NamedExec(sql, params)
	if err != nil {
		s.logger.Error("failed to execute sql for user step", err)
		return err
	}

	if cnt, err := res.RowsAffected(); err != nil {
		if cnt <= 0 {
			return errors.New("user courses no rows affected")
		}
	} else {
		return err
	}

	return nil
}

func (s *Store) Update(us *UserStep) error {
	sql := `UPDATE user_step set user_id =:userId, step_id =:stepId, completed=:completed, sent=:sent, 
		update_date =:updateDate where id =:id`

	params := mapUserStepToMap(us)
	params["update-date"] = time.Now()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Insert(us *UserStep) error {
	sql := `INSERT into user_step (id, user_id, step_id, completed, sent, create_date) 
		values(:id, :userId, :stepId, :completed, :sent, :createDate)`

	params := mapUserStepToMap(us)
	params["createDate"] = time.Now()
	params["id"] = utils.GetGuid()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}
