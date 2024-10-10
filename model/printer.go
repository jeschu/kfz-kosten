package model

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"kfz-kosten/lang"
	"sort"
	"strings"
	"time"
)

var fmt = message.NewPrinter(language.German)

//goland:noinspection GoUnhandledErrorResult
func (kfz *Kfz) PrintStats() {
	kfz.printHeader()
	fmt.Println()
	kfz.printKosten()
	kfz.printTanken()
	fmt.Println()
	kfz.printSummary()
}

//goland:noinspection GoUnhandledErrorResult
func (kfz *Kfz) printHeader() {
	label := fmt.Sprintf("%s [ %s ]", kfz.Name, kfz.Kennzeichen)
	fmt.Println(label)
	fmt.Println(strings.Repeat("=", len(label)))
}

//goland:noinspection GoUnhandledErrorResult
func (kfz *Kfz) printTanken() {
	fmt.Println("  Tanken - Aufstellung:")
	fmt.Println("    | Datum      | Tachostand | Preis      | Liter  | €/l       | Distanz | l/100km  | l/100km |")
	sort.SliceStable(kfz.Tanken, func(i, j int) bool {
		return false
	})
	startKm := -1
	lastKm := -1
	sumLiter := 0.
	sumDistanz := 0
	for _, tanken := range kfz.Tanken {
		if lastKm == -1 {
			lastKm = tanken.Km
		}
		if startKm == -1 && tanken.Art == Voll {
			startKm = tanken.Km
		}
		distanz := tanken.Km - lastKm
		sumDistanz += distanz
		var lastKmStr string
		if distanz > 0 {
			lastKmStr = fmt.Sprintf("%5dkm", distanz)
		} else {
			lastKmStr = "    -  "
		}
		if tanken.Art != Erst {
			sumLiter += tanken.Liter
		}
		verbrauchSeitTanken := tanken.Liter / (float64(distanz) / 100.)
		verbrauchBisher := sumLiter / (float64(sumDistanz) / 100.)
		verbrauchSym := "➡"
		if verbrauchSeitTanken > verbrauchBisher {
			verbrauchSym = "⬆"
		} else if verbrauchSeitTanken < verbrauchBisher {
			verbrauchSym = "⬇"
		}
		fmt.Printf("    | %s | %8dkm | %8.2f € | %5.2fl | %5.3f €/l | %s | %6.2f %s | %7.2f |\n",
			tanken.Datum.Format("02.01.2006"),
			tanken.Km,
			tanken.Kosten,
			tanken.Liter,
			tanken.Kosten/tanken.Liter,
			lastKmStr,
			verbrauchSeitTanken,
			verbrauchSym,
			verbrauchBisher,
		)
		lastKm = tanken.Km
	}
}

//goland:noinspection GoUnhandledErrorResult
func (kfz *Kfz) printKosten() {
	fmt.Println("  Kosten - Aufstellung:")
	sort.SliceStable(kfz.Kosten, func(i, j int) bool {
		return false
	})
	abschreibungLen := 0
	abschreibungen := make([]kostenAbschreibung, 0, len(kfz.Kosten))
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
		l := len([]rune(abschreibung))
		if l > abschreibungLen {
			abschreibungLen = l
		}
		abschreibungen = append(abschreibungen, kostenAbschreibung{kosten: &kosten, abschreibung: abschreibung})
	}
	fmt.Printf("    | Datum      | Tachostand | Preis      | Kategorie       | %s | 🏦| Bemerkung            |\n",
		lang.FixedString("Abschreibung", abschreibungLen, ""))
	for _, line := range abschreibungen {
		fa := " "
		switch line.kosten.AbschreibungFa {
		case Ja:
			fa = "✓"
		case Nein:
			fa = "𐄂"
		case Abschreibung:
			fa = "⏲"
		}
		fmt.Printf("    | %s | %8dkm | %8.2f € | %s | %s | %s | %s |\n",
			line.kosten.Datum.Format("02.01.2006"),
			line.kosten.Km,
			line.kosten.Kosten,
			lang.FixedString(line.kosten.Kategorie, 15, "…"),
			lang.FixedString(line.abschreibung, abschreibungLen, "…"),
			fa,
			lang.FixedString(line.kosten.Notiz, 20, "…"),
		)
	}
}

//goland:noinspection GoUnhandledErrorResult
func (kfz *Kfz) printSummary() {
	minKm, minDatum := kfz.MinKm()
	maxKm, maxDatum := kfz.MaxKm()
	gefahrenKm := float64(maxKm - minKm)
	liter, km, kostenTanken := kfz.StatTanken()
	totalKosten, anteilKosten, anteilKostenFa := kfz.StatKosten()

	fmt.Printf("  Start:    %9d km (%s)\n", minKm, minDatum.Format("02.01.2006"))
	fmt.Printf("  Aktuell:  %9d km (%s)\n", maxKm, maxDatum.Format("02.01.2006"))
	tage := time.Now().Sub(minDatum).Hours() / 24
	fmt.Printf("  Gefahren: %9.0f km in %.0f Tagen (%.0f km/Tag)\n", gefahrenKm, tage, gefahrenKm/tage)
	fmt.Println("\n  Tanken:")
	fmt.Printf("    %9.2fl (%.2fl/100km)   | %10.2f€ (%.3f€/l | %.2f€/km)\n",
		liter, liter/(km/100.0), kostenTanken, kostenTanken/liter, kostenTanken/km)
	fmt.Println("  Kosten:")
	fmt.Printf("    Total: %9.2f€ (%.2f/km) | Anteilig: %9.2f€ (%.2f/km) | Finanzamt: %9.2f€ (%.2f/km)\n",
		totalKosten, totalKosten/gefahrenKm,
		anteilKosten, anteilKosten/gefahrenKm,
		anteilKostenFa, anteilKostenFa/gefahrenKm,
	)
	fmt.Println("\n  Gesamt:")
	fmt.Printf("    Total: %9.2f€ (%.2f/km) | Anteilig: %9.2f€ (%.2f/km) | Finanzamt: %9.2f€ (%.2f/km)\n",
		kostenTanken+totalKosten, (kostenTanken+totalKosten)/gefahrenKm,
		kostenTanken+anteilKosten, (kostenTanken+anteilKosten)/gefahrenKm,
		kostenTanken+anteilKostenFa, (kostenTanken+anteilKostenFa)/gefahrenKm,
	)
}

type kostenAbschreibung struct {
	kosten       *Kosten
	abschreibung string
}
