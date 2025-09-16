package main

import (
	"fmt"
	"log"
	"lunar-backend-engineer-challenge/cmd/di"
	"os"
	"strings"
)

func main() {
	rocketsDI := di.Init()

	n, err := askForOptions(map[string]func() (int, error){
		"up":   func() (int, error) { return rocketsDI.Services.DatabaseMigrator.Up() },
		"down": func() (int, error) { return rocketsDI.Services.DatabaseMigrator.Down() },
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
}

func askForOptions(validOptions map[string]func() (int, error)) (int, error) {

	for {
		response := os.Args[1]

		response = strings.ToLower(strings.TrimSpace(response))
		result, ok := validOptions[response]
		if !ok {
			log.Fatalf("Invalid option provided %s", response)
		}

		return result()
	}
}
