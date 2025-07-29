package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"vikunja-mcp",
		"1.0.0",
	)

	registerTools(s)

	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}
