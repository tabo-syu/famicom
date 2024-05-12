//go:generate go run .
//go:generate gofmt -w ../../
package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/gocolly/colly"
)

const (
	success = 0
	failure = 1
)

func main() {
	if err := run(); err != nil {
		os.Exit(failure)
	}

	os.Exit(success)
}

//go:embed opcodes.go.tmpl
var opcodesTmpl string

func run() error {
	ops, err := scrape()
	if err != nil {
		return err
	}

	tmpl, err := template.New("opcodes").Parse(opcodesTmpl)
	if err != nil {
		return err
	}

	file, err := os.Create("../../internal/cpu/opcodes.gen.go")
	if err != nil {
		return err
	}

	if err := tmpl.Execute(file, ops); err != nil {
		return err
	}

	return nil
}

type status struct {
	C string
	Z string
	I string
	D string
	B string
	V string
	N string
}

type mode struct {
	Name   string
	Code   string
	Bytes  string
	Cycles string
}

type opcode struct {
	Name   string
	Status status
	Modes  []mode
}

func scrape() ([]opcode, error) {
	c := colly.NewCollector()

	// Fetch opcodes
	opcodes := []opcode{}

	table := 0
	opIndex := 0

	c.OnHTML("table", func(e *colly.HTMLElement) {
		table++

		if table == 114 {
			return
		}

		// opcode names table
		if table == 1 {
			e.ForEach("td a", func(i int, e *colly.HTMLElement) {
				opcodes = append(opcodes, opcode{
					Name:   strings.TrimSpace(e.Text),
					Status: status{},
					Modes:  []mode{},
				})
			})

			return
		}

		// status table
		if table%2 == 0 {
			e.ForEach("td", func(i int, e *colly.HTMLElement) {
				// behavior
				text := strings.TrimSpace(e.Text)
				if (i+1)%3 == 0 {
					switch (i + 1) / 3 {
					case 1:
						opcodes[opIndex].Status.C = text
					case 2:
						opcodes[opIndex].Status.Z = text
					case 3:
						opcodes[opIndex].Status.I = text
					case 4:
						opcodes[opIndex].Status.D = text
					case 5:
						opcodes[opIndex].Status.B = text
					case 6:
						opcodes[opIndex].Status.V = text
					case 7:
						opcodes[opIndex].Status.N = text
					default:
						panic("status")
					}
				}
			})

			return
		}

		// addressing mode table
		if table%2 == 1 {
			e.ForEach("tr", func(i int, e *colly.HTMLElement) {
				if i == 0 {
					return
				}

				// addressing mode
				mode := mode{}
				e.ForEach("td", func(i int, e *colly.HTMLElement) {
					text := strings.TrimSpace(e.Text)

					if i == 0 /* name */ {
						text := strings.NewReplacer(
							" ", "",
							",", "",
							"\n", "",
							"(", "",
							")", "",
							"$", "",
						).Replace(text)
						mode.Name = text
					} else /* code */ if i == 1 {
						mode.Code = strings.ReplaceAll(text, "$", "")
					} else /* bytes */ if i == 2 {
						mode.Bytes = text
					} else /* cycles */ if i == 3 {
						text := strings.NewReplacer(
							"      ", ", ",
							"\n", "",
							"(", "/*(",
							")", ")*/",
						).Replace(text)
						mode.Cycles = text
					}
				})
				opcodes[opIndex].Modes = append(opcodes[opIndex].Modes, mode)
			})

			opIndex++
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting URL:", r.URL.String())
	})

	if err := c.Visit("https://www.nesdev.org/obelisk-6502-guide/reference.html"); err != nil {
		return nil, err
	}

	return opcodes, nil
}
