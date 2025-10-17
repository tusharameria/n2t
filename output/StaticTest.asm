// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/7/MemoryAccess/StaticTest/StaticTest.vm

// Executes pop and push commands using the static segment.

// push constant 111
@111
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 333
@333
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 888
@888
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// pop static 8
@8
D=A
@STATIC
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

// pop static 3
@3
D=A
@STATIC
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

// pop static 1
@1
D=A
@STATIC
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

// push static 3
@3
D=A
@STATIC
A=D+M
D=M
@SP
M=M+1
A=M
A=A-1
M=D

// push static 1
@1
D=A
@STATIC
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

// push static 8
@8
D=A
@STATIC
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

