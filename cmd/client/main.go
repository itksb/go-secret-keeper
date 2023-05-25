package main

import (
	"fmt"
	"github.com/itksb/go-secret-keeper/internal/client"
	"log"
)

var (
	gitBuildCommit = "N/A"
	gitBuildTag    = "N/A"
	buildDateTime  = "N/A"
)

func main() {

	fmt.Printf("Client. Build commit: %s, build date and time: %s, build tag: %s \n",
		gitBuildCommit,
		buildDateTime,
		gitBuildTag,
	)

	cfg := client.NewConfig()
	cfg.UseFlags()
	_, err := cfg.UseJsonConfigFile()
	if err != nil {
		log.Fatalf("error while using config file: %s", err.Error())
		return
	}

	application := client.NewClientApp(*cfg)
	application.Run()
}
