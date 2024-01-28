package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: move to util.go
func PrintDesc(msg string) string {
	return fmt.Sprintf(">>>>> %s <<<<<", msg)
}

// if succeed, return string end with "\n"
func Readline() (string, error) {
	input := bufio.NewReader(os.Stdin)
	line, err := input.ReadString('\n')
	if err != nil {
		fmt.Println("input.ReadString err:", err)
		return "", err
	}
	return line[:len(line)-1], nil
}
