package user_course

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

func (s *Store) GetAll() ([]UserCourse, error) {
	var userCourses []UserCourse

	err := s.conn.Select(&userCourses, "select * from user_courses")
	if err != nil {
		s.logger.Error("failed to get all user courses")
		return nil, err
	}

	return userCourses, nil
}

func mapUserCourseToMap(uc *UserCourse) map[string]interface{} {
	return map[string]interface{}{
		"id":        uc.Id,
		"user_id":   uc.UserId,
		"course_id": uc.CourseId,
		"finished":  uc.Finished,
	}
}

func (s *Store) execute(sql string, params map[string]interface{}) error {
	res, err := s.conn.NamedExec(sql, params)
	if err != nil {
		s.logger.Error("failed to execute sql for user course", err)
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

func (s *Store) Update(uc *UserCourse) error {
	sql := `UPDATE user_courses set user_id=:user-id, course-id=:course-id,
		finished=:finished, update_date=:update-date where id=:id`

	params := mapUserCourseToMap(uc)
	params["update-date"] = time.Now()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Insert(uc *UserCourse) error {
	sql := `INSERT into user_courses (id, user_id, course_id, finished, create_date)
		values (:id, :user-id, :course-id, :finished, :create-date)`

	params := mapUserCourseToMap(uc)
	params["create-date"] = time.Now()
	params["id"] = utils.GetGuid()

	err := s.execute(sql, params)
	if err != nil {
		return err
	}

	return nil
}
