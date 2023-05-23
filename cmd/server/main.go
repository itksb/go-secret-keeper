package main

import (
	"fmt"
	"github.com/itksb/go-secret-keeper/internal/server"
	"github.com/itksb/go-secret-keeper/migrate"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	gitBuildCommit string = "N/A"
	gitBuildTag    string = "N/A"
	buildDateTime  string = "N/A"
)

func main() {
	fmt.Printf("Client. Build commit: %s, build date and time: %s, build tag: %s \n",
		gitBuildCommit,
		buildDateTime,
		gitBuildTag,
	)

	cfg := server.NewConfig()
	cfg.UseFlags()
	_, err := cfg.UseJsonConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	application := server.NewServerApp(*cfg, migrate.AppMigratorFunc(migrate.Migrate))
	if err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
		<-sigint
		application.GRPCServer.GracefulStop()
		close(doneCh)
	}()

	log.Println(application.Run())
	<-doneCh
	err = application.Stop()
	if err != nil {
		log.Printf("error while closing application %s", err)
	}

	log.Println("bye bye")

}
