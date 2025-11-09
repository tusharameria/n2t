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

// function Foo 2
(FOO)
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

@5\nA=D-A\nD=M\n@R14\nM=D\n

@SP\nM=M-1\nA=M\nD=M\n@ARG\nA=M\nM=D\n
@ARG\nD=M+1\n@SP\nM=D\n
@R13\nAM=M-1\nD=M\n@THAT\nM=D\n
@R13\nAM=M-1\nD=M\n@THIS\nM=D\n
@R13\nAM=M-1\nD=M\n@ARG\nM=D\n
@R13\nAM=M-1\nD=M\n@LCL\nM=D\n
(END)\n@END\n0;JMP\n