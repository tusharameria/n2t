package translator

import (
	"fmt"
	"strings"
)

type Counters struct {
	Eq int
	Gt int
	Lt int
}

const EQ_TRUE = "EQ_TRUE"
const EQ_END = "EQ_END"
const GT_TRUE = "GT_TRUE"
const GT_END = "GT_END"
const LT_TRUE = "LT_TRUE"
const LT_END = "LT_END"

const ADD = "add"
const SUB = "sub"
const NEG = "neg"
const EQ = "eq"
const GT = "gt"
const LT = "lt"
const AND = "and"
const OR = "or"
const NOT = "not"

var ArithmaticLogicArgs = map[string]string{
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

func ArithmaticLogicCommandTranslate(text string, counters *Counters) (string, error) {
	buff, ok := ArithmaticLogicArgs[text]
	if !ok {
		fmt.Printf("%s is not a valid alArgs key\n", text)
		return "", fmt.Errorf("%s is not a valid alArgs key", text)
	}
	if text == EQ || text == GT || text == LT {
		if text == EQ {
			buff = strings.ReplaceAll(buff, EQ_TRUE, fmt.Sprintf("EQ_TRUE_%d", counters.Eq))
			buff = strings.ReplaceAll(buff, EQ_END, fmt.Sprintf("EQ_END_%d", counters.Eq))
			counters.Eq++
		}
		if text == GT {
			buff = strings.ReplaceAll(buff, GT_TRUE, fmt.Sprintf("GT_TRUE_%d", counters.Gt))
			buff = strings.ReplaceAll(buff, GT_END, fmt.Sprintf("GT_END_%d", counters.Gt))
			counters.Gt++
		}
		if text == LT {
			buff = strings.ReplaceAll(buff, LT_TRUE, fmt.Sprintf("LT_TRUE_%d", counters.Lt))
			buff = strings.ReplaceAll(buff, LT_END, fmt.Sprintf("LT_END_%d", counters.Lt))
			counters.Lt++
		}
	}
	return buff, nil
}
