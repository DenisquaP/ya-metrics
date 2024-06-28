package main

import (
	"log"

	"github.com/DenisquaP/ya-metrics/internal/agent"
)

func main() {
	log.Println("agent run")
	agent.Run()
}
