package slacsops

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func DynamoSession(r string) *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess, aws.NewConfig().WithRegion(r))
	return svc
}

type DynaSession struct {
	Dynamodb  *dynamodb.DynamoDB
	Tablename string
	Hashkey   string
}

type EcsService struct {
	DynaSession DynaSession
	ServiceName string `json:"ServiceName"`
	Project     string `json:"Project"`
	Account     string `json:"Account"`
	Cluster     string `json:"Cluster"`
	Region      string `json:"Region"`
	EcsName     string `json:"EcsName"`
}

type EcsServices []EcsService

type EcsUpdate struct {
	EcsService EcsService
	Count      string
}

func (e *EcsService) Describe() string {
	textS := []string{e.ServiceName, e.Project, e.Account, e.Cluster, e.Region, e.EcsName}
	return strings.Join(textS, " ")
}

func (e *EcsService) PutData() error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ServiceName": {
				S: aws.String(e.ServiceName),
			},
			"Project": {
				S: aws.String(e.Project),
			},
			"Account": {
				S: aws.String(e.Account),
			},
			"Cluster": {
				S: aws.String(e.Cluster),
			},
			"Region": {
				S: aws.String(e.Region),
			},
			"EcsName": {
				S: aws.String(e.EcsName),
			},
		},
		TableName: aws.String(e.DynaSession.Tablename),
	}
	result, err := e.DynaSession.Dynamodb.PutItem(input)
	if err != nil {
		return nil
	}
	return nil
}

func (e *EcsService) GetData() (r *EcsService, err error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ServiceName": {
				S: aws.String(e.ServiceName),
			},
		},
		TableName: aws.String(e.DynaSession.Tablename),
	}

	result, err := e.DynaSession.Dynamodb.GetItem(input)
	err = dynamodbattribute.UnmarshalMap(result.Item, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (e *DynaSession) ScanData() (r EcsServices, err error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(e.Tablename),
	}

	result, err := e.Dynamodb.Scan(input)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (e *EcsService) DeleteData() error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ServiceName": {
				S: aws.String(e.ServiceName),
			},
		},
		TableName: aws.String(e.DynaSession.Tablename),
	}

	_, err := e.DynaSession.Dynamodb.DeleteItem(input)
	if err != nil {
		return err
	}
	return nil
}
