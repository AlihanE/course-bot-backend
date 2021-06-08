package reports

import (
	"backend/auth"
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

	e.GET("/report", service.getAllApi, auth.IsLoggedIn)
	e.POST("/report", service.addApi, auth.IsLoggedIn)
	e.PATCH("/report", service.patchApi, auth.IsLoggedIn)

	return service
}

func (s *Service) getAllApi(c echo.Context) error {
	reports, err := s.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, reports)
}

func (s *Service) GetAll() ([]Report, error) {
	reports, err := s.store.GetAll()
	if err != nil {
		s.logger.Error("get all reports failed", err)
		return nil, err
	}

	return reports, err
}

func (s *Service) addApi(c echo.Context) error {
	report := new(Report)

	err := c.Bind(report)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.Add(report)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Add(r *Report) error {
	err := s.store.Insert(r)
	if err != nil {
		s.logger.Error("insert report failed", err)
		return err
	}

	return nil
}

func (s *Service) patchApi(c echo.Context) error {
	report := new(Report)

	err := c.Bind(report)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.Patch(report)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Patch(r *Report) error {
	err := s.store.Update(r)
	if err != nil {
		s.logger.Error("update report failed", err)
		return err
	}

	return nil
}
