package generator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func WriteOutput(outputDir, outputFileName, msg string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	outFile, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("%s.asm", outputFileName)))
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Printf("%s\n", err)
		}
	}()

	_, err = writer.WriteString(msg)
	if err != nil {
		fmt.Printf("write error at line: %s\nerr : %v\n", msg, err)
		return err
	}
	return nil
}
