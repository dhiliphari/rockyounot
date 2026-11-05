package main

import (
	"log"

	"github.com/dhiliphari/rockyounot/data"
	"github.com/dhiliphari/rockyounot/internal/reader"
)

func main() {
	file, err := data.FS.Open("rockyou.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	lines := reader.Lines(file)

	count := 0
	for line := range lines {
		if line.Err != nil {
			log.Fatal(line.Err.Error())
		}
		count++
	}

	log.Printf("Total passwords: %d", count)
}
