package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/priyankasharma10/ReNew/server"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword() (string, error) {
	// Generate a salted hash of the password
	password := "Admin"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Convert the hashed password to a string and return it
	return string(hashedPassword), nil
}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	srv := server.SrvInit()
	go srv.Start()

	<-done
	logrus.Info("Graceful shutdown")
	srv.Stop()

}
