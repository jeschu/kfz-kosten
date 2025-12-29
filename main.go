package main

import (
	"fmt"
	"kfz-kosten/input"
	"kfz-kosten/lang"
	"kfz-kosten/model"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
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
			saveKfzs(kfzs)
			break
		case actionKosten:
			kosten(kfz)
			saveKfzs(kfzs)
			break
		case actionExcel:
			excel(kfz)
			break
		case actionSummary:
			loop = false
			break
		}
	}

	fmt.Println()
	kfz.PrintStats()

	saveKfzs(kfzs)
}

func saveKfzs(kfzs model.Kfzs) {
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

	km := input.ReadInt("  Tachostand: ")
	liter := input.ReadFloat64("  Liter: ")
	preis := input.ReadFloat64("  Kosten: ")
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
	f.SetCellStr(sheetName, "A2", "Kosten")
	f.SetCellStr(sheetName, "A3", "Datum")
	f.SetCellStr(sheetName, "B3", "Tachostand")
	f.SetCellStr(sheetName, "C3", "Preis")
	f.SetCellStr(sheetName, "D3", "Kategorie")
	f.SetCellStr(sheetName, "E3", "Abschreibung")
	f.SetCellStr(sheetName, "F3", "🏦")
	f.SetCellStr(sheetName, "G3", "Bemerkung")
	sort.SliceStable(kfz.Kosten, func(i, j int) bool { return false })
	row := 4
	for _, kosten := range kfz.Kosten {
		abschreibung := ""
		if kosten.AbschreibungZeit > 0 {
			abschreibung = fmt.Sprintf("%s (bis %s)",
				lang.FormatDuration(kosten.AbschreibungZeit),
				kosten.Datum.Add(kosten.AbschreibungZeit-24*time.Hour).Format("02.01.2006"),
			)
		} else if kosten.AbschreibungKm > 0 {
			abschreibung = fmt.Sprintf("%d km", kosten.AbschreibungKm)
		}
		fa := ""
		switch kosten.AbschreibungFa {
		case model.Ja:
			fa = "✓"
		case model.Nein:
			fa = "𐄂"
		case model.Abschreibung:
			fa = "⏲"
		}
		f.SetCellValue(sheetName, cell("A", row), kosten.Datum.Format("02.01.2006"))
		f.SetCellValue(sheetName, cell("B", row), kosten.Km)
		f.SetCellValue(sheetName, cell("C", row), kosten.Kosten)
		f.SetCellValue(sheetName, cell("D", row), kosten.Kategorie)
		f.SetCellValue(sheetName, cell("E", row), abschreibung)
		f.SetCellValue(sheetName, cell("F", row), fa)
		f.SetCellValue(sheetName, cell("G", row), kosten.Notiz)
		row++
	}

	row++
	f.SetCellStr(sheetName, cell("A", row), "Tanken")
	row++
	f.SetCellStr(sheetName, cell("A", row), "Datum")
	f.SetCellStr(sheetName, cell("B", row), "Tachostand")
	f.SetCellStr(sheetName, cell("C", row), "Preis")
	f.SetCellStr(sheetName, cell("D", row), "Liter")
	row++

	for _, tanken := range kfz.Tanken {
		f.SetCellValue(sheetName, cell("A", row), tanken.Datum.Format("02.01.2006"))
		f.SetCellValue(sheetName, cell("B", row), tanken.Km)
		f.SetCellValue(sheetName, cell("C", row), tanken.Kosten)
		f.SetCellValue(sheetName, cell("D", row), tanken.Liter)
		row++
	}
	f.SetActiveSheet(sheetIndex)
	if err = f.SaveAs("KFZ_" + sheetName + ".xlsx"); err != nil {
		log.Println(err)
	}
}

func cell(col string, row int) string { return fmt.Sprintf("%s%d", col, row) }
