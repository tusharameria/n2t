package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tusharameria/n2t/internal/filewalker"
	"github.com/tusharameria/n2t/internal/generator"
	"github.com/tusharameria/n2t/internal/parser"
	"github.com/tusharameria/n2t/internal/translator"
)

func main() {
	now := time.Now()
	fmt.Println("Starting Translator...")
	outputFileName := ""
	outputDir := "output"
	inputPath := ""
	isDirectory := false
	initialFileName := "Sys.vm"
	funcCounters := make(map[string]int)
	counters := translator.Counters{}

	flag.StringVar(&inputPath, "input", inputPath, "Path to the input file/directory")
	flag.Parse()

	info, err := os.Stat(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Input path does not exist: %s\n", inputPath)
			return
		}
		fmt.Printf("Error accessing input path: %s\n", err)
		return
	}

	pathChunks := strings.Split(inputPath, "/")
	lastName := pathChunks[len(pathChunks)-1]

	if info.IsDir() {
		isDirectory = true
		outputFileName = lastName
	} else {
		fileNameChunks := strings.Split(lastName, ".")
		ext := fileNameChunks[len(fileNameChunks)-1]
		if ext != "vm" {
			fmt.Printf("Input file is not a .vm file: %s\n", inputPath)
			return
		}
		outputFileName = fileNameChunks[0]
	}

	msg := ""

	if !isDirectory {
		msg, _, err = parser.ProcessFile(inputPath, outputFileName, &counters, funcCounters)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	} else {
		tree, err := filewalker.BuildDirTree(inputPath)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		var initFileInfo *filewalker.DirTree

		for i := 0; i < len(tree.Children); i++ {
			if tree.Children[i].Name == initialFileName {
				initFileInfo = tree.Children[i]
			}
		}

		if initFileInfo == nil {
			fmt.Println("no Sys.vm file present")
			return
		}

		calledFiles := make(map[string]bool)
		queue := []string{}

		msg += translator.GenerateBootstrapCode(strings.Replace(initialFileName, "vm", "init", -1))

		res, calledFileNames, err := parser.ProcessFile(initFileInfo.Path, initFileInfo.Name, &counters, funcCounters)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		msg += res

		for _, name := range calledFileNames {
			if _, ok := calledFiles[name]; !ok {
				queue = append(queue, name)
				calledFiles[name] = true
			}
		}

		for len(queue) > 0 {
			firstFileName := queue[0]
			queue = queue[1:]

			res, newFiles, err := parser.ProcessFile(filepath.Join(inputPath, fmt.Sprintf("%s.vm", firstFileName)), fmt.Sprintf("%s.vm", firstFileName), &counters, funcCounters)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}

			msg += res
			for _, name := range newFiles {
				if _, ok := calledFiles[name]; !ok {
					queue = append(queue, name)
					calledFiles[name] = true
				}
			}
		}

	}

	if err := generator.WriteOutput(outputDir, outputFileName, msg); err != nil {
		fmt.Printf("error in generating output : %v\n", err)
		return
	}

	fmt.Println("Closing Translator...")
	fmt.Printf("%d us\n", time.Since(now).Microseconds())
}
