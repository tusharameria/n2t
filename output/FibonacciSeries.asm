// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/8/ProgramFlow/FibonacciSeries/FibonacciSeries.vm
//
// Puts the first n elements of the Fibonacci series in the memory,
// starting at address addr. n and addr are given in argument[0] and
// argument[1], which must be initialized by the caller of this code.
//
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

// pop pointer 1
@SP
M=M-1
@SP
A=M
D=M
@THAT
M=D

// push constant 0
@0
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// pop that 0
@0
D=A
@THAT
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

// push constant 1
@1
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// pop that 1
@1
D=A
@THAT
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

// push constant 2
@2
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

//
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

//if-goto COMPUTE_ELEMENT
@SP
M=M-1
A=M
D=M
@COMPUTE_ELEMENT
D;JNE

//goto END
@END
0;JMP

//
//label COMPUTE_ELEMENT
(COMPUTE_ELEMENT)

// that[2] = that[0] + that[1]
// push that 0
@0
D=A
@THAT
A=D+M
D=M
@SP
M=M+1
A=M
A=A-1
M=D

// push that 1
@1
D=A
@THAT
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

// pop that 2
@2
D=A
@THAT
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

// THAT += 1 (updates the base address of that)
// push pointer 1
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1

// push constant 1
@1
D=A
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

// pop pointer 1
@SP
M=M-1
@SP
A=M
D=M
@THAT
M=D

// updates n-- and loops          
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

//goto LOOP
@LOOP
0;JMP

//
//label END
(END)

