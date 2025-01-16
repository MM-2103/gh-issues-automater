package main

import (
	"fmt"
	"os/exec"
)

func main() {
	//DONE: Run the shell command
	//TODO: Loop over it so it get executed multiple times
	//TODO: Add a array of keywords
	//TODO: Loop over the array of keywords so that they are used
	//TODO: Make the keywords array fill in the placeholders

	cmd := exec.Command("gh", "issue", "create", "-t", "%s", "-b", "test2")

	out, err := cmd.Output()
	if err != nil {
		fmt.Println("could not run command: ", err)
	}

	fmt.Println(string(out))
}
