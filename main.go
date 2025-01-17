package main

import (
	"fmt"
	"os/exec"
)

func main() {
	//DONE: Run the shell command
	//DONE: Loop over it so it get executed multiple times
	//DONE: Add a array of keywords
	//DONE: Loop over the array of keywords so that they are used
	//DONE: Make the keywords array fill in the placeholders

	titles := []string{"dashboard", "team", "permissions"}

	for _, title := range titles {
		cmd := exec.Command("gh", "issue", "create", "-t", title, "-F", "issue-template.md")

		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error creating issue with title %s: %s\n", title, err)
			continue
		}

		fmt.Printf("Output for title %s: %s\n", title, string(output))
	}

}
