package main

import (
	"time"

	"github.com/guregu/dynamo"
)

type Ecsservice struct {
	ServiceName string
	Project     string `dynamo:"Project"`
	Account     int    `dynamo:"Account"`
	Cluster     string `dynamo:"Cluster"`
	Region      string `dynamo:"Region"`
	EcsName     string `dynamo:"EcsName"`
	Time        time.Time
}

type Ecsservices []Ecsservice

func PutData(w *Ecsservice, table dynamo.Table) error {
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

func GetData(w string, hashkey string, table dynamo.Table) (Ecsservice, error) {
	var result Ecsservice
	err := table.Get(hashkey, w).
		One(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func ScanData(table dynamo.Table) (Ecsservices, error) {
	var results []Ecsservice
	err := table.Scan().All(&results)
	if err != nil {
		return results, err
	}
	return results, nil
}

func DeleteData(w string, hashkey string, table dynamo.Table) (Ecsservice, error) {
	var result Ecsservice
	err := table.Delete(hashkey, w).OldValue(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}
