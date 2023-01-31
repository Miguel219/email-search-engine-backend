package main

import (
	i "email-search-engine-backend/internal/importData"
	server "email-search-engine-backend/internal/server"
	"fmt"
	"strings"
)

func ask4confirmation(question string) bool {
	var s string

	fmt.Printf(question)
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func main() {
	if ask4confirmation("Do you want to import emails? (y/n): ") {
		i.ImportData()
	}
	server.CreateServer()
}
