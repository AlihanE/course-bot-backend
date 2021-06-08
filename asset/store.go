package asset

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

func (s *Store) GetAll() ([]Asset, error) {
	var assets []Asset

	err := s.conn.Select(&assets, "select id, course_id, step_id, link, text, picture from assets")
	if err != nil {
		s.logger.Error("failed to get all assets", err)
		return nil, err
	}

	return assets, nil
}

func mapAssetToMao(asset *Asset) map[string]interface{} {
	return map[string]interface{}{
		"id":       asset.Id,
		"courseId": asset.CourseId,
		"stepId":   asset.StepId,
		"link":     asset.Link,
		"text":     asset.Text,
		"picture":  asset.Picture,
		"type":     asset.Type,
	}
}

func (s *Store) execute(sql string, params map[string]interface{}) error {
	res, err := s.conn.NamedExec(sql, params)
	if err != nil {
		s.logger.Error("failed to execute sql for asset", err)
		return err
	}

	if cnt, err := res.RowsAffected(); err != nil {
		if cnt <= 0 {
			return errors.New("assets no rows affected")
		}
	} else {
		return err
	}

	return nil
}

func (s *Store) Update(asset *Asset) error {
	sql := `UPDATE assets set course_id=:courseId, step_id=:stepId, link=:link,
		text=:text, picture=:picture, update_date=:updateDate type=:type where id=:id`

	params := mapAssetToMao(asset)

	params["updateDate"] = time.Now()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Insert(asset *Asset) error {
	sql := `INSERT into assets (id, course_id, step_id, link, text, picture,
		create_date, type) values (:id, :courseId, :stepId, :link, :text, :picture, :createDate, :type)`

	params := mapAssetToMao(asset)
	params["createDate"] = time.Now()
	params["id"] = utils.GetGuid()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Delete(asset *Asset) error {
	sql := `DELETE from assets where id = :id`

	params := mapAssetToMao(asset)

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}
