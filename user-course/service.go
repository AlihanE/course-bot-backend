package user_course

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

	e.GET("/user-course", service.GetAll, auth.IsLoggedIn)
	e.POST("/user-course", service.Add, auth.IsLoggedIn)
	e.PATCH("/user-course", service.Patch, auth.IsLoggedIn)

	return service
}

func (s *Service) GetAll(c echo.Context) error {
	userCourses, err := s.store.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, userCourses)
}

func (s *Service) Add(c echo.Context) error {
	userCourse := new(UserCourse)

	err := c.Bind(userCourse)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Insert(userCourse)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) Patch(c echo.Context) error {
	userCourse := new(UserCourse)

	err := c.Bind(userCourse)
	if err != nil {
		return c.String(http.StatusBadRequest, "")
	}

	err = s.store.Update(userCourse)
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, "")
}
