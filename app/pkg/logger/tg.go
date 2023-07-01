package logger

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (a *apiLogger) InitTG() error {
	var err error
	a.tgBot, err = tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".

		Token:  a.cfg.Logger.TGToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	return err
}

func (a *apiLogger) SendTGLogMessage(message string) {
	_, err := a.tgBot.Send(tb.ChatID(a.cfg.Logger.ChatID), fmt.Sprintf("%s\n%s", a.makeAlertMessage(), message), &tb.SendOptions{DisableWebPagePreview: true})
	if err != nil {
		log.Err(err)
	}
}

func (a *apiLogger) makeAlertMessage() (message string) {
	message = "Alert: "
	for i, r := range a.cfg.Logger.AlertUsers {
		message += r
		if i != len(a.cfg.Logger.AlertUsers)-1 {
			message += " "
		}
	}
	return
}
