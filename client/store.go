package client

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

func (s *Store) GetClient(chatId string) ([]Client, error) {
	var client []Client

	err := s.conn.Select(&client, "select id, first_name, last_name, login, chat_id, active from clients where chat_id = $1", chatId)
	if err != nil {
		s.logger.Error("failed to get client by chatId", chatId, "error", err)
		return nil, err
	}

	return client, nil
}

func (s *Store) GetAll() ([]Client, error) {
	var clients []Client

	err := s.conn.Select(&clients, "select id, first_name, last_name, login, chat_id, active from clients")
	if err != nil {
		s.logger.Error("failed to get all clients", err)
		return nil, err
	}

	return clients, nil
}

func (s *Store) GetById(id string) ([]Client, error) {
	var client []Client

	err := s.conn.Select(&client, "select id, first_name, last_name, login, chat_id, active from clients where id = $1", id)
	if err != nil {
		s.logger.Error("failed to get client by Id", id, "error", err)
		return nil, err
	}

	return client, nil
}

func (s *Store) GetAllActive() ([]Client, error) {
	var clients []Client

	err := s.conn.Select(&clients, "select id, first_name, last_name, login, chat_id, active from clients where active=true")
	if err != nil {
		s.logger.Error("failed to get all active clients", err)
		return nil, err
	}

	return clients, nil
}

func mapClientToMap(cli *Client) map[string]interface{} {
	return map[string]interface{}{
		"id":        cli.Id,
		"firstName": cli.FirstName,
		"lastName":  cli.LastName,
		"login":     cli.Login,
		"chatId":    cli.ChatId,
		"active":    cli.Active,
	}
}

func (s *Store) execute(sql string, params map[string]interface{}) error {
	res, err := s.conn.NamedExec(sql, params)
	if err != nil {
		s.logger.Error("failed to execute sql for client", err)
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

func (s *Store) Update(cli *Client) error {

	sql := `UPDATE clients set first_name=:firstName, last_name=:lastName, ` +
		`login=:login, chat_id=:chatId, active=:active, update_date=:updateDate where id=:id`

	params := mapClientToMap(cli)
	params["updateDate"] = time.Now()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Insert(cli *Client) error {
	sql := `INSERT into clients (id, first_name, last_name, login, chat_id, active, create_date) values ` +
		`(:id, :firstName, :lastName, :login, :chatId, :active, :createDate)`

	params := mapClientToMap(cli)
	params["createDate"] = time.Now()
	params["id"] = utils.GetGuid()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}
