package slacsops

import (
	"os"
)

var (
	slackToken  = os.Getenv("SLACK_TOKEN")
	botToken    = os.Getenv("SLACK_BOT_TOKEN")
	postChannel = os.Getenv("POST_CHANNEL")
	hashkey     = "ServiceName"
	region      = "ap-northeast-1"
	tablename   = "ecs_operation"
)
