package termtable

import (
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

func Output(cols []string, rows [][]string) {
	OutputTo(cols, rows, os.Stdout)
}

func OutputTo(cols []string, rows [][]string, w io.Writer) {
	table := tablewriter.NewWriter(w)
	table.SetHeader(cols)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(rows) // Add Bulk Data
	table.Render()
}
