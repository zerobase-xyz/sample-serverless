package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/nlopes/slack"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "")
}

func echo(w http.ResponseWriter, r *http.Request) {
	var verificationToken string
	slackToken := os.Getenv("SLACK_TOKEN")
	flag.StringVar(&verificationToken, "token", slackToken, "Your Slash Verification Token")
	flag.Parse()
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !s.ValidateToken(verificationToken) {
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
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/echo", echo)
	ApexGatewayDisabled := os.Getenv("APEX_GATEWAY_DISABLED")
	if ApexGatewayDisabled == "true" {
		log.Fatal(http.ListenAndServe(":3000", nil))
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", nil))
	}
}
