package asset_history

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

	e.GET("/asset-history", service.getAssetHistoryApi, auth.IsLoggedIn)

	return service
}

func (s *Service) getAssetHistoryApi(c echo.Context) error {
	assetHis, err := s.store.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, assetHis)
}

func (s *Service) Add(assetHistory *AssetHistory) error {
	err := s.store.Insert(assetHistory)
	if err != nil {
		s.logger.Error("insert asset history failed", err)
		return err
	}

	return nil
}
