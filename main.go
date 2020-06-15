package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/ini.v1"
	tb "gopkg.in/tucnak/telebot.v2"
)

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	config, err := ini.Load("config.ini")
	checkErr(err)

	owner, err := config.Section("bot").Key("owner").Int()
	checkErr(err)

	logchannel, err := config.Section("bot").Key("logchannel").Int64()
	checkErr(err)

	watchchannel, err := config.Section("bot").Key("channeltowatch").Int64()
	checkErr(err)

	b, err := tb.NewBot(tb.Settings{
		Token:  config.Section("bot").Key("token").String(),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	checkErr(err)

	b.Handle("/setmessage", func(m *tb.Message) {
		if m.Sender.ID == owner {
			config.SaveTo("config.ini.bak")
			config.Section("bot").Key("newusermessage").SetValue(strings.SplitN(m.Text, " ", 2)[1])
			config.SaveTo("config.ini")
			b.Send(m.Chat, "Saved")
		}
	})

	b.Handle("/testmessage", func(m *tb.Message) {
		b.Updates <- tb.Update{
			Message: &tb.Message{
				Chat: &tb.Chat{
					ID:         -1001050101913,
					Title:      "Chat Title",
					Type:       tb.ChatGroup,
					Username:   "testUsername",
					InviteLink: "jiASOD8saHUDOSa",
				},
				UserJoined: &tb.User{
					ID:        55258520,
					Username:  "asstra",
					FirstName: "Astra",
					LastName:  "Vexton",
				},
			},
		}
	})

	b.Handle(tb.OnUserJoined, func(m *tb.Message) {
		if m.Chat.ID == watchchannel {

			var msg = config.Section("bot").Key("newusermessage").String()
			msg = formatMessage(m, msg)
			b.Send(&tb.Chat{ID: logchannel}, msg, &tb.SendOptions{ParseMode: tb.ModeHTML})
		}
	})

	fmt.Println("Bot running")

	b.Start()

}

func escapeMarkdownChar(s string) string {
	// replacer := strings.NewReplacer("*", "\\*", "_", "\\_", "[", "\\[", "`", "\\`")
	replacer := strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")
	return replacer.Replace(s)
}

func formatMessage(data *tb.Message, msg string) (message string) {
	firstLast := data.UserJoined.FirstName
	if data.UserJoined.LastName != "" {
		firstLast = fmt.Sprintf("%s %s", firstLast, data.UserJoined.LastName)
	}

	message = strings.Replace(msg, "$username", fmt.Sprintf("@%s", data.UserJoined.Username), -1)
	message = strings.Replace(message, "$userid", data.UserJoined.Recipient(), -1)
	message = strings.Replace(message, "$firstname", escapeMarkdownChar(data.UserJoined.FirstName), -1)
	message = strings.Replace(message, "$lastname", escapeMarkdownChar(data.UserJoined.LastName), -1)
	message = strings.Replace(message, "$firstlastname", escapeMarkdownChar(firstLast), -1)
	message = strings.Replace(message, "$chatid", data.Chat.Recipient()[4:], -1)
	message = strings.Replace(message, "$chattitle", data.Chat.Title, -1)
	message = strings.Replace(message, "$chatinvitelink", fmt.Sprintf("tg://join?invite=%s", data.Chat.InviteLink), -1)

	return
}
