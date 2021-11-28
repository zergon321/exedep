package main

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
)

var (
	dumpbinPath string
	exePath     string
)

func parseFlags() {
	flag.StringVar(&dumpbinPath, "dumpbin", "C:\\Program Files\\Microsoft Visual Studio"+
		"\\2022\\Community\\VC\\Tools\\MSVC\\14.30.30705\\bin\\Hostx64\\x64\\dumpbin.exe",
		"An absolute path to the 'dumpbin' command line utility")
	flag.StringVar(&exePath, "exe", "", "A path to the executable file to be analyzed")

	flag.Parse()
}

func createGetDependentsCommand(dumpbinPath string) string {
	return "& \"" + dumpbinPath + "\" /dependents"
}

func main() {
	parseFlags()
	var err error
	exePath, err = filepath.Abs(exePath)
	handleError(err)

	buffer := bytes.NewBuffer([]byte{})
	cmd := exec.Command("powershell", "-Command",
		createGetDependentsCommand(dumpbinPath)+" "+exePath)
	cmd.Stdout = buffer

	err = cmd.Run()
	handleError(err)

	fmt.Println(buffer.String())
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
