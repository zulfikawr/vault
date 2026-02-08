package main

import (
	"github.com/zulfikawr/vault/internal/server"
)

func main() {
	app := server.NewApp()
	app.Run()
}
