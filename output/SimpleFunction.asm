// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/8/FunctionCalls/SimpleFunction/SimpleFunction.vm
//
// Performs a simple calculation and returns the result.
// argument[0] and argument[1] must be set by the caller.
//
// [function SimpleFunction.test 2]
(SimpleFunction.test)
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
@0
D=A
@SP
A=M
M=D
@SP
M=M+1

// push local 0
@0
D=A
@LCL
A=D+M
D=M
@SP
M=M+1
A=M
A=A-1
M=D

// push local 1
@1
D=A
@LCL
A=D+M
D=M
@SP
M=M+1
A=M
A=A-1
M=D

// add
@SP
M=M-1
A=M
D=M
A=A-1
M=D+M

// not
@SP
A=M-1
M=!M

// push argument 0
@0
D=A
@ARG
A=D+M
D=M
@SP
M=M+1
A=M
A=A-1
M=D

// add
@SP
M=M-1
A=M
D=M
A=A-1
M=D+M

// push argument 1
@1
D=A
@ARG
A=D+M
D=M
@SP
M=M+1
A=M
A=A-1
M=D

// sub
@SP
M=M-1
A=M
D=M
A=A-1
M=M-D

// [return]
@LCL
D=M
@R13
M=D

@SP
M=M-1
A=M
D=M
@ARG
A=M
M=D

@ARG
D=M+1
@SP
M=D

@R13
AM=M-1
D=M
@THAT
M=D

@R13
AM=M-1
D=M
@THIS
M=D

@R13
AM=M-1
D=M
@ARG
M=D

@R13
AM=M-1
D=M
@LCL
M=D

(END)
@END
0;JMP


