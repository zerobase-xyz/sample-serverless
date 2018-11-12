package slacsops

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/nlopes/slack"
)

func Inter(w http.ResponseWriter, r *http.Request) {
	dyna := DynaSession{Tablename: tablename, Hashkey: hashkey}
	dyna.Dynamodb = DynamoSession(region)
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
	e := EcsService{DynaSession: dyna, ServiceName: form[0]}
	result, err := e.GetData()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	payload := EcsUpdate{EcsService: e, Count: form[1]}

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
