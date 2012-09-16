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
	log.Println("program ends")	
}