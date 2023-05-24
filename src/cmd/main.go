package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	repository "github.com/dscabral/ports/src/repository/sql"
	"github.com/dscabral/ports/src/service"
	_ "modernc.org/sqlite"
)

const (
	DBFile = "ports.db"
)

func main() {
	portsRepository := repository.NewPortRepository(DBFile)
	defer func() {
		fmt.Println("shutting down database connection")
		if err := portsRepository.Close(); err != nil {
			log.Fatalf("failed to shutdown database connection: %v", err)
		}
	}()

	err := portsRepository.Init()
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	portService := service.NewPortService(portsRepository)
	path := "ports.json"
	err = portService.SaveOrUpdatePortFromFile(path)
	if err != nil {
		log.Fatalf("failed to import and save the ports: %v", err)
	}

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-interruptChannel
	fmt.Printf("Shutting down %s\n", sig)

	// Clean up the database file after shutting down
	if err := os.Remove(DBFile); err != nil {
		log.Println(err)
	}
}
