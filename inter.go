package slacsops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

func Inter(w http.ResponseWriter, r *http.Request) {
	slackToken := os.Getenv("SLACK_TOKEN")
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.Token != slackToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if res.Actions[0].Value != "approve" {
		params := &slack.Msg{Text: "Request Reject!"}
		b, _ := json.Marshal(params)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}
	form := strings.Split(res.OriginalMessage.Text, " ")
	count, err := strconv.Itoi(form[1])
	ecs, err := GetData(form[0])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	payload := Ecsupdate{ecs, count}

	err = InvokeLambda(payload)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := &slack.Msg{Text: "Update Success"}
	b, _ := json.Marshal(params)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return
}
