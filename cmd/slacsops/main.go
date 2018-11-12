package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zerobase-xyz/slacsops"

	"github.com/apex/gateway"
)

func main() {
	ApexGatewayDisabled := os.Getenv("APEX_GATEWAY_DISABLED")
	http.HandleFunc("/slash", slacsops.Slash)
	http.HandleFunc("/interr", slacsops.Inter)

	if ApexGatewayDisabled == "true" {
		log.Fatal(http.ListenAndServe(":3000", nil))
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", nil))
	}
}
