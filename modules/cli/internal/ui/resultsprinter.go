package ui

import (
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"
)

// Printable is an interface that defines the methods required to print data.
type Printable interface {
	JSON() ([]byte, error)
	TableHeader() []string
	TableRows() [][]string
}

// ResultsPrinter handles printing of any Printable data.
type ResultsPrinter struct {
	Data      Printable
	PrintJSON bool
}

// Render executes the appropriate rendering method based on the PrintJSON flag.
func (rp *ResultsPrinter) Render() {
	if rp.PrintJSON {
		jsonData, err := rp.Data.JSON()
		if err != nil {
			log.WithError(err).Error("Error marshalling to JSON")
			return
		}
		pterm.Println(string(jsonData))
	} else {
		header := rp.Data.TableHeader()
		rows := rp.Data.TableRows()
		data := append([][]string{header}, rows...)
		output, err := pterm.DefaultTable.WithHasHeader(true).WithSeparator(" ").WithData(data).Srender()
		if err != nil {
			log.WithError(err).Error("Error rendering table")
			return
		}
		pterm.Println(output)
	}
}
