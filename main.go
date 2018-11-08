package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/nlopes/slack"
)

func main() {
	ApexGatewayDisabled := os.Getenv("APEX_GATEWAY_DISABLED")
	slackToken := os.Getenv("SLACK_TOKEN")
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
	table := db.Table("ecs_operation")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
			ecs := &Ecsservice{
				ServiceName: form[0],
				Project:     form[1],
				Account:     account,
				Cluster:     form[3],
				Region:      form[4],
				EcsName:     form[5],
			}

			err = PutData(ecs, table)
			if err != nil {
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

	})
	if ApexGatewayDisabled == "true" {
		log.Fatal(http.ListenAndServe(":3000", nil))
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", nil))
	}
}
