package slacsops

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func InvokeLambda(param Ecsupdate) error {
	sess := session.Must(session.NewSession())
	s := ecs.New(sess)
	update := ecs.UpdateServiceInput{}
	update.Cluster = &param.ecsservice.Cluster

	lamb := lambda.New(sess)
	payload, err := json.Marshal(param)
	input := &lambda.InvokeInput{
		FunctionName:   aws.String("ecs-update-service"),
		InvocationType: aws.String("Event"),
		Payload:        payload,
	}
	_, err = lamb.Invoke(input)
	if err != nil {
		return err
	}
	return nil

}
