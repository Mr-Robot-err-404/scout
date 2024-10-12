package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
)

func DuckDuckGooseEnumerator(l list.Items, i int) string {
	if l.At(i).Value() == "Goose" {
		return "Honk →"
	}
	return ""
}

func juicy_print() {
	l := list.New("Duck", "Duck", "Duck", "Duck", "Goose", "Duck", "Duck")
	l = l.Enumerator(DuckDuckGooseEnumerator)
	fmt.Println(l)
}

func print_table(headers []string, rows [][]string) {
	if len(rows) == 0 {
		info_msg_fatal("table is empty")
	}
	HeaderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	EvenRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("110"))
	OddRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("69"))

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return HeaderStyle
			case row%2 == 0:
				return EvenRowStyle
			default:
				return OddRowStyle
			}
		}).
		Headers(headers...).
		Rows(rows...)

	fmt.Println(t)
}
