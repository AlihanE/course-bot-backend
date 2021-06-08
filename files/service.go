package files

import (
	"backend/auth"
	"backend/utils"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
)

var (
	dirName = "/home/alikhan/sendfiles"
)

type Service struct {
	logger hclog.Logger
}

func NewService(l hclog.Logger, e *echo.Echo) (*Service, error) {
	if utils.IsTestEnvironment() {
		dirName = "./sendfiles"
	}

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.Mkdir(dirName, os.ModeDir)
		if err != nil {
			l.Error("failed to create dir", dirName, "error", err)
			return nil, err
		}
	}

	service := &Service{
		logger: l,
	}

	e.POST("/file", service.UploadFile, auth.IsLoggedIn)

	return service, nil
}

func (s *Service) UploadFile(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		s.logger.Error("file retrieve error", err)
		return c.String(http.StatusInternalServerError, "")
	}

	src, err := file.Open()
	if err != nil {
		s.logger.Error("create file error", err)
	}

	defer src.Close()

	filePath := dirName + "/" + file.Filename

	dst, err := os.Create(filePath)
	if err != nil {
		s.logger.Error("create destination file error", err)
		return c.String(http.StatusInternalServerError, "")
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		s.logger.Error("file write error", err)
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, `{"filePath": "`+filePath+`"}`)
}

func (s *Service) DeleteFile(filePath string) error {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if !info.IsDir() {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}
