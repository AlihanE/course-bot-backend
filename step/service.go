package step

import (
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service struct {
	store  *Store
	logger hclog.Logger
}

func NewService(l hclog.Logger, s *Store, e *echo.Echo) *Service {
	service := &Service{
		store:  s,
		logger: l,
	}

	e.GET("/step", service.GetAll)
	e.POST("/step", service.Add)
	e.PATCH("/step", service.Patch)

	return service
}

func (s *Service) GetAll(c echo.Context) error {
	steps, err := s.store.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, steps)
}

func (s *Service) Add(c echo.Context) error {
	step := new(Step)

	err := c.Bind(step)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Insert(step)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Patch(c echo.Context) error {
	step := new(Step)

	err := c.Bind(step)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Update(step)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}
