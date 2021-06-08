package ourbot

import (
	asset_history "backend/asset-history"
	"backend/auth"
	"backend/client"
	clientConversation "backend/client-conversation"
	"backend/ourbot/sender"
	"encoding/json"
	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Service struct {
	clientService    *client.Service
	botSender        *sender.Service
	logger           hclog.Logger
	clientConService *clientConversation.Service
	assetHistoryService *asset_history.Service
}

type SendData struct {
	Assets []Data `json:"assets"`
	Clients []string `json:"clients"`
}

type Data struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Link     string `json:"link"`
	FilePath string `json:"filePath"`
}

type ClientConversation struct {
	UserId   string `json:"userId"`
	ReportId string `json:"reportId"`
	Text     string `json:"text"`
}

func New(cs *client.Service, b *sender.Service, l hclog.Logger, e *echo.Echo,
	clientConv *clientConversation.Service, ah *asset_history.Service) *Service {
	service := &Service{
		clientService:    cs,
		botSender:        b,
		logger:           l,
		clientConService: clientConv,
		assetHistoryService: ah,
	}

	e.POST("/to-bot", service.sendToClients, auth.IsLoggedIn)
	e.POST("/client-conversation", service.converse, auth.IsLoggedIn)

	return service
}

func (s *Service) converse(c echo.Context) error {
	var conversationData ClientConversation

	err := c.Bind(&conversationData)
	if err != nil {
		s.logger.Error("request parse failed", err)
		return c.String(http.StatusBadRequest, "")
	}

	if len(conversationData.Text) == 0 {
		return c.String(http.StatusBadRequest, "")
	}

	clients, err := s.clientService.GetClientById(conversationData.UserId)
	if err != nil {
		s.logger.Error("request client failed", err)
		return c.String(http.StatusInternalServerError, "")
	}

	if len(clients) == 0 {
		s.logger.Error("client not found", err)
		return c.String(http.StatusInternalServerError, "")
	}

	for _, clientData := range clients {
		chatId, err := strconv.ParseInt(clientData.ChatId, 10, 64)
		if err != nil {
			s.logger.Error("failed to parse clients", clientData.Id, "chat id", clientData.ChatId)
			continue
		}

		err = s.botSender.SendText(chatId, conversationData.Text)
		if err != nil {
			s.logger.Error("failed to send to client", clientData.Id, "chat id", clientData.ChatId)
			continue
		}

		cliConv := &clientConversation.ClientConversation{
			UserId:   conversationData.UserId,
			ReportId: conversationData.ReportId,
			Text:     conversationData.Text,
		}

		err = s.clientConService.Add(cliConv)
		if err != nil {
			s.logger.Error("failed to save message for client", clientData.Id, "chat id", clientData.ChatId)
			continue
		}
	}

	return c.String(http.StatusOK, "")
}

func (s *Service) sendToClients(c echo.Context) error {
	var data SendData

	err := c.Bind(&data)
	if err != nil {
		s.logger.Error("request parse failed", err)
		return c.String(http.StatusBadRequest, "")
	}

	for _, client := range data.Clients {
		chatId, err := strconv.ParseInt(client, 10, 64)
		if err != nil {
			s.logger.Error("failed to parse clients chat id", client)
			continue
		}

		for _, piece := range data.Assets {
			switch piece.Type {
			case "text", "link":
				if len(piece.Text) > 0 {
					err = s.botSender.SendText(chatId, piece.Text)
					if err != nil {
						s.logger.Error("failed to send to client chat id", client, "data with type", piece.FilePath)
						continue
					}
				} else if len(piece.Link) > 0 {
					err = s.botSender.SendText(chatId, piece.Link)
					if err != nil {
						s.logger.Error("failed to send to client chat id", client, "data with type", piece.FilePath)
						continue
					}
				}
			case "picture":
				err = s.botSender.SendPhoto(chatId, piece.FilePath)
				if err != nil {
					s.logger.Error("failed to send to client chat id", client, "data with type", piece.FilePath)
					continue
				}
			case "audio":
				err = s.botSender.SendAudio(chatId, piece.FilePath)
				if err != nil {
					s.logger.Error("failed to send to client chat id", client, "data with type", piece.FilePath)
					continue
				}
			}
		}

	}

	js, err := json.Marshal(data)
	if err == nil {
		assetHistory := &asset_history.AssetHistory{
			Data:       string(js),
		}

		err = s.assetHistoryService.Add(assetHistory)
		if err != nil {
			s.logger.Error("insert asset history failed", err)
		}
	} else {
		s.logger.Error("asset history marshal failed", err)
	}

	return c.String(http.StatusOK, "")
}
