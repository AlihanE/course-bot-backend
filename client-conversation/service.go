package client_conversation

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

	e.GET("/conversation/:reportId", service.getClientApi, auth.IsLoggedIn)
	e.POST("/conversation", service.addApi, auth.IsLoggedIn)
	e.PATCH("/conversation", service.patchApi, auth.IsLoggedIn)

	return service
}

func (s *Service) getClientApi(c echo.Context) error {
	reportId := c.Param("reportId")

	conv, err := s.store.GetReportConversation(reportId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, conv)
}

func (s *Service) addApi(c echo.Context) error {
	conv := new(ClientConversation)

	err := c.Bind(conv)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.Add(conv)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Add(c *ClientConversation) error {
	err := s.store.Insert(c)
	if err != nil {
		s.logger.Error("insert client failed", err)
		return err
	}

	return nil
}

func (s *Service) patchApi(c echo.Context) error {
	conv := new(ClientConversation)

	err := c.Bind(conv)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Update(conv)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}
