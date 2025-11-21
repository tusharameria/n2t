package translator

import "fmt"

func GenerateBootstrapCode(initFileName string) string {
	msg := "// Bootstrap Code\n"
	msg += "@256\nD=A\n@SP\nM=D\n"
	msg += fmt.Sprintf("@%s.RETURN.%d\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n", initFileName, 0)
	msg += "@LCL\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
	msg += "@ARG\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
	msg += "@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
	msg += "@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n\n"
	msg += fmt.Sprintf("@%d\nD=A\n@5\nD=D+A\n@SP\nD=M-D\n@ARG\nM=D\n\n", 0)
	msg += "@SP\nD=M\n@LCL\nM=D\n\n"
	msg += fmt.Sprintf("@%s\n0;JMP\n\n(%s.RETURN.%d)\n", initFileName, initFileName, 0)
	return msg
}
