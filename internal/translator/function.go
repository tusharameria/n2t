package translator

import (
	"fmt"
	"strconv"
	"strings"
)

const FUNCTION = "function"
const RETURN = "return"
const CALL = "call"

var FunctionArgs = map[string]bool{
	FUNCTION: true,
	RETURN:   true,
	CALL:     true,
}

func FunctionCommandTranslate(argsFunction []string, fileName string, funcCounters map[string]int, calledFileNames []string) (string, []string, error) {

	argOne := argsFunction[0]
	text := fmt.Sprintf("// %v\n", argsFunction)
	if argOne == FUNCTION {
		functionName := argsFunction[1]
		numLocals, err := strconv.Atoi(argsFunction[2])
		if err != nil {
			fmt.Printf("Invalid number of local args for function: %s\n", text)
			return "", nil, err
		}
		text += fmt.Sprintf("(%s)\n", functionName)
		for i := 0; i < numLocals; i++ {
			text += "@0\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		}
	} else if argOne == RETURN {
		text += "@LCL\nD=M\n@R13\nM=D\n\n"
		text += "@5\nA=D-A\nD=M\n@R14\nM=D\n\n"
		text += "@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nM=D\n\n"
		text += "@ARG\nD=M+1\n@SP\nM=D\n\n"
		text += "@R13\nAM=M-1\nD=M\n@THAT\nM=D\n\n"
		text += "@R13\nAM=M-1\nD=M\n@THIS\nM=D\n\n"
		text += "@R13\nAM=M-1\nD=M\n@ARG\nM=D\n\n"
		text += "@R13\nAM=M-1\nD=M\n@LCL\nM=D\n\n"
		text += "@R14\nA=M\n0;JMP\n\n"
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
		text += "@LCL\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
		text += "@ARG\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
		text += "@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
		text += "@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
		text += fmt.Sprintf("@%d\nD=A\n@5\nD=D+A\n@SP\nD=M-D\n@ARG\nM=D\n\n", numArgs)
		text += "@SP\nD=M\n@LCL\nM=D\n\n"
		text += fmt.Sprintf("@%s\n0;JMP\n\n(%s.RETURN.%d)\n", functionName, functionName, num)
		funcCounters[functionName]++
	}
	return text, calledFileNames, nil
}
