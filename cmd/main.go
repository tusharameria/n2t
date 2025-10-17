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
const AND = "and"
const OR = "or"
const NOT = "not"

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

const EQ_TRUE = "EQ_TRUE"
const EQ_END = "EQ_END"
const GT_TRUE = "GT_TRUE"
const GT_END = "GT_END"
const LT_TRUE = "LT_TRUE"
const LT_END = "LT_END"

var alArgs = map[string]string{
	ADD: "@SP\nM=M-1\nA=M\nD=M\nA=A-1\nM=D+M\n",
	SUB: "@SP\nM=M-1\nA=M\nD=M\nA=A-1\nM=M-D\n",
	NEG: "@SP\nA=M-1\nM=-M\n",
	EQ:  fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\nA=A-1\nD=M-D\n@%s\nD;JEQ\n@SP\nA=M-1\nM=0\n@%s\n0;JMP\n(%s)\n@SP\nA=M-1\nM=-1\n(%s)\n", EQ_TRUE, EQ_END, EQ_TRUE, EQ_END),
	GT:  fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\nA=A-1\nD=M-D\n@%s\nD;JGT\n@SP\nA=M-1\nM=0\n@%s\n0;JMP\n(%s)\n@SP\nA=M-1\nM=-1\n(%s)\n", GT_TRUE, GT_END, GT_TRUE, GT_END),
	LT:  fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\nA=A-1\nD=M-D\n@%s\nD;JLT\n@SP\nA=M-1\nM=0\n@%s\n0;JMP\n(%s)\n@SP\nA=M-1\nM=-1\n(%s)\n", LT_TRUE, LT_END, LT_TRUE, LT_END),
	AND: "@SP\nM=M-1\nA=M\nD=M\nA=A-1\nM=D&M\n",
	OR:  "@SP\nM=M-1\nA=M\nD=M\nA=A-1\nM=D|M\n",
	NOT: "@SP\nA=M-1\nM=!M\n",
}

func main() {
	now := time.Now()
	fmt.Println("Starting Translator...")

	file, err := os.Open("tests/07/StackTest/StackTest.vm")
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
		eqCount := 0
		gtCount := 0
		ltCount := 0
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
					if text == EQ || text == GT || text == LT {
						if text == EQ {
							buff = strings.ReplaceAll(buff, EQ_TRUE, fmt.Sprintf("EQ_TRUE_%d", eqCount))
							buff = strings.ReplaceAll(buff, EQ_END, fmt.Sprintf("EQ_END_%d", eqCount))
							eqCount++
						}
						if text == GT {
							buff = strings.ReplaceAll(buff, GT_TRUE, fmt.Sprintf("GT_TRUE_%d", gtCount))
							buff = strings.ReplaceAll(buff, GT_END, fmt.Sprintf("GT_END_%d", gtCount))
							gtCount++
						}
						if text == LT {
							buff = strings.ReplaceAll(buff, LT_TRUE, fmt.Sprintf("LT_TRUE_%d", ltCount))
							buff = strings.ReplaceAll(buff, LT_END, fmt.Sprintf("LT_END_%d", ltCount))
							ltCount++
						}
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
