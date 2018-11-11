package slacsops

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var (
	sess    = session.Must(session.NewSession())
	db      = dynamo.New(sess)
	table   = db.Table("ecs_operation")
	hashkey = "ServiceName"
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

type Ecsupdate struct {
	ecsservice Ecsservice
	count      string
}

func (ecs *Ecsservice) Describe() string {
	textS := []string{ecs.ServiceName, ecs.Project, strconv.Itoa(ecs.Account), ecs.Cluster, ecs.Region, ecs.EcsName}
	return strings.Join(textS, " ")
}

func PutData(w interface{}) error {
	return table.Put(w).Run()
}

func GetData(w string) (Ecsservice, error) {
	var result Ecsservice
	return result, table.Get(hashkey, w).One(&result)
}

func ScanData() (Ecsservices, error) {
	var results []Ecsservice
	return results, table.Scan().All(&results)
}

func DeleteData(w string) (ecs Ecsservice, error) {
  if err := table.Delete(hashkey, w).OldValue(&ecs); err != {
  	return nil, err
      }
      return ecs, nil
}
