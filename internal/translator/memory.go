package translator

import "fmt"

const TEMP_INIT_ADDRESS = 5

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

var MemoryArgOne = map[string]bool{
	PUSH: true,
	POP:  true,
}
var MemoryArgTwo = map[string]string{
	LOCAL:    "LCL",
	ARGUMENT: "ARG",
	THIS:     "THIS",
	THAT:     "THAT",
	CONSTANT: "CONST",
	STATIC:   "STATIC",
	POINTER:  "POINT",
	TEMP:     "R5",
}

func MemorySegCommandTranslate(args []string, num uint32, fileName string) (string, error) {
	argOne := args[0]
	argTwo := args[1]

	str := ""
	if argOne == PUSH {
		argTwoVal, ok := MemoryArgTwo[argTwo]
		if !ok {
			return "", fmt.Errorf("%s doesn't exist in MemoryArgTwo", argTwo)
		}
		if argTwo == POINTER {
			buff := MemoryArgTwo[THIS]
			if num == 1 {
				buff = MemoryArgTwo[THAT]
			}
			str += fmt.Sprintf("@%s\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", buff)
		} else if argTwo == STATIC {
			str += fmt.Sprintf("@%s.%d\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n", fileName, num)
		} else {
			str += fmt.Sprintf("@%d\nD=A\n", num)
			if argTwo != CONSTANT {
				if argTwo == TEMP {
					str += fmt.Sprintf("@%d\nA=D+A\nD=M\n", TEMP_INIT_ADDRESS)
				} else {
					str += fmt.Sprintf("@%s\nA=D+M\nD=M\n", argTwoVal)
				}
			}
			str += "@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		}
	}

	if argOne == POP {
		if argTwo == CONSTANT {
			return "", fmt.Errorf("can't POP a constant")
		}
		argTwoVal, ok := MemoryArgTwo[argTwo]
		if !ok {
			return "", fmt.Errorf("%s doesn't exist in MemoryArgTwo", argTwo)
		}

		if argTwo == POINTER {
			buff := MemoryArgTwo[THIS]
			if num == 1 {
				buff = MemoryArgTwo[THAT]
			}
			str += fmt.Sprintf("@SP\nM=M-1\n@SP\nA=M\nD=M\n@%s\nM=D\n", buff)
		} else if argTwo == STATIC {
			str += fmt.Sprintf("@SP\nM=M-1\n@SP\nA=M\nD=M\n@%s.%d\nM=D\n", fileName, num)
		} else {
			str += fmt.Sprintf("@%d\nD=A\n", num)
			if argTwo == TEMP {
				str += fmt.Sprintf("@%d\nD=D+A\n", TEMP_INIT_ADDRESS)
			} else {
				str += fmt.Sprintf("@%s\nA=D+M\nD=A\n", argTwoVal)
			}
			str += "@R16\nM=D\n@SP\nM=M-1\nA=M\nD=M\n@R16\nA=M\nM=D\n"
		}
	}

	return fmt.Sprintf("// %s %s %d\n%s", argOne, argTwo, num, str), nil
}
