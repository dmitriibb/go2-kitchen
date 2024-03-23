package api

import (
	"fmt"
	"net/http"
)

func GetActiveWorkers(response *http.ResponseWriter, request *http.Request) {
	fmt.Fprint(*response, "No active workers")
}
