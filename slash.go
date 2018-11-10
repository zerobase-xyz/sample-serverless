package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/guregu/dynamo"
	"github.com/nlopes/slack"
)

func Slash(w http.ResponseWriter, r *http.Request) {
	slackToken := os.Getenv("SLACK_TOKEN")
	sess := session.Must(session.NewSession())
	db := dynamo.New(sess)
	table := db.Table("ecs_operation")
	hashkey := "ServiceName"
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if s.Token != slackToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	switch s.Command {
	case "/ecs-put-data":
		form := strings.Split(s.Text, " ")
		account, err := strconv.Atoi(form[2])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if len(form) != 6 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ecs := Ecsservice{
			ServiceName: form[0],
			Project:     form[1],
			Account:     account,
			Cluster:     form[3],
			Region:      form[4],
			EcsName:     form[5],
		}

		if err = PutData(ecs, table); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := &slack.Msg{Text: "Success"}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case "/ecs-get-data":
		var ecs Ecsservice
		ecs, err = GetData(s.Text, hashkey, table)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := &slack.Msg{Text: ecs.Describe()}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case "/ecs-list-data":
		var services Ecsservices
		services, err = ScanData(table)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var servicesText []string
		for _, ecs := range services {
			servicesText = append(servicesText, ecs.Describe())
		}

		text := strings.Join(servicesText, " ")

		params := &slack.Msg{Text: text}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case "/ecs-delete-data":
		var ecs Ecsservice
		ecs, err = DeleteData(s.Text, hashkey, table)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		params := &slack.Msg{Text: ecs.Describe()}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case "/ecs-update":
		var ecs Ecsservice
		var attachment = slack.Attachment{
			Text:       "Operator:{} ServiceName:{} Desired:{}",
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

		_, err = GetData(s.Text, hashkey, table)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		api := slack.New("YOUR_TOKEN_HERE")
		param := slack.PostMessageParameters{Attachments: attachments}

		_, _, err := api.PostMessage("CHANNEL_ID", "<!here>\n Do you want to permit this update?", param)

		params := &slack.Msg{Text: "Success!"}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		lamb := lambda.New(sess)
		payload, err := json.Marshal(ecs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		input := &lambda.InvokeInput{
			FunctionName:   aws.String("ecs-update-service"),
			InvocationType: aws.String("Event"),
			Payload:        payload,
		}
		_, err = lamb.Invoke(input)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case "/echo":
		params := &slack.Msg{Text: s.Text}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
