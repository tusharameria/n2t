package translator

import "fmt"

const LABEL = "label"
const GO_TO = "goto"
const IF_GO_TO = "if-goto" // Pop the topmost value from the stack, and if that value is not zero, jump to the label.

var BranchingArgs = map[string]bool{
	LABEL:    true,
	GO_TO:    true,
	IF_GO_TO: true,
}

func BranchingCommandTranslate(argsBranch []string, fileName string) string {

	argOne := argsBranch[0]
	argTwo := argsBranch[1]
	text := fmt.Sprintf("//%s %s\n", argOne, argTwo)
	if argOne == LABEL {
		text += fmt.Sprintf("(%s.%s)\n", fileName, argTwo)
	} else if argOne == GO_TO {
		text += fmt.Sprintf("@%s.%s\n0;JMP\n", fileName, argTwo)
	} else {
		text += fmt.Sprintf("@SP\nM=M-1\nA=M\nD=M\n@%s.%s\nD;JNE\n", fileName, argTwo)
	}
	return text
}
