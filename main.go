package main

import (
	"log"

	"github.com/joho/godotenv"
)

// loads values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	bootAPI()
	
	listen()
}

//boots REST API where zeta explorer connects to get information such as node health, registered nodes, connections
func bootAPI() {
	nodeAPI := NodeAPI{}
	go nodeAPI.Init()
}

//discovery node listening on peers.
//listens for registration, file fetch request, or file metainfo storage
func listen() {
	var dn DiscoveryNode
	dn.listen()
}
