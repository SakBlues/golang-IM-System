package pkg

import (
	"bufio"
	"fmt"
	"os"
)

// read a string end with "\n",
// and return result will remove the last "\n".
// use bufio.Reader to read a line,
// because fmt.scanln will return when reach a blank.
func Readline() (string, error) {
	input := bufio.NewReader(os.Stdin)
	line, err := input.ReadString('\n')
	if err != nil {
		fmt.Println("input.ReadString err:", err)
		return "", err
	}
	return line[:len(line)-1], nil
}