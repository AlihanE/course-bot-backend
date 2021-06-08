package receiver

import (
	"backend/client"
	"backend/ourbot/sender"
	"backend/reports"
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hashicorp/go-hclog"
	"strconv"
	"sync"
)

type Receiver struct {
	reportService *reports.Service
	clientService *client.Service
	botSender     *sender.Service
	bot           *botApi.BotAPI
	logger        hclog.Logger
}

func New(b *botApi.BotAPI, l hclog.Logger, rs *reports.Service, cs *client.Service, botSender *sender.Service) *Receiver {
	return &Receiver{
		reportService: rs,
		clientService: cs,
		bot:           b,
		logger:        l,
		botSender:     botSender,
	}
}

func (r *Receiver) Start(wg *sync.WaitGroup) error {
	u := botApi.NewUpdate(0)
	u.Timeout = 60

	updates, err := r.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		for update := range updates {
			if update.Message == nil || update.Message.From.IsBot { // ignore any non-Message Updates
				continue
			}
			r.logger.Info("username", update.Message.From.UserName, "message", update.Message.Text)

			cliTmp := client.Client{
				FirstName: update.Message.From.FirstName,
				LastName:  update.Message.From.LastName,
				Login:     update.Message.From.UserName,
				ChatId:    strconv.FormatInt(update.Message.Chat.ID, 10),
				Active:    false,
			}

			cli, err := r.getClient(&cliTmp)
			if err != nil {
				r.logger.Error("getClient failed", err)
			}

			if cli.Active {
				rep := &reports.Report{
					UserId:   cli.Id,
					CourseId: "1",
					StepId:   "1",
					Text:     update.Message.Text,
					Accepted: false,
				}

				err := r.reportService.Add(rep)
				if err != nil {
					r.logger.Error("add report failed", err)
				}

				err = r.botSender.SendText(334558947, "Пришел отчет от пользователя с логином "+update.Message.From.UserName+
					" и ФИО "+update.Message.From.LastName+" "+update.Message.From.FirstName)
				if err != nil {
					r.logger.Error("send notification failed", err)
				}

				err = r.botSender.SendText(1193970058, "Пришел отчет от пользователя с логином "+update.Message.From.UserName+
					" и ФИО "+update.Message.From.LastName+" "+update.Message.From.FirstName)
				if err != nil {
					r.logger.Error("send notification failed", err)
				}
			}
		}
		wg.Done()
	}()

	return nil
}

func (r *Receiver) getClient(client *client.Client) (*client.Client, error) {
	cli, err := r.clientService.GetClient(client.ChatId)
	if err != nil {
		r.logger.Error("get clients failed", err)
		return nil, err
	}

	r.logger.Info("clients", cli)

	if len(cli) == 0 {
		err := r.clientService.Add(client)
		if err != nil {
			r.logger.Error("get clients failed", err)
			return nil, err
		}
	}

	cli, err = r.clientService.GetClient(client.ChatId)
	if err != nil {
		r.logger.Error("get clients failed", err)
		return nil, err
	}

	return &cli[0], err
}
