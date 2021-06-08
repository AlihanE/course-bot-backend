package client

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

	e.GET("/client", service.getAllApi, auth.IsLoggedIn)
	e.GET("/client/:login", service.getClientApi, auth.IsLoggedIn)
	e.POST("/client", service.addApi, auth.IsLoggedIn)
	e.PATCH("/client", service.patchApi, auth.IsLoggedIn)

	return service
}

func (s *Service) getClientApi(c echo.Context) error {
	login := c.Param("login")

	clients, err := s.GetClient(login)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	if len(clients) == 0 {
		return c.String(http.StatusNotFound, "")
	}
	return c.JSON(http.StatusOK, clients[0])
}

func (s *Service) GetClient(chatId string) ([]Client, error) {
	clients, err := s.store.GetClient(chatId)
	if err != nil {
		s.logger.Error("get client by login failed", err)
		return nil, err
	}

	return clients, err
}

func (s *Service) GetClientById(id string) ([]Client, error) {
	clients, err := s.store.GetById(id)
	if err != nil {
		s.logger.Error("get client by login failed", err)
		return nil, err
	}

	return clients, err
}

func (s *Service) getAllApi(c echo.Context) error {
	clients, err := s.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, clients)
}

func (s *Service) GetAll() ([]Client, error) {
	clients, err := s.store.GetAll()
	if err != nil {
		s.logger.Error("get clients failed", err)
		return nil, err
	}

	return clients, err
}

func (s *Service) GetAllActive() ([]Client, error) {
	clients, err := s.store.GetAllActive()
	if err != nil {
		s.logger.Error("get clients failed", err)
		return nil, err
	}

	return clients, err
}

func (s *Service) addApi(c echo.Context) error {
	client := new(Client)

	err := c.Bind(client)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.Add(client)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Add(c *Client) error {
	err := s.store.Insert(c)
	if err != nil {
		s.logger.Error("insert client failed", err)
		return err
	}

	return nil
}

func (s *Service) patchApi(c echo.Context) error {
	client := new(Client)

	err := c.Bind(client)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Update(client)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Patch(c *Client) error {
	err := s.store.Update(c)
	if err != nil {
		s.logger.Error("update client failed", err)
		return err
	}

	return nil
}
