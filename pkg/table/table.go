package table

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func GetTableTemplate(title string, flagPageSize int, flagPageStyle string) (t table.Writer) {
	t = table.NewWriter()
	t.SetTitle(title)

	switch flagPageStyle {
	case "bright":
		t.SetStyle(table.StyleColoredBright)
	case "dark":
		t.SetStyle(table.StyleColoredCyanWhiteOnBlack)
	case "light":
		t.SetStyle(table.StyleLight)
	}

	t.Style().Title.Align = text.AlignCenter

	// t.SetColumnConfigs([]table.ColumnConfig{{Name: "Description", WidthMax: 80}})

	t.SetPageSize(flagPageSize)

	return t
}
