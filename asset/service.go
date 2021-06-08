package asset

import (
	"backend/auth"
	"backend/files"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service struct {
	store       *Store
	logger      hclog.Logger
	fileService *files.Service
}

func NewService(l hclog.Logger, s *Store, fileService *files.Service, e *echo.Echo) *Service {
	service := &Service{
		store:       s,
		logger:      l,
		fileService: fileService,
	}

	e.GET("/asset", service.GetAll, auth.IsLoggedIn)
	e.POST("/asset", service.Add, auth.IsLoggedIn)
	e.PATCH("/asset", service.Patch, auth.IsLoggedIn)
	e.DELETE("/asset", service.Delete, auth.IsLoggedIn)

	return service
}

func (s *Service) GetAll(c echo.Context) error {
	clients, err := s.store.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, clients)
}

func (s *Service) Add(c echo.Context) error {
	a := new(Asset)

	err := c.Bind(a)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Insert(a)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Patch(c echo.Context) error {
	a := new(Asset)

	err := c.Bind(a)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Update(a)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Delete(c echo.Context) error {
	a := new(Asset)

	err := c.Bind(a)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Delete(a)
	if err != nil {
		s.logger.Error("delete asset db failed", err)
		return c.String(http.StatusInternalServerError, "")
	}

	switch a.Type {
	case "audio", "picture":
		err = s.fileService.DeleteFile(a.Link)
		if err != nil {
			s.logger.Error("delete asset file failed", err)
			return c.String(http.StatusInternalServerError, "")
		}
	}

	return nil
}
