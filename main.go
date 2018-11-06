package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/nlopes/slack"
)

func echo(w http.ResponseWriter, r *http.Request) {
	slackToken := os.Getenv("SLACK_TOKEN")
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
	case "/echo":
		params := &slack.Msg{Text: s.Text}
		b, err := json.Marshal(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	case "/test":
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

func main() {
	http.HandleFunc("/", echo)
	ApexGatewayDisabled := os.Getenv("APEX_GATEWAY_DISABLED")
	if ApexGatewayDisabled == "true" {
		log.Fatal(http.ListenAndServe(":3000", nil))
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", nil))
	}
}
