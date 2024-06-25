package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load env vars: %v", err)
	}

	dbURL := getEnvVar("dbURL")
	addr := getEnvVar("ADDR")
	s := NewAPIServer(addr, dbURL)

	go FetchFeedsWorker(s, 3)
	// err = s.RunSetup(10)
	// if err != nil {
	// 	s.Logger.Warn("Failed to populate DB", "error", err.Error())
	// }
	log.Fatal(s.Run())
}

func getEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%v env var not found!", key)
	}
	return value
}
