// label LOOP
(LOOP)
// goto LOOP
@LOOP
0;JMP
// if-goto LOOP
@SP
M=M-1
A=M
D=M
@LOOP
D;JNE

// call Foo 1
// save return address
@FOO.RETURN
D=A
@SP
A=M
M=D
@SP
M=M+1

// save caller LCL
@LCL
D=M
@SP
A=M
M=D
@SP
M=M+1

// save caller ARG
@ARG
D=M
@SP
A=M
M=D
@SP
M=M+1

// save caller THIS
@THIS
D=M
@SP
A=M
M=D
@SP
M=M+1

// save caller THAT
@THAT
D=M
@SP
A=M
M=D
@SP
M=M+1

// reposition ARG for callee (ARG = SP - 5 - nArgs)
@SP
D=M
@5
D=D-A
@1
D=D-A
@ARG
M=D

// reposition LCL for callee (LCL = SP)
@SP
D=M
@LCL
M=D

@FOO
0;JMP

(FOO.RETURN)

// function Foo 2
(FOO)
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

// return
@LCL
D=M
@R13
M=D

@5
A=D-A
D=M
@R14
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

@FOO.RETURN
0;JMP
