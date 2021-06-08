package client_conversation

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

func (s *Store) GetReportConversation(reportId string) ([]ClientConversation, error) {
	var clientConv []ClientConversation

	err := s.conn.Select(&clientConv, "select id, user_id, report_id, text, create_date from client_report_conversation where report_id = $1", reportId)
	if err != nil {
		s.logger.Error("failed to get conversation by reportId", reportId, "error", err)
		return nil, err
	}

	return clientConv, nil
}

func mapStepToMap(clientConv *ClientConversation) map[string]interface{} {
	return map[string]interface{}{
		"id":       clientConv.Id,
		"userId":   clientConv.UserId,
		"reportId": clientConv.ReportId,
		"text":     clientConv.Text,
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

func (s *Store) Update(clientConv *ClientConversation) error {
	sql := `UPDATE steps set user_id=:userId, reportId=:report_id, text=:text, update_date=:updateDate, ` +
		`where id=:id`

	params := mapStepToMap(clientConv)
	params["updateDate"] = time.Now()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Insert(clientConv *ClientConversation) error {
	sql := `INSERT into client_report_conversation (id, user_id, report_id, text, create_date) ` +
		`values (:id, :userId, :reportId, :text, :createDate)`

	params := mapStepToMap(clientConv)
	params["createDate"] = time.Now()
	params["id"] = utils.GetGuid()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}
