package main

import "diveGames"

func main() {
	r := InitializeServer(diveGames.KrakenURL)
	r.Run(":8080")
}
