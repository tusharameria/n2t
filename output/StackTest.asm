// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/7/StackArithmetic/StackTest/StackTest.vm

// Executes a sequence of arithmetic and logical operations on the stack. 

// push constant 17
@17
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 17
@17
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// eq
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@EQ_TRUE_0
D;JEQ
@SP
A=M-1
M=0
@EQ_END_0
0;JMP
(EQ_TRUE_0)
@SP
A=M-1
M=-1
(EQ_END_0)

// push constant 17
@17
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 16
@16
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// eq
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@EQ_TRUE_1
D;JEQ
@SP
A=M-1
M=0
@EQ_END_1
0;JMP
(EQ_TRUE_1)
@SP
A=M-1
M=-1
(EQ_END_1)

// push constant 16
@16
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 17
@17
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// eq
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@EQ_TRUE_2
D;JEQ
@SP
A=M-1
M=0
@EQ_END_2
0;JMP
(EQ_TRUE_2)
@SP
A=M-1
M=-1
(EQ_END_2)

// push constant 892
@892
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 891
@891
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// lt
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@LT_TRUE_0
D;JLT
@SP
A=M-1
M=0
@LT_END_0
0;JMP
(LT_TRUE_0)
@SP
A=M-1
M=-1
(LT_END_0)

// push constant 891
@891
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 892
@892
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// lt
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@LT_TRUE_1
D;JLT
@SP
A=M-1
M=0
@LT_END_1
0;JMP
(LT_TRUE_1)
@SP
A=M-1
M=-1
(LT_END_1)

// push constant 891
@891
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 891
@891
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// lt
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@LT_TRUE_2
D;JLT
@SP
A=M-1
M=0
@LT_END_2
0;JMP
(LT_TRUE_2)
@SP
A=M-1
M=-1
(LT_END_2)

// push constant 32767
@32767
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 32766
@32766
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// gt
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@GT_TRUE_0
D;JGT
@SP
A=M-1
M=0
@GT_END_0
0;JMP
(GT_TRUE_0)
@SP
A=M-1
M=-1
(GT_END_0)

// push constant 32766
@32766
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 32767
@32767
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// gt
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@GT_TRUE_1
D;JGT
@SP
A=M-1
M=0
@GT_END_1
0;JMP
(GT_TRUE_1)
@SP
A=M-1
M=-1
(GT_END_1)

// push constant 32766
@32766
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 32766
@32766
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// gt
@SP
M=M-1
A=M
D=M
A=A-1
D=M-D
@GT_TRUE_2
D;JGT
@SP
A=M-1
M=0
@GT_END_2
0;JMP
(GT_TRUE_2)
@SP
A=M-1
M=-1
(GT_END_2)

// push constant 57
@57
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 31
@31
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// push constant 53
@53
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

// push constant 112
@112
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

// neg
@SP
A=M-1
M=-M

// and
@SP
M=M-1
A=M
D=M
A=A-1
M=D&M

// push constant 82
@82
D=A
@SP
M=M+1
A=M
A=A-1
M=D

// or
@SP
M=M-1
A=M
D=M
A=A-1
M=D|M

// not
@SP
A=M-1
M=!M

