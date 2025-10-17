package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const PUSH = "push"
const POP = "pop"

const LOCAL = "local"
const ARGUMENT = "argument"
const THIS = "this"
const THAT = "that"
const CONSTANT = "constant"
const STATIC = "static"
const POINTER = "pointer"
const TEMP = "temp"

const TEMP_INIT_ADDRESS = 5

const ADD = "add"
const SUB = "sub"
const NEG = "neg"
const EQ = "eq"
const GT = "gt"
const LT = "lt"

var argOnes = map[string]bool{
	PUSH: true,
	POP:  true,
}
var argTwos = map[string]string{
	LOCAL:    "LCL",
	ARGUMENT: "ARG",
	THIS:     "THIS",
	THAT:     "THAT",
	CONSTANT: "CONST",
	STATIC:   "STATIC",
	POINTER:  "POINT",
	TEMP:     "R5",
}

var alArgs = map[string]string{
	ADD: "@SP\nM=M-1\nA=M\nD=M\nA=A-1\nM=D+M\n",
	SUB: "@SP\nM=M-1\nA=M\nD=M\nA=A-1\nM=M-D\n",
	NEG: "@SP\nA=M-1\nM=-M\n",
	EQ:  "",
	GT:  "",
	LT:  "",
}

func main() {
	now := time.Now()
	fmt.Println("Starting Translator...")

	file, err := os.Open("tests/07/PointerTest/PointerTest.vm")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer file.Close()

	filePath := file.Name()
	fileNameExt := strings.Split(filePath, "/")
	fileName := strings.Split(fileNameExt[len(fileNameExt)-1], ".")[0]

	outputDir := "output"

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	outFile, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("%s.asm", fileName)))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer outFile.Close()

	done := make(chan struct{})
	messages := make(chan string, 1000)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	var once sync.Once
	shutdown := func() {
		once.Do(func() {
			close(done)
		})
	}

	go func() {
		scanner := bufio.NewScanner(file)
		i := 0
		for {
			if scanner.Scan() {
				i++
				text := scanner.Text()
				args, num, ok := isValidMemorySegCommand(text)
				if ok {
					buff, err := translateMemorySegCommant(args, num)
					if err != nil {
						fmt.Printf("%s\n", err)
						return
					}
					text = buff
				}
				if isValidALCommand(text) {
					buff, ok := alArgs[text]
					if !ok {
						fmt.Printf("%s is not a valid alArgs key\n", text)
						return
					}
					text = fmt.Sprintf("// %s\n%s", text, buff)
				}
				messages <- text + "\n"
				// time.Sleep(1 * time.Second)
			} else {
				close(messages)
				shutdown()
				return
			}
		}
	}()

	writer := bufio.NewWriter(outFile)
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Printf("%s\n", err)
		}
	}()

	for msg := range messages {
		_, err := writer.WriteString(msg)
		if err != nil {
			fmt.Printf("write error at line: %s\nerr : %v\n", msg, err)
			shutdown()
			return
		}
	}

	go func() {
		sig := <-sigCh
		fmt.Printf("\nReceived signal : %v\n", sig)
		shutdown()
	}()

	<-done

	signal.Stop(sigCh)
	fmt.Println("Closing Translator...")
	fmt.Printf("%d us\n", time.Now().Sub(now).Microseconds())
	return
}

func isValidMemorySegCommand(line string) ([]string, uint32, bool) {
	words := strings.Split(line, " ")
	if len(words) != 3 {
		return nil, 0, false
	}
	argOne := words[0]
	if _, okArgOne := argOnes[argOne]; !okArgOne {
		return nil, 0, false
	}
	argTwo := words[1]
	if _, okArgTwo := argTwos[argTwo]; !okArgTwo {
		return nil, 0, false
	}
	argThree := words[2]
	num, err := strconv.Atoi(argThree)
	if err != nil {
		return nil, 0, false
	}

	return words[:2], uint32(num), true
}

func translateMemorySegCommant(args []string, num uint32) (string, error) {
	argOne := args[0]
	argTwo := args[1]

	// push local 2
	// @2
	// D=A
	// @LCL
	// A=D+M
	// D=M
	// @SP
	// M=M+1
	// A=M
	// A=A-1
	// M=D

	// push temp 2
	// @2
	// D=A
	// @5
	// A=D+A
	// D=M
	// @SP
	// M=M+1
	// A=M
	// A=A-1
	// M=D

	str := ""
	if argOne == PUSH {
		argTwoVal, ok := argTwos[argTwo]
		if !ok {
			return "", fmt.Errorf("%s doesn't exist in argTwos", argTwo)
		}
		if argTwo == POINTER {
			buff := argTwos[THIS]
			if num == 1 {
				buff = argTwos[THAT]
			}
			str += fmt.Sprintf("@%s\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", buff)
		} else {
			str += fmt.Sprintf("@%d\nD=A\n", num)
			if argTwo != CONSTANT {
				if argTwo == TEMP {
					str += fmt.Sprintf("@%d\nA=D+A\nD=M\n", TEMP_INIT_ADDRESS)
				} else {
					str += fmt.Sprintf("@%s\nA=D+M\nD=M\n", argTwoVal)
				}
			}
			str += fmt.Sprintf("@SP\nM=M+1\nA=M\nA=A-1\nM=D\n")
		}
	}

	// pop local 3
	// @3
	// D=A
	// @LCL
	// A=D+M
	// D=A
	// @R16
	// M=D
	// @SP
	// M=M-1
	// A=M
	// D=M
	// @R16
	// A=M
	// M=D

	// pop temp 3
	// @3
	// D=A
	// @5
	// D=D+A
	// @R16
	// M=D
	// @SP
	// M=M-1
	// A=M
	// D=M
	// @R16
	// A=M
	// M=D

	if argOne == POP {
		if argTwo == CONSTANT {
			return "", fmt.Errorf("Can't POP a constant")
		}
		argTwoVal, ok := argTwos[argTwo]
		if !ok {
			return "", fmt.Errorf("%s doesn't exist in argTwos", argTwo)
		}

		if argTwo == POINTER {
			buff := argTwos[THIS]
			if num == 1 {
				buff = argTwos[THAT]
			}
			str += fmt.Sprintf("@SP\nM=M-1\n@SP\nA=M\nD=M\n@%s\nM=D\n", buff)
		} else {
			str += fmt.Sprintf("@%d\nD=A\n", num)
			if argTwo == TEMP {
				str += fmt.Sprintf("@%d\nD=D+A\n", TEMP_INIT_ADDRESS)
			} else {
				str += fmt.Sprintf("@%s\nA=D+M\nD=A\n", argTwoVal)
			}
			str += fmt.Sprintf("@R16\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R16\nA=M\nM=D\n")
		}
	}

	return fmt.Sprintf("// %s %s %d\n%s", argOne, argTwo, num, str), nil
}

func isValidALCommand(line string) bool {
	if _, ok := alArgs[line]; !ok {
		return false
	}
	return true
}
