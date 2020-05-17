package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sasbury/mini"
	tb "gopkg.in/tucnak/telebot.v2"
)

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	config, err := mini.LoadConfiguration("config.ini")
	checkErr(err)

	b, err := tb.NewBot(tb.Settings{
		Token:  config.StringFromSection("bot", "token", ""),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	checkErr(err)

	b.Handle(tb.OnUserJoined, func(m *tb.Message) {
		log.Printf("OnUserJoined - %+v\n", m)

		if m.Chat.ID == -1001090101913 {
			fmt.Printf("%+v\n", m)

			var msg = "\\[Gearheads\\]\nNew member\n\n👤"
			var name = fmt.Sprintf("[%s", escapeMarkdownChar(m.UserJoined.FirstName))
			if m.UserJoined.LastName != "" {
				name = fmt.Sprintf("%s %s](tg://user?id=%s)", name, escapeMarkdownChar(m.UserJoined.LastName), m.UserJoined.Recipient())
			} else {
				name = fmt.Sprintf("%s](tg://user?id=%s)", name, m.UserJoined.Recipient())
			}

			if m.UserJoined.Username != "" {
				msg = fmt.Sprintf("%s %s \\(@%s\\)", msg, name, m.UserJoined.Username)
			}

			msg = fmt.Sprintf("%s\n\\#new\\_member \\#c%s \\#u%s", msg, m.Chat.Recipient()[4:], m.UserJoined.Recipient())

			b.Send(&tb.Chat{ID: config.IntegerFromSection("bot", "logchannel", 0)}, msg, &tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
		}
	})

	b.Start()

}

func escapeMarkdownChar(s string) string {
	replacer := strings.NewReplacer("*", "\\*", "_", "\\_", "[", "\\[", "`", "\\`")
	return replacer.Replace(s)
}
