package user_step

import (
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service struct {
	store  *Store
	logger hclog.Logger
}

func NewService(l hclog.Logger, s *Store) *Service {
	return &Service{
		store:  s,
		logger: l,
	}
}

func (s *Service) GetAllActive(c echo.Context) error {
	userSteps, err := s.store.GetAllActive()

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, userSteps)
}

func (s *Service) Add(c echo.Context) error {
	userSteps := new(UserStep)

	err := c.Bind(userSteps)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Insert(userSteps)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Patch(c echo.Context) error {
	userStep := new(UserStep)

	err := c.Bind(userStep)
	if err != nil {
		c.String(http.StatusBadRequest, "")
	}

	err = s.store.Update(userStep)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}
