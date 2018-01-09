package main

import (
	"log"

	"github.com/KitlerUA/CSVParser/parser"
)

const defaultCSVFileName = "List of Actions.xlsx"

func main() {
	if err := parser.Parse(defaultCSVFileName, ""); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("Successfully parsed and saved")
}
