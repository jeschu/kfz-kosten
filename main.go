package main

import (
	"kfz-kosten/model"
	"log"
)

func main() {
	kfzs, err := model.LoadKfzs()
	if err != nil {
		log.Printf("Error loading Kfzs: %v", err)
	}
	PrintStats(kfzs["Mini"])
	if err := model.SaveKfzs(kfzs); err != nil {
		log.Printf("Error saving Kfzs: %v", err)
	}
}
