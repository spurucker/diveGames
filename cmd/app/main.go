package main

import (
	"diveGames"
	"diveGames/internal/app"
)

func main() {
	r := app.NewApp(diveGames.KrakenURL)
	r.Run(":8080")
}
