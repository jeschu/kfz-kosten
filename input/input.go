package input

import (
	"bufio"
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var CR = string([]byte{13})

func ReadString(msg string) string {
	for {
		fmt.Print(msg)
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println(err)
		} else {
			return input[:len(input)-1]
		}
	}
}

func Read[T any](msg string, mapper func(string) (T, error)) T {
	for {
		fmt.Print(msg)
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println(err)
		} else {
			if v, e := mapper(input[:len(input)-1]); e == nil {
				return v
			} else {
				log.Println(e)
			}
		}
	}
}

func ReadSelectionMapped[T any](msg string, mapping map[string]T, defaultValue T, options ...string) T {
	mapper := func(s string) T {
		if v, ok := mapping[s]; ok {
			return v
		} else {
			return defaultValue
		}
	}
	return ReadSelectionMapper(msg, mapper, options...)
}

func ReadSelectionMapper[T any](msg string, mapper func(string) T, options ...string) T {
	selection := readSelection(msg, options...)
	mapped := mapper(selection)
	fmt.Printf("\033[1A%s%s\n", msg, mapped)
	return mapped
}

func ReadSelection(msg string, options ...string) string {
	selection := readSelection(msg, options...)
	fmt.Printf("\033[1A%s%s\n", msg, selection)
	return selection
}

func readSelection(msg string, options ...string) string {
	var input string
	var b = make([]byte, 1)
	for {
		fmt.Println(msg)
		oldState := makeRaw()
		_, err := os.Stdin.Read(b)
		restore(oldState)
		if err != nil {
			log.Print(err)
		} else {
			input = strings.ToLower(string(b))
			for _, valid := range options {
				if input == strings.ToLower(valid) {
					return input
				}
			}
			fmt.Printf("ungültige Eingabe: '%s'\n\n", input)
		}
	}
}

func ReadDateInPast(msg string) time.Time {
	for {
		date := startOfDay(time.Now())
		day, month, year := splitDate(date)
		message := msg + " ([T][.M][.J]): "
		fmt.Print(message)

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println(err)
		} else {
			parts := strings.Split(input[:len(input)-1], ".")
			if parts[0] == "" {
				return date
			}
			day, err = strconv.Atoi(parts[0])
			if err != nil {
				fmt.Printf("ungültiger Tag: '%s'\n", parts[0])
				continue
			}
			if len(parts) > 1 {
				month, err = strconv.Atoi(parts[1])
				if err != nil {
					fmt.Printf("ungültiger Monat: '%s'\n", parts[1])
					continue
				}
			}
			if len(parts) > 2 {
				year, err = strconv.Atoi(parts[2])
				if err != nil {
					fmt.Printf("ungültiges Jahr: '%s'\n", parts[2])
					continue
				}
			}
			inputDate := newDate(year, month, day)
			for inputDate.After(date) {
				if inputDate.Day() > date.Day() {
					inputDate = inputDate.AddDate(0, -1, 0)
				}
				if inputDate.Month() > date.Month() {
					inputDate = inputDate.AddDate(-1, 0, 0)
				}
			}
			fmt.Printf("\033[1A%s%s\n", msg, inputDate.Format("02.01.2006"))
			return inputDate
		}
	}
}

func splitDate(date time.Time) (d int, m int, y int) {
	return date.Day(), int(date.Month()), date.Year()
}

func newDate(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func startOfDay(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

}

func makeRaw() *term.State {
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	return state
}

func restore(state *term.State) {
	err := term.Restore(int(os.Stdin.Fd()), state)
	if err != nil {
		log.Fatal(err)
	}
}
