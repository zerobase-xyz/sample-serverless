package slacsops

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var (
	slackToken  = os.Getenv("SLACK_TOKEN")
	botToken    = os.Getenv("SLACK_BOT_TOKEN")
	postChannel = os.Getenv("POST_CHANNEL")
	hashkey     = "ServiceName"
	region      = "ap-northeast-1"
	tablename   = "ecs_operation"
)

func Slash(w http.ResponseWriter, r *http.Request) {
	dyna := DynaSession{Tablename: tablename, Hashkey: hashkey}
	dyna.Dynamodb = DynamoSession(region)
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
		form := strings.Split(s.Text, " ")
		if len(form) != 6 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		e := EcsService{
			DynaSession: dyna,
			ServiceName: form[0],
			Project:     form[1],
			Account:     form[2],
			Cluster:     form[3],
			Region:      form[4],
			EcsName:     form[5],
		}

		if err = e.PutData(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
		}
		params := &slack.Msg{Text: "Success"}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/ecs-get-data":
		e := EcsService{DynaSession: dyna, ServiceName: s.Text}
		result, err := e.GetData()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := &slack.Msg{Text: result.Describe()}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/ecs-list-data":
		result, err := dyna.ScanData()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var results []string
		for _, v := range result {
			results = append(results, v.Describe())
		}
		params := &slack.Msg{Text: strings.Join(results, "\n")}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/ecs-delete-data":
		e := EcsService{DynaSession: dyna, ServiceName: s.Text}
		err := e.DeleteData()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		params := &slack.Msg{Text: "Success!"}
		b, err = json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "/ecs-update":
		e := EcsService{DynaSession: dyna, ServiceName: s.Text}
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

		_, err = e.GetData()
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
