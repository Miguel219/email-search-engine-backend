package server

import (
	"fmt"
	"net/http"
)

const host = "http://localhost:4080/api"
const user = "admin"
const password = "Complexpass#123"

func getUrlByIndex(index string, _type string) string {
	return fmt.Sprintf("%s/%s/%s", host, index, _type)
}
func getUrl(_type string) string {
	return fmt.Sprintf("%s/%s", host, _type)
}

func setHeaders(req *http.Request) *http.Request {
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/x-ndjson")
	return req
}
