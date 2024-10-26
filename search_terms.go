package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/chzyer/readline"
)

func get_user_input(msg string, required bool) []string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Fprint(os.Stdout, msg)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		err_fatal(err)
	}
	str := scanner.Text()
	if len(str) == 0 && required {
		err := fmt.Errorf("that field is required good sir")
		err_fatal(err)
	}
	q := parse_input(str)
	return q
}

func edit_user_input(msg string, prev string) (string, error) {
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("102"))
	input, _ := readline.New(s.Render(msg))
	defer input.Close()

	input.WriteStdin([]byte(prev))

	str, err := input.Readline()
	if err != nil {
		return "", err
	}
	return str, nil
}
