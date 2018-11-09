package main

import (
	"strconv"
	"strings"

	"github.com/guregu/dynamo"
)

type Ecsservice struct {
	ServiceName string `json:"ServiceName"`
	Project     string `json:"Project"`
	Account     int    `json:"Account"`
	Cluster     string `json:"Cluster"`
	Region      string `json:"Region"`
	EcsName     string `json:"EcsName"`
}

type Ecsservices []Ecsservice

func (ecs *Ecsservice) Describe() string {
	textS := []string{ecs.ServiceName, ecs.Project, strconv.Itoa(ecs.Account), ecs.Cluster, ecs.Region, ecs.EcsName}
	text := strings.Join(textS, " ")
	return text
}

func PutData(w Ecsservice, table dynamo.Table) error {
	err := table.Put(w).Run()
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
