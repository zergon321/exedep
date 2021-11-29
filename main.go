package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zergon321/arrayqueue"
)

var (
	dumpbinPath string
	exePath     string
	dllPath     string
)

func parseFlags() {
	flag.StringVar(&dumpbinPath, "dumpbin", "C:\\Program Files\\Microsoft Visual Studio"+
		"\\2022\\Community\\VC\\Tools\\MSVC\\14.30.30705\\bin\\Hostx64\\x64\\dumpbin.exe",
		"An absolute path to the 'dumpbin' command line utility")
	flag.StringVar(&exePath, "exe", "", "A path to the executable file to be analyzed")
	flag.StringVar(&dllPath, "dll", "C:\\msys64\\mingw64\\bin",
		"A path to the directory where all the possible DLLs are stored")

	flag.Parse()
}

func createGetDependentsCommand(dumpbinPath string) string {
	return "& \"" + dumpbinPath + "\" /dependents"
}

func getDependents(targetPath string) ([]string, error) {
	if !filepath.IsAbs(targetPath) {
		var err error
		targetPath, err = filepath.Abs(targetPath)

		if err != nil {
			return nil, err
		}
	}

	buffer := bytes.NewBuffer([]byte{})
	cmd := exec.Command("powershell", "-Command",
		createGetDependentsCommand(dumpbinPath)+" "+targetPath)
	cmd.Stdout = buffer

	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	out := buffer.String()

	return getDLLsFromOutput(out), nil
}

func getDLLsFromOutput(out string) []string {
	lines := strings.Split(out, "\n")
	dlls := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasSuffix(line, ".dll") {
			dlls = append(dlls, line)
		}
	}

	return dlls
}

func containsDLL(dllEntries []fs.FileInfo, dllName string) bool {
	contains := false

	for _, dllEntry := range dllEntries {
		if dllEntry.Name() == dllName {
			contains = true
			break
		}
	}

	return contains
}

func contains(strs []string, str string) bool {
	contains := false

	for _, strEntry := range strs {
		if strEntry == str {
			contains = true
			break
		}
	}

	return contains
}

func main() {
	parseFlags()
	var err error

	entries, err := ioutil.ReadDir(dllPath)
	handleError(err)
	dllEntries := []fs.FileInfo{}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".dll") {
			dllEntries = append(dllEntries, entry)
		}
	}

	directDependencies, err := getDependents(exePath)
	handleError(err)
	inspectQueue, err := arrayqueue.NewQueue(64)
	handleError(err)

	for _, dependency := range directDependencies {
		inspectQueue.Enqueue(dependency)
	}

	bundleDLLs := []string{}

	for inspectQueue.Length() > 0 {
		item, err := inspectQueue.Dequeue()
		handleError(err)
		dependency, ok := item.(string)

		if !ok {
			handleError(fmt.Errorf(
				"wrong dependency item type: %s", item))
		}

		if contains(bundleDLLs, dependency) {
			continue
		}

		if containsDLL(dllEntries, dependency) {
			bundleDLLs = append(bundleDLLs, dependency)
			fmt.Println(dependency)
		} else {
			continue
		}

		dependencyPath := filepath.Join(dllPath, dependency)
		dependencies, err := getDependents(dependencyPath)
		handleError(err)

		for _, dependency := range dependencies {
			inspectQueue.Enqueue(dependency)
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
