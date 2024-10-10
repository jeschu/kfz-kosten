package main

import (
	"fmt"
	"kfz-kosten/model"
	"log"
)

func main() {
	kfzs, err := model.LoadKfzs()
	if err != nil {
		log.Printf("Error loading Kfzs: %v", err)
	}

	fmt.Println("")

	for _, kfz := range kfzs {
		kfz.PrintStats()
	}

	if err := model.SaveKfzs(kfzs); err != nil {
		log.Printf("Error saving Kfzs: %v", err)
	}
}
