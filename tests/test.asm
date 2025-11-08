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