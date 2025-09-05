package output

import (
	"io"
	"strings"

	"github.com/mdryaan/vaultenv/internal/vault"
	"github.com/olekukonko/tablewriter"
)

// TableFormatter outputs entries as an ASCII table.
type TableFormatter struct{}

func (f *TableFormatter) WriteEntries(w io.Writer, entries []vault.Entry, showValues bool) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"KEY", "VALUE", "TAGS", "UPDATED"})
	table.SetBorder(true)
	table.SetAutoWrapText(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
	)

	for _, e := range entries {
		value := maskValue(e.Value)
		if showValues {
			value = e.Value
		}
		tags := strings.Join(e.Tags, ", ")
		updated := e.UpdatedAt.Format("2006-01-02 15:04")
		table.Append([]string{e.Key, value, tags, updated})
	}

	table.Render()
	return nil
}
