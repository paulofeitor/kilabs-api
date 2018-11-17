package main

import (
	"github.com/paulofeitor/kilabs-api/app"
	"github.com/paulofeitor/kilabs-api/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
