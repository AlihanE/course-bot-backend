package sender

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/go-hclog"
)

type Service struct {
	bot    *botApi.BotAPI
	logger hclog.Logger
}

func New(bot *botApi.BotAPI, logger hclog.Logger) *Service {
	return &Service{
		bot:    bot,
		logger: logger,
	}
}

func (s *Service) SendText(chatId int64, text string) error {
	msg := botApi.NewMessage(chatId, text)

	return s.send(msg)
}

func (s *Service) SendPhoto(chatId int64, filePath string) error {
	msg := botApi.NewPhotoUpload(chatId, filePath)

	return s.send(msg)
}

func (s *Service) send(c botApi.Chattable) error {
	_, err := s.bot.Send(c)
	if err != nil {
		s.logger.Error("send message failed", err)
		return err
	}

	return nil
}

func (s *Service) SendAudio(chatId int64, filePath string) error {
	msg := botApi.NewAudioUpload(chatId, filePath)

	return s.send(msg)
}
