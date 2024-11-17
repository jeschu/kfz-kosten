package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"kfz-kosten/input"
	"kfz-kosten/model"
	"log"
	"strconv"
	"time"
)

type actionType int

func (a actionType) String() string {
	switch a {
	case actionTanken:
		return "Tanken erfassen"
	case actionKosten:
		return "Kosten erfassen"
	case actionExcel:
		return "Excel Export"
	case actionSummary:
		return "Übersicht anzeigen"
	default:
		return fmt.Sprintf("?%d", int(a))
	}
}

const (
	actionTanken actionType = iota
	actionKosten
	actionExcel
	actionSummary
)

func main() {
	kfzs, err := model.LoadKfzs()
	if err != nil {
		log.Printf("Error loading Kfzs: %v", err)
	}

	var kfz model.Kfz
	for _, kfz = range kfzs {
		break
	}

	loop := true
	for loop {
		fmt.Printf("[t] %s\n", actionTanken)
		fmt.Printf("[k] %s\n", actionKosten)
		fmt.Printf("[e] %s\n", actionExcel)
		fmt.Printf("[␍] %s\n", actionSummary)
		action := input.ReadSelectionMapped(
			"Was möchtest du tun? ",
			map[string]actionType{"t": actionTanken, "k": actionKosten, "e": actionExcel},
			actionSummary,
			"t", "k", "e", input.CR,
		)
		fmt.Print("\n\n")
		switch action {
		case actionTanken:
			tanken(kfz)
		case actionKosten:
			kosten(kfz)
		case actionExcel:
			excel(kfz)
		case actionSummary:
			loop = false
		}
	}

	fmt.Println()
	kfz.PrintStats()

	if err := model.SaveKfzs(kfzs); err != nil {
		log.Printf("Error saving Kfzs: %v", err)
	}
}

func tanken(kfz model.Kfz) {
	fmt.Println("Tanken erfassen:")

	date := input.ReadDateInPast("  Datum")

	art := input.ReadSelectionMapped[model.TankArt](
		"  Art ([V]oll-, [T]eil-, [E]rstbetankung): ",
		map[string]model.TankArt{"t": model.Teil, "e": model.Erst},
		model.Voll,
		"v", "t", "e", input.CR,
	)

	km := input.Read[int]("  Tachostand: ", func(s string) (int, error) { return strconv.Atoi(s) })
	liter := input.Read[float64]("  Liter: ", func(s string) (float64, error) { return strconv.ParseFloat(s, 64) })
	preis := input.Read[float64]("  Kosten: ", func(s string) (float64, error) { return strconv.ParseFloat(s, 64) })
	sorte := input.ReadString("  Sorte: ")

	kfz.Tanken = append(kfz.Tanken, model.Tanken{Datum: date, Art: art, Km: km, Liter: liter, Kosten: preis, Sorte: sorte})
}

func kosten(kfz model.Kfz) {
	fmt.Println("Kosten erfassen:")
}

func excel(kfz model.Kfz) {
	var err error
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	sheetName := strconv.Itoa(time.Now().Year())
	sheetIndex := f.GetActiveSheetIndex()
	f.SetSheetName(f.GetSheetName(sheetIndex), sheetName)
	f.SetCellStr(sheetName, "A1", "KFZ Kosten "+sheetName+" - "+kfz.Name+" ["+kfz.Kennzeichen+"]")
	f.SetActiveSheet(sheetIndex)
	if err = f.SaveAs("KFZ_" + sheetName + ".xlsx"); err != nil {
		log.Println(err)
	}
}
