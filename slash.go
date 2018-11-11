package slacsops

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

var slackToken = os.Getenv("SLACK_TOKEN")
var botToken = os.Getenv("SLACK_BOT_TOKEN")
var postChannel = os.Getenv("POST_CHANNEL")

func EcsPutData(text string) error {
	form := strings.Split(text, " ")
	account, err := strconv.Atoi(form[2])
	if err != nil {
		return err
	}
	if len(form) != 6 {
		panic("Invid request")
	}
	ecs := Ecsservice{
		ServiceName: form[0],
		Project:     form[1],
		Account:     account,
		Cluster:     form[3],
		Region:      form[4],
		EcsName:     form[5],
	}

	if err = PutData(ecs); err != nil {
		return err
	}
	return nil

}

func Slash(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if s.Token != slackToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var b []byte

	switch s.Command {
	case "/ecs-put-data":
		err := EcsPutData(s.Text)
		params := &slack.Msg{Text: "Success"}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/ecs-get-data":
		ecs, err := GetData(s.Text)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := &slack.Msg{Text: ecs.Describe()}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/ecs-list-data":
		services, err := ScanData()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var serviceslist []string
		for _, ecs := range services {
			serviceslist = append(serviceslist, ecs.Describe())
		}
		params := &slack.Msg{Text: strings.Join(serviceslist, "\n")}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/ecs-delete-data":
		ecs, err := DeleteData(s.Text)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		params := &slack.Msg{Text: ecs.Describe()}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/ecs-update":
		form := strings.Split(s.Text, " ")
		if len(form) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		text := "Operator:" + s.UserName + " ServiceName:" + form[0] + " Desired:" + form[1]
		var attachment = slack.Attachment{
			Text:       text,
			CallbackID: "ecs-update",
			Color:      "#3AA3E3",
			Actions: []slack.AttachmentAction{
				{
					Name:  "approve",
					Text:  "Approve",
					Style: "primary",
					Type:  "button",
					Value: "approve",
				},
				{
					Name:  "reject",
					Text:  "Reject",
					Style: "danger",
					Type:  "button",
					Value: "reject",
				},
			},
		}
		attachments := []slack.Attachment{attachment}

		_, err = GetData(s.Text)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		api := slack.New(botToken)
		param := slack.PostMessageParameters{Attachments: attachments}

		_, _, err := api.PostMessage(postChannel, "<!here>\n Do you want to permit this update?", param)

		params := &slack.Msg{Text: "Success!"}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	case "/echo":
		params := &slack.Msg{Text: s.Text}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

}
