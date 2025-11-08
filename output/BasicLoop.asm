// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/8/ProgramFlow/BasicLoop/BasicLoop.vm
//
// Computes the sum 1 + 2 + ... + n and pushes the result onto
// the stack. The value n is given in argument[0], which must be 
// initialized by the caller of this code.
//
// push constant 0
@0
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// pop local 0
@0
D=A
@LCL
A=D+M
D=A
@R16
M=D
@SP
M=M-1
A=M
D=M
@R16
A=M
M=D

//label LOOP
(LOOP)

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

// add
@SP
M=M-1
A=M
D=M
A=A-1
M=D+M

// pop local 0
@0
D=A
@LCL
A=D+M
D=A
@R16
M=D
@SP
M=M-1
A=M
D=M
@R16
A=M
M=D

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

// push constant 1
@1
D=A
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

// pop argument 0
@0
D=A
@ARG
A=D+M
D=A
@R16
M=D
@SP
M=M-1
A=M
D=M
@R16
A=M
M=D

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

//if-goto LOOP
@SP
M=M-1
A=M
D=M
@LOOP
D;JNE

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

