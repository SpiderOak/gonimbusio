package main

import (
	"log"
)

func main() {
	log.Println("program starts")
	flags, err := loadFlags(); if err != nil {
		log.Fatalf("Unable to load flags: %s", err)
	}
	log.Printf("flags = %v", flags)
	config, err := LoadConfig(flags.configPath); if err != nil {
		log.Fatalf("Unable to load config from %s %s", flags.configPath, err)
	}
	log.Printf("config = %v", config)
	log.Println("program ends")	
}