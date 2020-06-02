package main

import (
	"factory/in_progress/app"
	"factory/in_progress/config"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)

	app.Run(":3000")
}
