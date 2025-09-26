// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/7/MemoryAccess/BasicTest/BasicTest.vm

// Executes pop and push commands.

// push constant 10
@10
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

// push constant 21
@21
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 22
@22
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// pop argument 2
@2
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

// pop argument 1
@1
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

// push constant 36
@36
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// pop this 6
@6
D=A
@THIS
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

// push constant 42
@42
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 45
@45
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// pop that 5
@5
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

// push constant 510
@510
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// pop temp 6
@6
D=A
@5
D=D+A
@R16
M=D
@SP
M=M-1
A=M
D=M
@R16
A=M
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

// push that 5
@5
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

// push this 6
@6
D=A
@THIS
A=D+M
D=M
@SP
M=M+1
A=M
A=A-1
M=D

// push this 6
@6
D=A
@THIS
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

// sub
@SP
M=M-1
A=M
D=M
A=A-1
M=M-D

// push temp 6
@6
D=A
@5
A=D+A
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

