package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
)

var sp = create_spinner()

func create_spinner() *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Color("yellow")
	return s
}

func load(msg string) {
	sp.Suffix = " " + msg
	sp.Start()
}

func success_msg(msg string) {
	sp.Stop()
	t := "\u2713 "
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	fmt.Println(s.Render(t) + msg)
}
func success_resp() {
	d := "DONE"
	s := lipgloss.NewStyle().Background(lipgloss.Color("22")).PaddingRight(1).PaddingLeft(1).MarginRight(1)
	fmt.Println(s.Render(d))
}

func err_msg(msg string) {
	sp.Stop()
	t := "\u2717 "
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	fmt.Println(s.Render(t) + msg)
}

func log_err_queue(queue []error) {
	for _, err := range queue {
		err_resp(err)
	}
}

func err_resp(err error) {
	e := "ERROR"
	s := lipgloss.NewStyle().Background(lipgloss.Color("1")).PaddingRight(1).PaddingLeft(1).MarginRight(1)
	fmt.Printf(s.Render(e)+"%v\n", err)
}
func err_fatal(err error) {
	err_resp(err)
	os.Exit(0)
}

func info_msg(msg string) {
	t := "\u276F "
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	fmt.Println(s.Render(t) + msg)
}

func info_msg_fatal(msg string) {
	info_msg(msg)
	os.Exit(0)
}

func print_title(msg string) {
	s := lipgloss.NewStyle().Bold(true)
	fmt.Println(s.Render(strings.ToUpper(msg)))
}

func print_title_with_bg(msg string) {
	s := lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("4")).PaddingRight(1).PaddingLeft(1).MarginRight(1)
	fmt.Println(s.Render(strings.ToUpper(msg)))
}
