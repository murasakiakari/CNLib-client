package utility

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text, _ := reader.ReadString('\r')
	text = strings.Trim(text, "\r")
	return text
}
