package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// Use struct tags much like the standard JSON library,
// you can embed anonymous structs too!
type ecsservice struct {
	ServiceName string
	Project     string `dynamo:"Project"`
	Account     int    `dynamo:"Account"`
	Cluster     string `dynamo:"Cluster"`
	Region      string `dynamo:"Region"`
	EcsName     string `dynamo:"EcsName"`
	Time        time.Time
}

type ecsservices []ecsservice

const hashkey = "Servicename"

func PutData(w *ecsservice, table *dynamo.Table) error {
	err := table.Put(w).Run()
	if err != nil {
		return err
	}
	err = table.Put(w).Run()
	if err != nil {
		return err
	}
	return nil
}

func GetData(w *ecsservice, table *dynamo.Table) (ecsservice, error) {
	var result ecsservice
	err := table.Get(hashkey, w.ServiceName).
		One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func ScanData(w *ecsservice, table *dynamo.Table) (ecsservices, error) {
	var results []ecsservice
	err := table.Scan().All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func DeleteData(w *ecsservice, table *dynamo.Table) (ecsservice, error) {
	var result ecsservice
	err := table.Delete(hashkey, w.ServiceName).OldValue(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func main() {
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
	table := db.Table("ecs_operation")
	w := ecsservice{ServiceName: "test", Time: time.Now(), Project: "hello", Account: 21904811, Cluster: "test", Region: "ap-northeasti-1", EcsName: "test"}
}
