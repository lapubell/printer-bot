package main

import (
	"errors"
	"os"

	"github.com/lapubell/printer-bot/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err := processEnv()
	if err != nil {
		panic(err)
	}

	s, err := server.New()
	if err != nil {
		panic(err)
	}

	err = s.Serve()
	if err != nil {
		panic(err)
	}
}

func processEnv() error {
	requiredENV := []string{
		"APP_PORT",
	}

	for _, env := range requiredENV {
		check := os.Getenv(env)

		if check == "" {
			return errors.New("missing env: " + env)
		}
	}

	return nil
}
