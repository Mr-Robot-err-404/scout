package main

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
)

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

func err_resp(err error) {
	e := "ERROR"
	s := lipgloss.NewStyle().Background(lipgloss.Color("1")).PaddingRight(1).PaddingLeft(1).MarginRight(1)
	fmt.Printf(s.Render(e)+"%v\n", err)
}
func err_fatal(err error) {
	e := "ERROR"
	s := lipgloss.NewStyle().Background(lipgloss.Color("1")).PaddingRight(1).PaddingLeft(1).MarginRight(1)
	fmt.Printf(s.Render(e)+"%v\n", err)
	os.Exit(0)
}

func info_msg(msg string) {
	t := "\u276F "
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	fmt.Println(s.Render(t) + msg)

}

func logging_time() {
	s := create_spinner()
	msg := "fetch channels"
	load(msg)
	time.Sleep(4 * time.Second)
	s.Stop()
	success_msg(msg)

	msg = "create playlist"
	load(msg)
	time.Sleep(4 * time.Second)
	s.Stop()
	success_msg(msg)

	msg = "insert items"
	load(msg)
	time.Sleep(4 * time.Second)
	s.Stop()
	err_msg(msg)
	err := fmt.Errorf("failed to insert items")
	err_fatal(err)

}
