package asset_history

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

func (s *Store) GetAll() ([]AssetHistory, error) {
	var assetHistory []AssetHistory

	err := s.conn.Select(&assetHistory, "select id, data, create_date from asset_history")
	if err != nil {
		s.logger.Error("failed to get all asset history", err)
		return nil, err
	}

	return assetHistory, nil
}

func (s *Store) Insert(assetHistory *AssetHistory) error {
	sql := "insert into asset_history (id, data, create_date) values (:id, :data, :createDate)"

	params := map[string]interface{}{
		"id": utils.GetGuid(),
		"data": assetHistory.Data,
		"createDate": time.Now(),
	}

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