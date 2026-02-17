package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"

	"eventryx.api_service/config"
	"eventryx.api_service/internal/database"
)

func migrate(operation string) error {
	database.InitConnectionString(
		config.Config.DbHost, config.Config.DbUsername,
		config.Config.DbPassword, config.Config.DbPort, config.Config.DbName,
	)

	commandArgs := []string{"-dir", config.Config.DirMigrations, "postgres", database.ConnectionString}

	operationArgs := strings.Split(operation, " ")
	commandArgs = append(commandArgs, operationArgs...)

	cmd := exec.Command("goose", commandArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		panic(errors.New("error running migrations: " + err.Error()))
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate <operation>")
	}

	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = migrate(strings.Join(os.Args[1:], " "))
	if err != nil {
		log.Fatal(err)
	}
}
