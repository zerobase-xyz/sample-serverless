package slacsops

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	hashkey   = "ServiceName"
	region    = "ap-northeast-1"
	tablename = "ecs_operation"
)

func DynamoSess() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess, aws.NewConfig().WithRegion("region"))
	return svc
}

type Ecsservice struct {
	ServiceName string `json:"ServiceName"`
	Project     string `json:"Project"`
	Account     string `json:"Account"`
	Cluster     string `json:"Cluster"`
	Region      string `json:"Region"`
	EcsName     string `json:"EcsName"`
}

type Ecsservices []Ecsservice

type Ecsupdate struct {
	ecsservice Ecsservice
	count      string
}

func (e *Ecsservice) Describe() string {
	textS := []string{e.ServiceName, e.Project, e.Account, e.Cluster, e.Region, e.EcsName}
	return strings.Join(textS, " ")
}

func (e *Ecsservice) PutData() error {
	svc := DynamoSess()
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
		TableName: aws.String(tablename),
	}
	result, err := svc.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				fmt.Println(dynamodb.ErrCodeConditionalCheckFailedException, aerr.Error())
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}
	return nil
}

func (e *Ecsservice) GetData(k string) (r *Ecsservice, err error) {
	svc := DynamoSess()
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ServiceName": {
				S: aws.String(e.ServiceName),
			},
		},
		TableName: aws.String(tablename),
	}

	r, err := svc.GetItem(input)
	return r, nil
}

func ScanData() (Ecsservices, error) {
	var results []Ecsservice
	return results, table.Scan().All(&results)
}

func DeleteData(w string) (e Ecsservice, err error) {
	err = table.Delete(hashkey, w).OldValue(&e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
