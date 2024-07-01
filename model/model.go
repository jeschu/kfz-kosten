package model

import (
	cfg "github.com/jeschu/go-config"
	"math"
	"time"
)

type Kfzs map[string]Kfz

type Kfz struct {
	Name        string   `yaml:"name"`
	Kennzeichen string   `yaml:"kennzeichen"`
	Kosten      []Kosten `yaml:"kosten"`
	Tanken      []Tanken `yaml:"tanken"`
}

type Kosten struct {
	Datum            time.Time     `yaml:"datum"`
	Km               int           `yaml:"km"`
	Kategorie        string        `yaml:"kategorie"`
	AbschreibungKm   int           `yaml:"abschreibung_km"`
	AbschreibungZeit time.Duration `yaml:"abschreibung_zeit"`
	AbschreibungFa   bool          `yaml:"abschreibung_fa"`
	Kosten           float64       `yaml:"kosten"`
	Notiz            string        `yaml:"notiz"`
}

type TankArt int

const (
	Erst TankArt = iota
	Teil
	Voll
)

type Tanken struct {
	Datum  time.Time `yaml:"datum"`
	Art    TankArt   `yaml:"art"`
	Km     int       `yaml:"km"`
	Liter  float64   `yaml:"liter"`
	Kosten float64   `yaml:"kosten"`
	Sorte  string    `yaml:"sorte"`
}

func (tanken Tanken) Len() int {
	return 0
}

func (kfz *Kfz) MaxKm() (int, time.Time) {
	kst := kfz.Kosten[0]
	km := float64(kst.Km)
	datum := kst.Datum
	for _, kosten := range kfz.Kosten {
		km = math.Max(km, float64(kosten.Km))
		if kosten.Datum.After(datum) {
			datum = kosten.Datum
		}
	}
	for _, tanken := range kfz.Tanken {
		km = math.Max(km, float64(tanken.Km))
		if tanken.Datum.After(datum) {
			datum = tanken.Datum
		}
	}
	return int(km), datum
}

func (kfz *Kfz) MinKm() (int, time.Time) {
	kst := kfz.Kosten[0]
	km := float64(kst.Km)
	datum := kst.Datum
	for _, kosten := range kfz.Kosten {
		km = math.Min(km, float64(kosten.Km))
		if kosten.Datum.Before(datum) {
			datum = kosten.Datum
		}
	}
	for _, tanken := range kfz.Tanken {
		km = math.Min(km, float64(tanken.Km))
		if tanken.Datum.Before(datum) {
			datum = tanken.Datum
		}
	}
	return int(km), datum
}

func (kfz *Kfz) StatTanken() (float64, float64, float64) {
	liter := 0.0
	kosten := 0.0
	kmMin := float64(kfz.Tanken[0].Km)
	kmMax := kmMin
	for _, tanken := range kfz.Tanken {
		liter += tanken.Liter
		kosten += tanken.Kosten
		kmMin = math.Min(kmMin, float64(tanken.Km))
		kmMax = math.Max(kmMax, float64(tanken.Km))
	}
	return liter, kmMax - kmMin, kosten
}

func (kfz *Kfz) StatKosten() (float64, float64, float64) {
	heute := time.Now().Truncate(24 * time.Hour)
	kmInt, _ := kfz.MaxKm()
	km := float64(kmInt)
	anteilKosten := 0.0
	anteilKostenFa := 0.0
	totalKosten := 0.0
	for _, kst := range kfz.Kosten {
		kosten := kst.Kosten
		totalKosten += kosten
		abschreibungKm := float64(kst.AbschreibungKm)
		abschreibungZeit := kst.AbschreibungZeit
		if abschreibungZeit > 0 {
			d := heute.Sub(kst.Datum)
			if d > abschreibungZeit {
				anteilKosten += kosten
				if kst.AbschreibungFa {
					anteilKostenFa += kosten
				}
			} else {
				anteil := kosten / float64(abschreibungZeit/d)
				anteilKosten += anteil
				if kst.AbschreibungFa {
					anteilKostenFa += anteil
				}
			}
		} else if abschreibungKm > 0.0 {
			kostenStart := float64(kst.Km)
			if kostenStart+abschreibungKm > km {
				anteilKosten += kosten
				if kst.AbschreibungFa {
					anteilKostenFa += kosten
				}
			} else {
				kstKm := float64(kst.Km)
				anteilKm := km - kstKm
				anteil := kosten / (anteilKm / abschreibungKm)
				anteilKosten += anteil
				if kst.AbschreibungFa {
					anteilKostenFa += anteil
				}
			}
		} else {
			anteilKosten += kosten
			if kst.AbschreibungFa {
				anteilKostenFa += kosten
			}
		}
	}
	return totalKosten, anteilKosten, anteilKostenFa
}

func LoadKfzs() (Kfzs, error) {
	kfzs := Kfzs{}
	err := cfg.ReadConfigYaml("kfzs.yaml", &kfzs)
	return kfzs, err
}

func SaveKfzs(kfzs Kfzs) error {
	return cfg.WriteConfigYaml("kfzs.yaml", &kfzs)
}
