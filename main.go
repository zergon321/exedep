package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func main() {
	buffer := bytes.NewBuffer([]byte{})
	cmd := exec.Command("powershell", "-Command", "echo \"Lol\"")
	cmd.Stdout = buffer

	err := cmd.Run()
	handleError(err)

	fmt.Println(string(buffer.Bytes()))
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
