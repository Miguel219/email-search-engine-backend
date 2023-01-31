package main

import (
	i "email-search-engine-backend/internal/importData"
	server "email-search-engine-backend/internal/server"
)

func main() {
	if false {
		i.ImportData()
	}
	server.CreateServer()
}
