package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tusharameria/n2t/internal/translator"
)

func ProcessFile(path, fullName string, counters *translator.Counters, funcCounters map[string]int) (string, []string, error) {

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		return "", nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	finalText := ""
	parts := strings.Split(fullName, ".")
	fileName := parts[0]
	calledFileNames := []string{}

	for {
		if scanner.Scan() {
			i++
			text, comment := cleanText(scanner.Text())
			if text == "" {
				text = fmt.Sprintf("//%s", comment)
			} else {
				args, num, ok := isValidMemorySegCommand(text)
				if ok {
					buff, err := translator.MemorySegCommandTranslate(args, num, fileName)
					if err != nil {
						fmt.Printf("%s\n", err)
						return "", nil, err
					}
					text = buff
				}
				if isValidALCommand(text) {
					buff, err := translator.ArithmaticLogicCommandTranslate(text, counters)
					if err != nil {
						fmt.Printf("%s\n", err)
						return "", nil, err
					}
					text = fmt.Sprintf("// %s\n%s", text, buff)
				}
				argsBranch, ok := isValidBranchingCommand(text)
				if ok {
					text = translator.BranchingCommandTranslate(argsBranch, fileName)
				}
				argsFunction, ok := isValidFunctionCommand(text)
				if ok {
					text, calledFileNames, err = translator.FunctionCommandTranslate(argsFunction, fileName, funcCounters, calledFileNames)
					if err != nil {
						fmt.Printf("%s\n", err)
						return "", nil, err
					}
				}
			}
			finalText += text + "\n"
			// time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return finalText, calledFileNames, nil
}

func cleanText(line string) (string, string) {
	words := strings.SplitN(line, "//", 2)
	text := ""
	if len(words) == 2 {
		text = words[1]
	}
	return strings.TrimRight(strings.TrimLeft(words[0], " \t"), " \t"), text
}

func isValidMemorySegCommand(line string) ([]string, uint32, bool) {
	words := strings.Split(line, " ")
	if len(words) != 3 {
		return nil, 0, false
	}
	argOne := words[0]
	if _, okArgOne := translator.MemoryArgOne[argOne]; !okArgOne {
		return nil, 0, false
	}
	argTwo := words[1]
	if _, okArgTwo := translator.MemoryArgTwo[argTwo]; !okArgTwo {
		return nil, 0, false
	}
	argThree := words[2]
	num, err := strconv.Atoi(argThree)
	if err != nil {
		return nil, 0, false
	}

	return words[:2], uint32(num), true
}

func isValidBranchingCommand(line string) ([]string, bool) {
	words := strings.Split(line, " ")
	if len(words) != 2 {
		return nil, false
	}
	argOne := words[0]
	if _, okArgOne := translator.BranchingArgs[argOne]; !okArgOne {
		return nil, false
	}
	return words, true
}

func isValidFunctionCommand(line string) ([]string, bool) {
	words := strings.Split(line, " ")
	length := len(words)
	if length > 3 || length == 0 {
		return nil, false
	}
	argOne := words[0]
	if _, okArgOne := translator.FunctionArgs[argOne]; !okArgOne {
		return nil, false
	}
	if length == 1 && argOne != translator.RETURN {
		return nil, false
	}
	if length == 3 && argOne != translator.FUNCTION && argOne != translator.CALL {
		return nil, false
	}
	return words, true
}

func isValidALCommand(line string) bool {
	if _, ok := translator.ArithmaticLogicArgs[line]; !ok {
		return false
	}
	return true
}
