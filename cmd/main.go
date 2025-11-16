package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

const LABEL = "label"
const GO_TO = "goto"
const IF_GO_TO = "if-goto" // Pop the topmost value from the stack, and if that value is not zero, jump to the label.

var branchingArgs = map[string]bool{
	LABEL:    true,
	GO_TO:    true,
	IF_GO_TO: true,
}

const FUNCTION = "function"
const RETURN = "return"
const CALL = "call"

var functionArgs = map[string]bool{
	FUNCTION: true,
	RETURN:   true,
	CALL:     true,
}

type DirTree struct {
	Name     string
	Path     string
	Children []*DirTree
}

func main() {
	now := time.Now()
	fmt.Println("Starting Translator...")
	outputFileName := ""
	outputDir := "output"
	inputPath := ""
	isDirectory := false
	initialFileName := "Sys.vm"
	funcCounters := make(map[string]int)

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
		msg, _, err = processFile(inputPath, outputFileName, false, funcCounters)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
	} else {
		tree, err := buildDirTree(inputPath)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		var initFileInfo *DirTree

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

		// Bootstrap the initial SP assembly code
		buffName := strings.Replace(initialFileName, "vm", "init", -1)
		msg += "// Bootstrap Code\n"
		msg += "@256\nD=A\n@SP\nM=D\n"
		msg += fmt.Sprintf("@%s.RETURN.%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n", buffName, 0)
		msg += fmt.Sprintf("@LCL\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n")
		msg += fmt.Sprintf("@ARG\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n")
		msg += fmt.Sprintf("@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n")
		msg += fmt.Sprintf("@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n")
		msg += fmt.Sprintf("@%d\nD=A\n@5\nD=D+A\n@SP\nD=M-D\n@ARG\nM=D\n\n", 0)
		msg += fmt.Sprintf("@SP\nD=M\n@LCL\nM=D\n\n")
		msg += fmt.Sprintf("@%s\n0;JMP\n\n(%s.RETURN.%d)\n", buffName, buffName, 0)

		res, calledFileNames, err := processFile(initFileInfo.Path, initFileInfo.Name, true, funcCounters)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		msg += res

		fmt.Printf("calledFileNames : %s\n", calledFileNames)

		for _, name := range calledFileNames {
			if _, ok := calledFiles[name]; !ok {
				queue = append(queue, name)
				calledFiles[name] = true
			}
		}

		for len(queue) > 0 {
			firstFileName := queue[0]
			queue = queue[1:]

			res, newFiles, err := processFile(filepath.Join(inputPath, fmt.Sprintf("%s.vm", firstFileName)), fmt.Sprintf("%s.vm", firstFileName), false, funcCounters)
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

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	outFile, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("%s.asm", outputFileName)))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
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
		return
	}

	fmt.Println("Closing Translator...")
	fmt.Printf("%d us\n", time.Now().Sub(now).Microseconds())
	return
}

func processFile(path, fullName string, isSys bool, funcCounters map[string]int) (string, []string, error) {

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("%s\n", err)
		return "", nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	eqCount := 0
	gtCount := 0
	ltCount := 0
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
					buff, err := translateMemorySegCommand(args, num)
					if err != nil {
						fmt.Printf("%s\n", err)
						return "", nil, err
					}
					text = buff
				}
				if isValidALCommand(text) {
					buff, ok := alArgs[text]
					if !ok {
						fmt.Printf("%s is not a valid alArgs key\n", text)
						return "", nil, fmt.Errorf("%s is not a valid alArgs key", text)
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
				argsBranch, ok := isValidBranchingCommand(text)
				if ok {
					argOne := argsBranch[0]
					argTwo := argsBranch[1]
					text = fmt.Sprintf("//%s %s\n", argOne, argTwo)
					if argOne == LABEL {
						text += fmt.Sprintf("(%s.%s)\n", fileName, argTwo)
					} else if argOne == GO_TO {
						text += fmt.Sprintf("@%s.%s\n0;JMP\n", fileName, argTwo)
					} else {
						text += fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\n@%s.%s\nD;JNE\n", fileName, argTwo)
					}
				}
				argsFunction, ok := isValidFunctionCommand(text)
				if ok {
					argOne := argsFunction[0]
					text = fmt.Sprintf("// %v\n", argsFunction)
					if argOne == FUNCTION {
						functionName := argsFunction[1]
						numLocals, err := strconv.Atoi(argsFunction[2])
						if err != nil {
							fmt.Printf("Invalid number of local args for function: %s\n", text)
							return "", nil, err
						}
						text += fmt.Sprintf("(%s)\n", functionName)
						for i := 0; i < numLocals; i++ {
							text += fmt.Sprintf("@0\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n")
						}
					} else if argOne == RETURN {
						text += fmt.Sprintf("@LCL\nD=M\n@R13\nM=D\n\n")
						text += fmt.Sprintf("@5\nA=D-A\nD=M\n@R14\nM=D\n\n")
						text += fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nM=D\n\n")
						text += fmt.Sprintf("@ARG\nD=M+1\n@SP\nM=D\n\n")
						text += fmt.Sprintf("@R13\nAM=M-1\nD=M\n@THAT\nM=D\n\n")
						text += fmt.Sprintf("@R13\nAM=M-1\nD=M\n@THIS\nM=D\n\n")
						text += fmt.Sprintf("@R13\nAM=M-1\nD=M\n@ARG\nM=D\n\n")
						text += fmt.Sprintf("@R13\nAM=M-1\nD=M\n@LCL\nM=D\n\n")
						text += fmt.Sprintf("@R14\nA=M\n0;JMP\n\n")
					} else {
						functionName := argsFunction[1]
						parts := strings.Split(functionName, ".")
						if parts[0] != fileName {
							calledFileNames = append(calledFileNames, parts[0])
						}
						numArgs, err := strconv.Atoi(argsFunction[2])
						if err != nil {
							fmt.Printf("error in string conversion : %v", err)
							return "", nil, err
						}
						num := 0
						if val, ok := funcCounters[functionName]; ok {
							num = val
						} else {
							funcCounters[functionName] = 0
						}
						text += fmt.Sprintf("@%s.RETURN.%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n", functionName, num)
						text += fmt.Sprintf("@LCL\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n")
						text += fmt.Sprintf("@ARG\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n")
						text += fmt.Sprintf("@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n")
						text += fmt.Sprintf("@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n")
						text += fmt.Sprintf("@%d\nD=A\n@5\nD=D+A\n@SP\nD=M-D\n@ARG\nM=D\n\n", numArgs)
						text += fmt.Sprintf("@SP\nD=M\n@LCL\nM=D\n\n")
						text += fmt.Sprintf("@%s\n0;JMP\n\n(%s.RETURN.%d)\n", functionName, functionName, num)
						funcCounters[functionName]++
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

func buildDirTree(root string) (*DirTree, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	var dirTree *DirTree
	if !info.IsDir() {
		parts := strings.Split(info.Name(), ".")
		ext := parts[len(parts)-1]
		if ext == "vm" {
			dirTree = &DirTree{
				Name: info.Name(),
				Path: root,
			}
		}
	} else {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, err
		}
		dirTree = &DirTree{
			Name: info.Name(),
			Path: root,
		}
		var children []*DirTree

		for _, entry := range entries {
			childPath := filepath.Join(root, entry.Name())
			childNode, err := buildDirTree(childPath)
			if err != nil {
				return nil, err
			}
			if childNode != nil {
				children = append(children, childNode)
			}
		}
		if len(children) != 0 {
			dirTree.Children = children
		}
	}
	return dirTree, nil
}

func printDirTree(tree *DirTree, level int) {
	trail := ""
	for i := 0; i < level; i++ {
		trail += "-"
	}
	fmt.Printf("%s%s\n", trail, tree.Name)
	for i := 0; i < len(tree.Children); i++ {
		printDirTree(tree.Children[i], level+1)
	}
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

func isValidBranchingCommand(line string) ([]string, bool) {
	words := strings.Split(line, " ")
	if len(words) != 2 {
		return nil, false
	}
	argOne := words[0]
	if _, okArgOne := branchingArgs[argOne]; !okArgOne {
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
	if _, okArgOne := functionArgs[argOne]; !okArgOne {
		return nil, false
	}
	if length == 1 && argOne != RETURN {
		return nil, false
	}
	if length == 3 && argOne != FUNCTION && argOne != CALL {
		return nil, false
	}
	return words, true
}

func translateMemorySegCommand(args []string, num uint32) (string, error) {
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
			str += fmt.Sprintf("@SP\nA=M\nM=D\n@SP\nM=M+1\n")
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

func cleanText(line string) (string, string) {
	words := strings.SplitN(line, "//", 2)
	text := ""
	if len(words) == 2 {
		text = words[1]
	}
	return strings.TrimRight(strings.TrimLeft(words[0], " \t"), " \t"), text
}
